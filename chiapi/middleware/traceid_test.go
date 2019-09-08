package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

func TestTraceID(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	middleware.TraceID(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Len(t, w.Header().Get("X-Trace-Id"), 36)
}

func TestTraceID_Existing(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("X-Trace-Id", "my-cool-trace-id")

	middleware.TraceID(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "my-cool-trace-id", w.Header().Get("X-Trace-Id"))
}
