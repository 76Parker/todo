package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"todo/internal/api/apierr"
	"todo/internal/api/validator"
	"todo/internal/entities/domain"
	"todo/internal/entities/dto"
	"todo/internal/usecase/task"

	"github.com/76Parker/golib/ctxlib"
	"github.com/c2h5oh/datasize"
)

const (
	maxBodySize = datasize.MB * 5
)

func sendJSONError(w http.ResponseWriter, err apierr.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", err.RequestID)
	w.WriteHeader(err.Status)
	_ = json.NewEncoder(w).Encode(err)
}

// TaskHandler handle request related to Task
type TaskHandler struct {
	taskSvc   TaskService
	validator DtoValidator
}

// NewTaskHandler is a constructor for TaskHandler
func NewTaskHandler(taskSvc TaskService) *TaskHandler {
	return &TaskHandler{
		taskSvc:   taskSvc,
		validator: validator.MustDtoValidator(),
	}
}

// DtoValidator validate all input DTO's which implemented dto.DTO interface
type DtoValidator interface {
	Validate(dto dto.DTO) error
}

// TaskService create Task and check domain rules
type TaskService interface {
	Create(ctx context.Context, cmd task.CreateCommand) (domain.Task, error)
	ReadByID(ctx context.Context, id int64) (domain.Task, error)
}

// Create external HTTP API for create Task
func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	requestID := ctxlib.RequestID(r.Context())

	var taskDto dto.CreateTask

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))
	defer func() {
		_ = r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&taskDto); err != nil {
		log.Warn("decode CreateTaskDTO failed", "error", err)
		sendJSONError(w, apierr.BadRequestError("invalid json body", requestID))
		return
	}

	if err := t.validator.Validate(taskDto); err != nil {
		log.Warn("validate CreateTaskDTO failed", "error", err)
		sendJSONError(w, apierr.ValidationFailedError(err.Error(), requestID))
		return
	}

	cmd := createTaskDtoToCommand(taskDto)
	domainTask, err := t.taskSvc.Create(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidStatus) {
			log.Warn("invalid status", "error", err)
			sendJSONError(w, apierr.BadRequestError("invalid status", requestID))
			return
		}
		log.Error("taskSvc.Create failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("new task created", "id", domainTask.ID())

	readDto := fromDomainTaskToReadTaskDTO(domainTask)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(readDto); err != nil {
		log.Warn("encode dto.ReadTask failed", "error", err)
	}
}

func (t *TaskHandler) Read(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	requestID := ctxlib.RequestID(r.Context())
	id := r.PathValue("id")

	taskID, err := validator.ExtractAndValidateTaskID(id)
	if err != nil {
		log.Warn("ExtractAndValidateTaskID failed", "error", err, "id", id)
		sendJSONError(w, apierr.ValidationFailedError(err.Error(), requestID))
		return
	}

	domainTask, err := t.taskSvc.ReadByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			log.Warn("task not found", "id", id)
			w.Header().Set("X-Request-ID", requestID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("read task failed", "error", err, "id", id)
		w.Header().Set("X-Request-ID", requestID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	readDto := fromDomainTaskToReadTaskDTO(domainTask)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readDto); err != nil {
		log.Warn("encode dto.ReadTaskDTO failed", "error", err)
	}
	log.Info("read task completed", "id", id)
}
