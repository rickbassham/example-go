# router
--
    import "github.com/rickbassham/example-go/chiapi/router"

Package router is used to define the routes to our request handlers.

## Usage

#### func  NewRouter

```go
func NewRouter(h Handler, log *zap.Logger, nr newrelic.Application, tokenAuth *jwtauth.JWTAuth, version, corsOrigin string) http.Handler
```
NewRouter creates a new CORS enabled router for our API. All requests will be
logged and instrumented with New Relic.

#### type Handler

```go
type Handler interface {
	Health(w http.ResponseWriter, r *http.Request)
	Protected(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	Unauthorized(w http.ResponseWriter, r *http.Request)
}
```

Handler exposes the functions for handling web requests.
