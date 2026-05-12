package middleware

import (
	"net/http"
	"todo/internal/utils/ctxutils"

	"github.com/76Parker/golib/loglib"
	"github.com/rs/xid"
)

type Log struct {
	logger loglib.Logger
}

const requestIDHeader = "X-Request-ID"
const maxHeaderLen = 128

func NewLogMW(l loglib.Logger) *Log {
	return &Log{
		logger: l,
	}
}

func (l *Log) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" || len(requestID) > maxHeaderLen {
			requestID = xid.New().String()
		}
		ctx = ctxutils.SetRequestID(ctx, requestID)

		w.Header().Set(requestIDHeader, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
