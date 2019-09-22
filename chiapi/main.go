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
	"github.com/go-redis/redis/v7"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/chiapi/handler"
	"github.com/rickbassham/example-go/chiapi/router"
	"github.com/rickbassham/example-go/chiapi/server"
	"github.com/rickbassham/example-go/pkg/cache"
	"github.com/rickbassham/example-go/pkg/env"
	"github.com/rickbassham/example-go/pkg/logging"
)

type config struct {
	env.Config
	ListenAddress string `env:"LISTEN_ADDRESS,required"`
	JWTAuthSecret string `env:"JWT_AUTH_SECRET,required"`
	CORSOrigin    string `env:"CORS_ORIGIN,required"`
	RedisAddress  string `env:"REDIS_ADDRESS,required"`
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
	log := logging.Initialize(c.Config)

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

	rc := redis.NewClient(&redis.Options{
		Addr: c.RedisAddress,
	})

	_, err = rc.Ping().Result()
	if err != nil {
		log.Error("error pinging redis", zap.Error(err))
		return
	}

	appCache := cache.New(rc)

	h := handler.New(appCache)

	r := router.NewRouter(h, log, nr, jwtAuth, c.BuildGitTag, c.CORSOrigin)

	err = startHTTPServer(r, log, c.ListenAddress)
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
