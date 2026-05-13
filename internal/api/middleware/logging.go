package middleware

import (
	"net/http"
	"todo/internal/utils/ctxutils"

	"github.com/76Parker/golib/loglib"
	"github.com/rs/xid"
)

const requestIDHeader = "X-Request-ID"
const maxHeaderLen = 128

func RequestID(log loglib.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" || len(requestID) > maxHeaderLen {
				requestID = xid.New().String()
			}
			reqLog := log.With("request_id", requestID)
			ctx = ctxutils.SetRequestID(ctx, requestID)
			ctx = ctxutils.SetLoggerInContext(ctx, reqLog)

			w.Header().Set(requestIDHeader, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
