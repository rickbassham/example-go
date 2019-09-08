// Package server contains an HTTP server implementation that supports graceful
// shutdowns. By gracefully shutting down, we allow existing requests to finish
// before shutting the server down.
package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

// HTTPServer is an interface that the go http.Server struct will satisfy.
type HTTPServer interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}

// GracefulHTTPServer will run a new http server that supports graceful shutdown.
type GracefulHTTPServer struct {
	svr HTTPServer
	l   net.Listener
	log *zap.Logger

	timeout time.Duration
}

// NewGracefulHTTPServer creates a new GracefulHTTPServer.
func NewGracefulHTTPServer(log *zap.Logger, svr HTTPServer, l net.Listener, timeout time.Duration) *GracefulHTTPServer {
	return &GracefulHTTPServer{
		log:     log,
		svr:     svr,
		l:       l,
		timeout: timeout,
	}
}

// Run starts the http server, waiting for a SIGINT signal to
// the process. Once it has received one of those, or if it fails to start at all,
// the function will return. If we are shutting down due to a signal, we call
// Shutdown on the http.Server we are using. This function will wait for existing
// requests to be completed before returning. We only wait up to 30 seconds to
// do the graceful shutdown. After that, we just kill the connections.
func (s *GracefulHTTPServer) Run() error {
	errs := make(chan error, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		err := s.svr.Serve(s.l)
		if err != nil {
			errs <- err
		}
	}()

	// This select statement will block until we can read from EITHER our errs
	// channel or the stop channel. The stop channel will get a value when we get
	// a SIGINT signal. The errs channel will get a value if we failed
	// to start the server.
	select {
	case err := <-errs:
		s.log.Error("")
		return err
	case sig := <-stop:
		s.log.Info("server shutdown request received", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	err := s.svr.Shutdown(ctx)
	cancel() // Cancel the timeout, since we already finished.

	return err
}
