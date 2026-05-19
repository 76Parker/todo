package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/internal/api"

	"github.com/76Parker/golib/loglib"
	"github.com/stretchr/testify/assert"
)

func TestTaskReadRoutePattern(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks/15", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestTaskQueryRoutePattern(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks?title=test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestTaskUpdateRoutePattern(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/15", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestTaskAddTagsRoutePattern(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodPost, "/v1/tasks/15/tags", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestTaskDeleteTagsRoutePattern(t *testing.T) {
	r := New(Handlers{Tasks: api.NewTaskHandler(nil)}, loglib.NewMockLogger())

	req := httptest.NewRequest(http.MethodDelete, "/v1/tasks/15/tags", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}
