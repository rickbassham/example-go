# handler
--
    import "github.com/rickbassham/example-go/chiapi/handler"

Package handler implements the request handlers for the API.

## Usage

#### type Handler

```go
type Handler struct {
}
```

Handler has all the functions needed to serve our api.

#### func (*Handler) Health

```go
func (h *Handler) Health(w http.ResponseWriter, r *http.Request)
```
Health always returns a 200 response.

#### func (*Handler) NotFound

```go
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request)
```
NotFound is called when the request is for an unknown resource.

#### func (*Handler) Protected

```go
func (h *Handler) Protected(w http.ResponseWriter, r *http.Request)
```
Protected is a protected endpoint example.

#### func (*Handler) Unauthorized

```go
func (h *Handler) Unauthorized(w http.ResponseWriter, r *http.Request)
```
Unauthorized is called when the request is not authorized.

#### type SimpleResponse

```go
type SimpleResponse struct {
	XMLName xml.Name `json:"-" xml:"simpleResponse"`
	TraceID string   `json:"trace_id" xml:"traceId,attr"`
	Message string   `json:"message" xml:",innerxml"`
}
```

SimpleResponse is used to send a meaningful message back to the caller, with
trace id to debug later.
