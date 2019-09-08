# server
--
    import "github.com/rickbassham/example-go/chiapi/server"

Package server contains an HTTP server implementation that supports graceful
shutdowns. By gracefully shutting down, we allow existing requests to finish
before shutting the server down.

## Usage

#### type GracefulHTTPServer

```go
type GracefulHTTPServer struct {
}
```

GracefulHTTPServer will run a new http server that supports graceful shutdown.

#### func  NewGracefulHTTPServer

```go
func NewGracefulHTTPServer(log *zap.Logger, svr HTTPServer, l net.Listener, timeout time.Duration) *GracefulHTTPServer
```
NewGracefulHTTPServer creates a new GracefulHTTPServer.

#### func (*GracefulHTTPServer) Run

```go
func (s *GracefulHTTPServer) Run() error
```
Run starts the http server, waiting for a SIGINT signal to the process. Once it
has received one of those, or if it fails to start at all, the function will
return. If we are shutting down due to a signal, we call Shutdown on the
http.Server we are using. This function will wait for existing requests to be
completed before returning. We only wait up to 30 seconds to do the graceful
shutdown. After that, we just kill the connections.

#### type HTTPServer

```go
type HTTPServer interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}
```

HTTPServer is an interface that the go http.Server struct will satisfy.
