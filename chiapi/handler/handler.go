// Package handler implements the request handlers for the API.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/pkg/identity"
	"github.com/rickbassham/example-go/pkg/logging"
	"github.com/rickbassham/example-go/pkg/tracing"
)

// Cache defines the funcs needed to get stuff from redis.
type Cache interface {
	GetValue(ctx context.Context) (string, error)
}

// Handler has all the functions needed to serve our api.
type Handler struct {
	cache Cache
}

// New creates a new handler.
func New(c Cache) *Handler {
	return &Handler{
		cache: c,
	}
}

// Health always returns a 200 response.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusOK, &SimpleResponse{
		TraceID: tracing.FromContext(r.Context()),
		Message: "OK",
	})
}

// Cached gets a value from redis and writes it to the response.
func (h *Handler) Cached(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := logging.FromContext(ctx)
	l.Info("test")

	traceID := tracing.FromContext(ctx)

	v, err := h.cache.GetValue(ctx)
	if err != nil {
		txn := newrelic.FromContext(ctx)
		txn.NoticeError(err) // nolint

		//l := middleware.GetLogger(ctx)
		l.Error("error getting value from cache", zap.Error(err))

		writeResponse(r, w, http.StatusInternalServerError, &SimpleResponse{
			TraceID: traceID,
			Message: "error getting value from cache",
		})

		return
	}

	writeResponse(r, w, http.StatusOK, &SimpleResponse{
		TraceID: traceID,
		Message: v,
	})
}

// Protected is a protected endpoint example.
func (h *Handler) Protected(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := identity.FromContext(ctx)
	traceID := tracing.FromContext(ctx)
	id := routeParamInt(ctx, "id")

	writeResponse(r, w, http.StatusOK, &SimpleResponse{
		TraceID: traceID,
		Message: fmt.Sprintf("your username is: %s; the id you requested is: %d", user, id),
	})
}

// Unauthorized is called when the request is not authorized.
func (h *Handler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusUnauthorized, &SimpleResponse{
		TraceID: tracing.FromContext(r.Context()),
		Message: "unauthorized",
	})
}

// NotFound is called when the request is for an unknown resource.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusNotFound, &SimpleResponse{
		TraceID: tracing.FromContext(r.Context()),
		Message: "not found",
	})
}

func routeParamInt(ctx context.Context, name string) int {
	// this func should only be called for params that are guaranteed to be ints.
	val, _ := strconv.Atoi(chi.RouteContext(ctx).URLParam("id")) // nolint
	return val
}
