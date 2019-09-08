# middleware
--
    import "github.com/rickbassham/example-go/chiapi/middleware"

Package middleware contains various server middleware to help with requests,
including logging, instrumentation, authentication, and debugging middleware.

## Usage

#### func  Authenticator

```go
func Authenticator(unauthorized http.Handler) func(next http.Handler) http.Handler
```
Authenticator is used to verify a valid token is on the request. The param
unauthorized should be a request handler to render your 401 response. If it is
nil, a simple default response will be written.

#### func  GetLogger

```go
func GetLogger(ctx context.Context) *zap.Logger
```
GetLogger retrieves the logger from the context. If there is no logger on the
context, it will return a valid NoOp logger.

#### func  GetTraceID

```go
func GetTraceID(ctx context.Context) string
```
GetTraceID retrieves the traceID from the request context.

#### func  GetUser

```go
func GetUser(ctx context.Context) string
```
GetUser retrieves the user from the context.

#### func  Logger

```go
func Logger(log *zap.Logger) func(next http.Handler) http.Handler
```
Logger middleware logs each request and adds the logger to the request context
for use in request handlers.

#### func  NewRelicChiRouter

```go
func NewRelicChiRouter(nr newrelic.Application) func(next http.Handler) http.Handler
```
NewRelicChiRouter will start a newrelic transaction for the request. It will
name the transaction to match the chi route. It will also add custom attributes
to the transaction, such as the trace id, any url parameters, and query string
parameters. These attributes make it easy to tie your New Relic traces to logs.
This middleware can be used in place of using newrelic.WrapHandleFunc
everywhere.

#### func  TraceID

```go
func TraceID(next http.Handler) http.Handler
```
TraceID tries to read the X-Trace-Id request header, and if empty, creates a new
traceID to use. This traceID will be added to the response X-Trace-Id header and
to the request context.

#### func  User

```go
func User(next http.Handler) http.Handler
```
User get's user email from JWT token and adds it to the request context.

#### func  Version

```go
func Version(v string) func(next http.Handler) http.Handler
```
Version middleware adds the X-Version header to the response.

#### func  WithLogger

```go
func WithLogger(ctx context.Context, l *zap.Logger) context.Context
```
WithLogger adds the logger to the request context.

#### func  WithTraceID

```go
func WithTraceID(ctx context.Context, traceID string) context.Context
```
WithTraceID adds the traceID to the request context.

#### func  WithUser

```go
func WithUser(ctx context.Context, user string) context.Context
```
WithUser adds the user to the request context.
