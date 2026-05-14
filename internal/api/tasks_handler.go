package api

import (
	"net/http"
	"todo/internal/utils/ctxutils"
)

type TaskHandler struct{}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	_ = ctxutils.GetLoggerFromContext(r.Context())
}
