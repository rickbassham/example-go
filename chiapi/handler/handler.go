// Package handler implements the request handlers for the API.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

// Handler has all the functions needed to serve our api.
type Handler struct {
}

// Health always returns a 200 response.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusOK, &SimpleResponse{
		TraceID: middleware.GetTraceID(r.Context()),
		Message: "OK",
	})
}

// Protected is a protected endpoint example.
func (h *Handler) Protected(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.GetUser(ctx)
	traceID := middleware.GetTraceID(ctx)
	id := routeParamInt(ctx, "id")

	writeResponse(r, w, http.StatusOK, &SimpleResponse{
		TraceID: traceID,
		Message: fmt.Sprintf("your username is: %s; the id you requested is: %d", user, id),
	})
}

// Unauthorized is called when the request is not authorized.
func (h *Handler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusUnauthorized, &SimpleResponse{
		TraceID: middleware.GetTraceID(r.Context()),
		Message: "unauthorized",
	})
}

// NotFound is called when the request is for an unknown resource.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeResponse(r, w, http.StatusNotFound, &SimpleResponse{
		TraceID: middleware.GetTraceID(r.Context()),
		Message: "not found",
	})
}

func routeParamInt(ctx context.Context, name string) int {
	// this func should only be called for params that are guaranteed to be ints.
	val, _ := strconv.Atoi(chi.RouteContext(ctx).URLParam("id")) // nolint
	return val
}
