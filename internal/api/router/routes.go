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

type Handlers struct {
	Tasks *api.TaskHandler
}

func New(handlers Handlers, log loglib.Logger) http.Handler {
	mux := http.NewServeMux()
	registerSwaggerRoutes(mux)
	registerV1TaskRoutes(mux, handlers.Tasks, log)
	return mux
}

func registerV1TaskRoutes(mux *http.ServeMux, taskHandler *api.TaskHandler, log loglib.Logger) {
	requestID := middleware.RequestID(log)
	createHandler := http.HandlerFunc(taskHandler.Create)
	mux.Handle("POST "+apiV1+"/tasks", middleware.Recover(requestID(createHandler)))
}
