// Package router uses for register HTTP-routs
package router

import (
	"net/http"
	"todo/internal/api"
	"todo/internal/api/middleware"

	"github.com/76Parker/golib/loglib"
)

const (
	apiV1 string = "/v1"
)

// Handlers it's a container for handlers from package api
type Handlers struct {
	Tasks *api.TaskHandler
}

// New constructor for server http.Handler
func New(handlers Handlers, log loglib.Logger) http.Handler {
	mux := http.NewServeMux()
	registerSwaggerRoutes(mux)
	registerV1TaskRoutes(mux, handlers.Tasks, log)
	return mux
}

func registerV1TaskRoutes(mux *http.ServeMux, taskHandler *api.TaskHandler, log loglib.Logger) {
	requestID := middleware.RequestID(log)
	createHandler := http.HandlerFunc(taskHandler.Create)
	readHandler := http.HandlerFunc(taskHandler.Read)
	updateHandler := http.HandlerFunc(taskHandler.Update)
	mux.Handle("POST "+apiV1+"/tasks", middleware.Recover(requestID(createHandler)))
	mux.Handle("GET "+apiV1+"/tasks/{id}", middleware.Recover(requestID(readHandler)))
	mux.Handle("PATCH "+apiV1+"/tasks/{id}", middleware.Recover(requestID(updateHandler)))
}
