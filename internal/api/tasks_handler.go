package api

import (
	"net/http"
)

type TaskHandler struct{}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {}
