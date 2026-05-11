package api

import (
	"net/http"

	"github.com/76Parker/golib/loglib"
)

type TasksHandler struct {
	log loglib.Logger
}

func (t *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {}
