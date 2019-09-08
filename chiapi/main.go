// Service chiapi is an example API written in Go. It is fully logged, and instrumented with New
// Relic. It also shows basic JWT parsing for authentication, and use of the go-chi router.
package main

import (
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
)

func main() {
	var err error

	defer func() {
		// If we are existing due to an error, be sure to set the exit code appropriately.
		if err != nil {
			os.Exit(1)
		}
	}()

	logEnc := zap.NewProductionEncoderConfig()
	logEnc.EncodeTime = zapcore.ISO8601TimeEncoder
	logEnc.TimeKey = "timestamp"

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEnc),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevel(),
	))

	log = log.With(
		zap.String("app_name", "chiapi"),
		zap.String("environment", "development"),
		zap.String("team", "rickbassham"),
		zap.String("run_id", uuid.New().String()),
		zap.Time("start_time", time.Now()),
	)

	log.Info("initializing")

	nr, err := newrelic.NewApplication(newrelic.Config{
		AppName: "chiapi-production",
		Labels: map[string]string{
			"Team":        "rickbassham",
			"Environment": "development",
		},
		License: "0123456789012345678901234567890123456789",
	})
	if err != nil {
		log.Error("error creating newrelic app", zap.Error(err))
		return
	}

	// Give new relic 30 seconds to send instrumentation before terminating.
	defer nr.Shutdown(30 * time.Second)

	jwtAuth := jwtauth.New("HS256", []byte("env.AuthSecret"), nil)

	h := &handler.Handler{}

	r := router.NewRouter(h, log, nr, jwtAuth, "0.0.1", "http://localhost")

	err = startHTTPServer(r, log, ":3000")
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
