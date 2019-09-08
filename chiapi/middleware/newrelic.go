package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	newrelic "github.com/newrelic/go-agent"
)

// NewRelicChiRouter will start a newrelic transaction for the request. It will
// name the transaction to match the chi route. It will also add custom attributes
// to the transaction, such as the trace id, any url parameters, and query string
// parameters. These attributes make it easy to tie your New Relic traces to logs.
// This middleware can be used in place of using newrelic.WrapHandleFunc
// everywhere.
func NewRelicChiRouter(nr newrelic.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			txn := nr.StartTransaction(r.URL.Path, w, r)

			r = r.WithContext(newrelic.NewContext(r.Context(), txn))

			next.ServeHTTP(txn, r)

			rctx := chi.RouteContext(r.Context())

			txn.AddAttribute("X-Trace-Id", GetTraceID(r.Context())) // nolint

			for i := range rctx.URLParams.Keys {
				if rctx.URLParams.Keys[i] == "*" {
					continue
				}

				txn.AddAttribute(rctx.URLParams.Keys[i], rctx.URLParams.Values[i]) // nolint
			}

			for k, v := range r.URL.Query() {
				txn.AddAttribute(k, v[0]) // nolint
			}

			routePattern := strings.Join(rctx.RoutePatterns, "")
			routePattern = strings.Replace(routePattern, "/*", "", -1)

			if routePattern != "" {
				txn.SetName(routePattern) // nolint
			}

			txn.End() // nolint
		}

		return http.HandlerFunc(fn)
	}
}
