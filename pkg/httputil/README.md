# httputil
--
    import "github.com/rickbassham/example-go/pkg/httputil"


## Usage

#### func  APIKeyTransport

```go
func APIKeyTransport(apiKey string, old http.RoundTripper) http.RoundTripper
```
APIKeyTransport is an http.RoundTripper that will add basic authentication
headers to all requests.

#### func  BasicAuthTransport

```go
func BasicAuthTransport(username, password string, old http.RoundTripper) http.RoundTripper
```
BasicAuthTransport is an http.RoundTripper that will add basic authentication
headers to all requests.

#### func  DefaultLogTransport

```go
func DefaultLogTransport(log *zap.Logger, old http.RoundTripper) http.RoundTripper
```
DefaultLogTransport will add the given logger to all request contexts.

#### func  HeaderTransport

```go
func HeaderTransport(key, value string, old http.RoundTripper) http.RoundTripper
```
HeaderTransport is an http.RoundTripper that will add a header to all requests.

#### func  LogTransport

```go
func LogTransport(old http.RoundTripper) http.RoundTripper
```
LogTransport will log every outgoing request.

#### func  TraceIDTransport

```go
func TraceIDTransport(old http.RoundTripper) http.RoundTripper
```
TraceIDTransport ensures a trace id header is set on all outgoing requests.
