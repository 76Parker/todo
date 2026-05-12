package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyPost_Allowed(t *testing.T) {

	// TEST CODE
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusCreated)
	})
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	rec := httptest.NewRecorder()
	OnlyPost(next).ServeHTTP(rec, req)

	// EQUALITY
	assert.Equal(t, true, nextCalled)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestOnlyPost_NotAllowed(t *testing.T) {
	// TEST CODE
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		nextCalled = true
	})
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	OnlyPost(next).ServeHTTP(rec, req)

	// EQUALITY

	assert.Equal(t, false, nextCalled)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	allowHeader := rec.Header().Get("Allow")
	assert.Equal(t, http.MethodPost, allowHeader)
}
