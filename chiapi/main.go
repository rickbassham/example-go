// Service chiapi is an example API written in Go. It is fully logged, and instrumented with New
// Relic. It also shows basic JWT parsing for authentication, and use of the go-chi router.
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rickbassham/example-go/chiapi/handler"
	"github.com/rickbassham/example-go/chiapi/router"
	"github.com/rickbassham/example-go/chiapi/server"
	"github.com/rickbassham/example-go/pkg/env"
)

type config struct {
	env.Config
	ListenAddress string `env:"LISTEN_ADDRESS,required"`
	JWTAuthSecret string `env:"JWT_AUTH_SECRET,required"`
	CORSOrigin    string `env:"CORS_ORIGIN,required"`
}

func main() {
	var err error

	defer func() {
		// If we are existing due to an error, be sure to set the exit code appropriately.
		if err != nil {
			os.Exit(1)
		}
	}()

	var c config
	err = env.Load(&c)
	log := startLogger(c)

	if err != nil {
		log.Error("error initializing environment", zap.Error(err))
		return
	}

	log.Info("initializing")

	nr, err := startNewRelic(c)
	if err != nil {
		log.Error("error creating newrelic app", zap.Error(err))
		return
	}
	// Give new relic 30 seconds to send instrumentation before terminating.
	defer nr.Shutdown(30 * time.Second)

	jwtAuth := jwtauth.New("HS256", []byte(c.JWTAuthSecret), nil)

	h := &handler.Handler{}

	r := router.NewRouter(h, log, nr, jwtAuth, c.BuildGitTag, c.CORSOrigin)

	err = startHTTPServer(r, log, c.ListenAddress)
}

func startLogger(c config) *zap.Logger {
	logEnc := zap.NewProductionEncoderConfig()
	logEnc.EncodeTime = zapcore.ISO8601TimeEncoder
	logEnc.TimeKey = "timestamp"

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEnc),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevel(),
	))

	log = log.With(
		zap.String("app_name", c.AppName),
		zap.String("environment", c.Environment),
		zap.String("build_git_hash", c.BuildGitHash),
		zap.String("build_git_tag", c.BuildGitTag),
		zap.Time("build_date", c.BuildDate),
		zap.String("team", "rickbassham"),
		zap.String("run_id", uuid.New().String()),
		zap.Time("start_time", time.Now()),
	)

	return log
}

func startNewRelic(c config) (newrelic.Application, error) {
	nr, err := newrelic.NewApplication(newrelic.Config{
		AppName: fmt.Sprintf("%s-%s", c.AppName, c.Environment),
		Labels: map[string]string{
			"Team":        "rickbassham",
			"Environment": c.Environment,
			"Version":     c.BuildGitTag,
		},
		License: c.NewRelicLicense,
	})
	if err != nil {
		return nil, err
	}

	err = nr.WaitForConnection(30 * time.Second)
	if err != nil {
		return nil, err
	}

	return nr, err
}

func startHTTPServer(h http.Handler, log *zap.Logger, serverAddr string) error {
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: h,
	}

	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Error("error starting listener", zap.Error(err))
		return err
	}

	log.Info("starting server", zap.String("addr", serverAddr))

	// Start the http server.
	s := server.NewGracefulHTTPServer(log, httpServer, ln, 30*time.Second)
	return s.Run()
}
