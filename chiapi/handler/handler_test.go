package handler_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/rickbassham/example-go/chiapi/handler"
	"github.com/rickbassham/example-go/pkg/identity"
	"github.com/rickbassham/example-go/pkg/tracing"
)

func TestHealth_JSON(t *testing.T) {
	h := &handler.Handler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	r = r.WithContext(tracing.WithTraceID(r.Context(), "my-trace-id"))

	h.Health(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, "{\"trace_id\":\"my-trace-id\",\"message\":\"OK\"}\n", string(body))
}

func TestHealth_XML(t *testing.T) {
	h := &handler.Handler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	r = r.WithContext(tracing.WithTraceID(r.Context(), "my-trace-id"))
	r.Header.Add("Accept", "text/xml")

	h.Health(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, "<simpleResponse traceId=\"my-trace-id\">OK</simpleResponse>", string(body))
}

func TestProtected(t *testing.T) {
	h := &handler.Handler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/protected/1", nil)
	r = r.WithContext(tracing.WithTraceID(r.Context(), "my-trace-id"))
	r = r.WithContext(withChiParam(r.Context(), "id", "1"))
	r = r.WithContext(identity.WithUser(r.Context(), "my-user-name"))

	h.Protected(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, "{\"trace_id\":\"my-trace-id\",\"message\":\"your username is: my-user-name; the id you requested is: 1\"}\n", string(body))
}

func withChiParam(ctx context.Context, key, value string) context.Context {
	var rctx *chi.Context
	var ok bool
	if rctx, ok = ctx.Value(chi.RouteCtxKey).(*chi.Context); !ok {
		rctx = chi.NewRouteContext()
	}

	rctx.URLParams.Add(key, value)

	return context.WithValue(ctx, chi.RouteCtxKey, rctx)
}
