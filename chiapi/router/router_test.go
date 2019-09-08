package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-chi/jwtauth"

	"go.uber.org/zap"

	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rickbassham/example-go/chiapi/handler"
	"github.com/rickbassham/example-go/chiapi/router"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Health(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *mockHandler) Protected(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *mockHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *mockHandler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

type mockNewRelicApp struct {
	mock.Mock
}

func (m *mockNewRelicApp) StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	return m.Called(name, w, r).Get(0).(newrelic.Transaction)
}

func (m *mockNewRelicApp) RecordCustomEvent(eventType string, params map[string]interface{}) error {
	return m.Called(eventType, params).Error(0)
}

func (m *mockNewRelicApp) RecordCustomMetric(name string, value float64) error {
	return m.Called(name, value).Error(0)
}

func (m *mockNewRelicApp) WaitForConnection(timeout time.Duration) error {
	return m.Called(timeout).Error(0)
}

func (m *mockNewRelicApp) Shutdown(timeout time.Duration) {
	m.Called(timeout)
}

type mockNewRelicTxn struct {
	mock.Mock

	w http.ResponseWriter
}

func (m *mockNewRelicTxn) Header() http.Header {
	m.Called()
	return m.w.Header()
}

func (m *mockNewRelicTxn) Write(d []byte) (int, error) {
	m.Called(d)
	return m.w.Write(d)
}

func (m *mockNewRelicTxn) WriteHeader(statusCode int) {
	m.Called(statusCode)
	m.w.WriteHeader(statusCode)
}

func (m *mockNewRelicTxn) End() error {
	return m.Called().Error(0)
}

func (m *mockNewRelicTxn) Ignore() error {
	return m.Called().Error(0)
}

func (m *mockNewRelicTxn) SetName(name string) error {
	return m.Called(name).Error(0)
}

func (m *mockNewRelicTxn) NoticeError(err error) error {
	return m.Called(err).Error(0)
}

func (m *mockNewRelicTxn) AddAttribute(key string, value interface{}) error {
	return m.Called(key, value).Error(0)
}

func (m *mockNewRelicTxn) SetWebRequest(r newrelic.WebRequest) error {
	return m.Called(r).Error(0)
}

func (m *mockNewRelicTxn) SetWebResponse(w http.ResponseWriter) newrelic.Transaction {
	return m.Called(w).Get(0).(newrelic.Transaction)
}

func (m *mockNewRelicTxn) StartSegmentNow() newrelic.SegmentStartTime {
	return m.Called().Get(0).(newrelic.SegmentStartTime)
}

func (m *mockNewRelicTxn) CreateDistributedTracePayload() newrelic.DistributedTracePayload {
	return m.Called().Get(0).(newrelic.DistributedTracePayload)
}

func (m *mockNewRelicTxn) AcceptDistributedTracePayload(t newrelic.TransportType, payload interface{}) error {
	return m.Called(t, payload).Error(0)
}

func (m *mockNewRelicTxn) Application() newrelic.Application {
	return m.Called().Get(0).(newrelic.Application)
}

func (m *mockNewRelicTxn) BrowserTimingHeader() (*newrelic.BrowserTimingHeader, error) {
	args := m.Called()
	return args.Get(0).(*newrelic.BrowserTimingHeader), args.Error(1)
}

func (m *mockNewRelicTxn) NewGoroutine() newrelic.Transaction {
	return m.Called().Get(0).(newrelic.Transaction)
}

func TestRouter_Health(t *testing.T) {
	log := zap.NewExample()

	h := &mockHandler{}
	nr := &mockNewRelicApp{}
	txn := &mockNewRelicTxn{}

	nr.On("StartTransaction", "/health", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		txn.w = args.Get(1).(http.ResponseWriter)
	}).Return(txn)
	txn.On("Header")
	txn.On("AddAttribute", "X-Trace-Id", mock.Anything).Return(nil)
	txn.On("SetName", "/health").Return(nil)
	txn.On("End").Return(nil)
	h.On("Health", mock.Anything, mock.Anything).Return()

	rtr := router.NewRouter(h, log, nr, nil, "my-version", "http://example.com")

	s := httptest.NewServer(rtr)
	defer s.Close()

	req, err := http.NewRequest("GET", s.URL+"/health", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	h.AssertExpectations(t)
	nr.AssertExpectations(t)
	txn.AssertExpectations(t)
}

func TestRouter_NotFound(t *testing.T) {
	log := zap.NewExample()

	h := &mockHandler{}
	nr := &mockNewRelicApp{}
	txn := &mockNewRelicTxn{}

	nr.On("StartTransaction", "/abcd", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		txn.w = args.Get(1).(http.ResponseWriter)
	}).Return(txn)

	txn.On("Header")
	txn.On("WriteHeader", 404)
	txn.On("AddAttribute", "X-Trace-Id", mock.Anything).Return(nil)
	txn.On("End").Return(nil)
	h.On("NotFound", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(404)
	}).Return()

	rtr := router.NewRouter(h, log, nr, nil, "my-version", "http://example.com")

	s := httptest.NewServer(rtr)
	defer s.Close()

	req, err := http.NewRequest("GET", s.URL+"/abcd", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 404, resp.StatusCode)

	h.AssertExpectations(t)
	nr.AssertExpectations(t)
	txn.AssertExpectations(t)
}

func TestRouter_NoToken(t *testing.T) {
	log := zap.NewExample()

	h := &mockHandler{}
	nr := &mockNewRelicApp{}
	txn := &mockNewRelicTxn{}

	nr.On("StartTransaction", "/protected/1", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		txn.w = args.Get(1).(http.ResponseWriter)
	}).Return(txn)

	txn.On("Header")
	txn.On("WriteHeader", 401)
	txn.On("AddAttribute", "X-Trace-Id", mock.Anything).Return(nil)
	txn.On("SetName", "/protected").Return(nil)
	txn.On("End").Return(nil)

	// make sure we call the Unauthorized func on this request.
	h.On("Unauthorized", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(401)
	}).Return()

	rtr := router.NewRouter(h, log, nr, nil, "my-version", "http://example.com")

	s := httptest.NewServer(rtr)
	defer s.Close()

	req, err := http.NewRequest("GET", s.URL+"/protected/1", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)

	// All mocks used should assert that they met their expectations.
	h.AssertExpectations(t)
	nr.AssertExpectations(t)
	txn.AssertExpectations(t)
}

func TestRouter_ValidToken(t *testing.T) {
	log := zap.NewExample()

	h := &handler.Handler{}
	nr := &mockNewRelicApp{}
	txn := &mockNewRelicTxn{}

	nr.On("StartTransaction", "/protected/1", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		txn.w = args.Get(1).(http.ResponseWriter)
	}).Return(txn)

	txn.On("Header")
	txn.On("WriteHeader", 200)
	txn.On("Write", mock.Anything)
	txn.On("AddAttribute", "X-Trace-Id", mock.Anything).Return(nil)
	txn.On("AddAttribute", "id", "1").Return(nil)
	txn.On("SetName", "/protected/{id:[0-9]+}").Return(nil)
	txn.On("End").Return(nil)

	signingKey := []byte("my-key")

	auth := jwtauth.New("HS256", signingKey, nil)

	rtr := router.NewRouter(h, log, nr, auth, "my-version", "http://example.com")

	s := httptest.NewServer(rtr)
	defer s.Close()

	req, err := http.NewRequest("GET", s.URL+"/protected/1", nil)
	require.NoError(t, err)

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "myemail@example.com",
	})

	tokenString, err := tok.SignedString(signingKey)
	require.NoError(t, err)

	req.Header.Set("Authorization", "BEARER "+tokenString)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), "\"message\":\"your username is: myemail@example.com; the id you requested is: 1\"")

	// All mocks used should assert that they met their expectations.
	nr.AssertExpectations(t)
	txn.AssertExpectations(t)
}
