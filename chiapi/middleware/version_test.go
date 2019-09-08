package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

func TestVersion(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	middleware.Version("my-version")(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "my-version", w.Header().Get("X-Version"))
}
