package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"

	"github.com/rickbassham/example-go/chiapi/middleware"
	"github.com/rickbassham/example-go/pkg/identity"
)

func TestUser_NoToken(t *testing.T) {
	var user string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user = identity.FromContext(r.Context())
		w.WriteHeader(200)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	middleware.User(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", user)
}

func TestUser_WithToken(t *testing.T) {
	var user string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user = identity.FromContext(r.Context())
		w.WriteHeader(200)
	})

	tok := &jwt.Token{
		Valid: true,
		Claims: jwt.MapClaims{
			"email": "my-user",
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	r = r.WithContext(jwtauth.NewContext(r.Context(), tok, nil))

	middleware.User(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "my-user", user)
}
