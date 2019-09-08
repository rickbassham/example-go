package httputil_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rickbassham/example-go/pkg/httputil"
	"github.com/stretchr/testify/mock"
)

type mockTransport struct {
	mock.Mock
}

func (m *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	args := m.Called(r)

	return args.Get(0).(*http.Response), args.Error(1)
}

func TestAPIKeyTransport(t *testing.T) {
	old := &mockTransport{}

	old.On("RoundTrip", mock.Anything).Return(&http.Response{}, nil).Run(func(args mock.Arguments) {
		r := args.Get(0).(*http.Request)

		assert.Equal(t, "apikey", r.Header.Get("X-Api-Key"))
	})

	rt := httputil.APIKeyTransport("apikey", old)

	r, err := http.NewRequest("GET", "http://test/api", nil)
	require.NoError(t, err)

	_, err = rt.RoundTrip(r)

	require.NoError(t, err)
}

func TestBasicAuthTransport(t *testing.T) {
	old := &mockTransport{}

	old.On("RoundTrip", mock.Anything).Return(&http.Response{}, nil).Run(func(args mock.Arguments) {
		r := args.Get(0).(*http.Request)

		assert.Equal(t, "Basic dGVzdDpwYXNz", r.Header.Get("Authorization"))
	})

	rt := httputil.BasicAuthTransport("test", "pass", old)

	r, err := http.NewRequest("GET", "http://test/api", nil)
	require.NoError(t, err)

	_, err = rt.RoundTrip(r)

	require.NoError(t, err)
}
