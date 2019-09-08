package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-chi/jwtauth"

	"github.com/stretchr/testify/assert"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

func TestAuthenticator_MissingToken(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	middleware.Authenticator(nil)(h).ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code)
}

func TestAuthenticator_InvalidToken(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	tok := &jwt.Token{
		Valid:  false,
		Claims: jwt.MapClaims{},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	r = r.WithContext(jwtauth.NewContext(r.Context(), tok, nil))

	middleware.Authenticator(nil)(h).ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code)
}

func TestAuthenticator_ValidToken(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	tok := &jwt.Token{
		Valid:  true,
		Claims: jwt.MapClaims{},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	r = r.WithContext(jwtauth.NewContext(r.Context(), tok, nil))

	middleware.Authenticator(nil)(h).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
