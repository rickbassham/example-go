package handler

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

// SimpleResponse is used to send a meaningful message back to the caller, with
// trace id to debug later.
type SimpleResponse struct {
	XMLName xml.Name `json:"-" xml:"simpleResponse"`
	TraceID string   `json:"trace_id" xml:"traceId,attr"`
	Message string   `json:"message" xml:",innerxml"`
}

func writeResponse(r *http.Request, w http.ResponseWriter, status int, resp interface{}) {
	accept := r.Header.Get("Accept")

	if strings.HasPrefix(accept, "text/xml") {
		writeXMLResponse(r.Context(), w, status, resp)
	} else {
		writeJSONResponse(r.Context(), w, status, resp)
	}
}

func writeXMLResponse(ctx context.Context, w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	err := xml.NewEncoder(w).Encode(resp)
	if err != nil {
		middleware.GetLogger(ctx).Error("error writing response", zap.Error(err))
		newrelic.FromContext(ctx).NoticeError(err) // nolint
	}
}

func writeJSONResponse(ctx context.Context, w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		middleware.GetLogger(ctx).Error("error writing response", zap.Error(err))
		newrelic.FromContext(ctx).NoticeError(err) // nolint
	}
}
