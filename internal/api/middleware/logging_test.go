package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/internal/utils/ctxutils"

	"github.com/76Parker/golib/loglib"
	"github.com/stretchr/testify/assert"
)

func TestLog_RequestID_Auto(t *testing.T) {
	var gotRequestID string
	var logger loglib.Logger

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = ctxutils.RequestID(r.Context())
		logger = ctxutils.GetLoggerFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mockLogger := loglib.NewMockLogger()

	logMW := RequestID(mockLogger)

	logMW(next).ServeHTTP(rec, req)

	assert.NotEmpty(t, gotRequestID)
	assert.NotNil(t, logger)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLog_RequestID_Header(t *testing.T) {
	var gotRequestID string
	var logger loglib.Logger
	testRequestID := "test-request-id"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = ctxutils.RequestID(r.Context())
		logger = ctxutils.GetLoggerFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(requestIDHeader, testRequestID)
	rec := httptest.NewRecorder()

	mockLogger := loglib.NewMockLogger()

	logMW := RequestID(mockLogger)

	logMW(next).ServeHTTP(rec, req)

	assert.NotNil(t, logger)
	assert.NotEmpty(t, gotRequestID)
	assert.Equal(t, testRequestID, gotRequestID)
	assert.Equal(t, http.StatusOK, rec.Code)
}
