package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/internal/api"

	"github.com/76Parker/golib/loglib"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerUIRoute(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = req.Body.Close()
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/html")
	assert.Contains(t, rec.Body.String(), "/swagger/openapi.yaml")
}

func TestOpenAPIRoute(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
		_ = req.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/yaml")
	assert.Contains(t, rec.Body.String(), "openapi: 3.0.3")
	assert.Contains(t, rec.Body.String(), "/v1/tasks")
	assert.Contains(t, rec.Body.String(), "/v1/tasks/{id}")
	assert.Contains(t, rec.Body.String(), "patch:")
}
