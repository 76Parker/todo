package routes

import (
	"net/http"
	"todo/internal/api"
)

type apiVersion string

const (
	apiV1 apiVersion = "/v1"
)

type MW func(next http.Handler) http.Handler // Middleware

type Handlers struct {
	Tasks *api.TaskHandler
}
