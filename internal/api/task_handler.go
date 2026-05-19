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
	Create(ctx context.Context, createCmd task.CreateCommand) (domain.Task, error)
	ReadByID(ctx context.Context, id int64) (domain.Task, error)
	UpdateByID(ctx context.Context, updateCmd task.UpdateCommand) (domain.Task, error)
	AddTagByID(ctx context.Context, addTagCmd task.TagCommand) (domain.Task, error)
	DeleteTagByID(ctx context.Context, deleteTagCmd task.TagCommand) (domain.Task, error)
	QueryTasks(ctx context.Context, queryCmd task.QueryCommand) ([]domain.Task, error)
}

// Create API handler for create Task (POST /v1/tasks)
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

// Read API handler for read Task by ID (GET /v1/tasks/{id})
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

// Update API handler for update Task by ID (PATCH /v1/tasks/{id}) - updates a task's title, category, description, or status
func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	requestID := ctxlib.RequestID(r.Context())
	id := r.PathValue("id")
	taskID, err := validator.ExtractAndValidateTaskID(id)
	if err != nil {
		log.Warn("ExtractAndValidateTaskID failed", "error", err, "id", id)
		sendJSONError(w, apierr.ValidationFailedError(err.Error(), requestID))
		return
	}

	var updateDto dto.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&updateDto); err != nil {
		log.Warn("decode dto.UpdateTask failed", "error", err, "id", id)
		sendJSONError(w, apierr.ValidationFailedError(err.Error(), requestID))
		return
	}
	if err := t.validator.Validate(updateDto); err != nil {
		log.Warn("validate dto.UpdateTask failed", "error", err, "id", id)
		sendJSONError(w, apierr.ValidationFailedError(err.Error(), requestID))
		return
	}

	cmd := task.UpdateCommand{
		ID:          taskID,
		Title:       updateDto.Title,
		Status:      updateDto.Status,
		Description: updateDto.Description,
		Category:    updateDto.Category,
	}
	updatedTask, err := t.taskSvc.UpdateByID(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			log.Warn("task not found", "id", taskID)
			w.Header().Set("X-Request-ID", requestID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrInvalidStatus) {
			log.Warn("invalid status", "status", updateDto.Status, "id", taskID)
			sendJSONError(w, apierr.BadRequestError(err.Error(), requestID))
			return
		}
		if errors.Is(err, domain.ErrNoChanges) {
			log.Warn("no changes in update request", "id", taskID)
			sendJSONError(w, apierr.BadRequestError("no changes in update task request", requestID))
			return
		}
		log.Error("update task failed", "error", err, "id", taskID)
		w.Header().Set("X-Request-ID", requestID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	readDto := fromDomainTaskToReadTaskDTO(updatedTask)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readDto); err != nil {
		log.Warn("encode dto.ReadTaskDTO failed", "error", err)
	}
	log.Info("update task completed", "id", id)
}

// AddTag API handler for adding tags to a task (POST /tasks/{id}/tags)
func (t *TaskHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	id := r.PathValue("id")
	requestID := ctxlib.RequestID(r.Context())
	taskID, err := validator.ExtractAndValidateTaskID(id)
	if err != nil {
		log.Warn("add tags: ExtractAndValidateTaskID failed", "task_id", id, "err", err.Error())
		sendJSONError(w, apierr.BadRequestError("invalid id", requestID))
		return
	}
	var tag dto.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		log.Warn("add tags: decode tags failed", "error", err.Error(), "task_id", taskID)
		sendJSONError(w, apierr.BadRequestError("invalid json body", requestID))
		return
	}
	if err := t.validator.Validate(tag); err != nil {
		log.Warn("validate tag failed", "error", err.Error())
		sendJSONError(w, apierr.BadRequestError(err.Error(), requestID))
		return
	}
	cmd := task.TagCommand{
		ID:  taskID,
		Tag: tag.Tag,
	}
	task, err := t.taskSvc.AddTagByID(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			log.Warn("add tags: task_id not found", "err", err.Error(), "task_id", taskID)
			w.Header().Set("X-Request-ID", requestID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("add tags failed", "error", err.Error(), "task_id", taskID)
		w.Header().Set("X-Request-ID", requestID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	readDTO := fromDomainTaskToReadTaskDTO(task)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readDTO); err != nil {
		log.Warn("encode dto.ReadTaskDTO failed", "error", err)
	}
	log.Info("add tags completed", "task_id", taskID)
}

// DeleteTag API handler for deleting tags from a task (DELETE /tasks/{id}/tags)
func (t *TaskHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	requestID := ctxlib.RequestID(r.Context())
	id := r.PathValue("id")
	taskID, err := validator.ExtractAndValidateTaskID(id)
	if err != nil {
		log.Warn("ExtractAndValidateTaskID failed", "task_id", id, "err", err.Error())
		sendJSONError(w, apierr.BadRequestError("invalid id", requestID))
		return
	}

	var tag dto.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		log.Warn("decode tag failed", "error", err.Error())
		sendJSONError(w, apierr.BadRequestError("invalid json body", requestID))
		return
	}

	if err := t.validator.Validate(tag); err != nil {
		log.Warn("validate tag failed", "error", err.Error())
		sendJSONError(w, apierr.BadRequestError(err.Error(), requestID))
		return
	}
	deleteCmd := task.TagCommand{
		ID:  taskID,
		Tag: tag.Tag,
	}

	task, err := t.taskSvc.DeleteTagByID(r.Context(), deleteCmd)
	if err != nil {
		if errors.Is(err, domain.ErrTagOrTaskNotFound) {
			log.Warn("delete tags: tag or task not found", "err", err.Error(), "task_id", taskID)
			sendJSONError(w, apierr.NotFoundError("tag or task not found", requestID))
			return
		}
		log.Error("delete tags failed", "error", err.Error(), "task_id", taskID)
		w.Header().Set("X-Request-ID", requestID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	readDTO := fromDomainTaskToReadTaskDTO(task)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readDTO); err != nil {
		log.Warn("delete tags: encode dto.ReadTaskDTO failed", "error", err, "task_id", taskID)
	}
	log.Info("delete tags completed", "task_id", taskID)
}

// QueryTasks API handler for searching tasks by filter. Supports only title (GET /tasks?title=)
func (t *TaskHandler) QueryTasks(w http.ResponseWriter, r *http.Request) {
	log := ctxlib.GetLoggerFromContext(r.Context())
	requestID := ctxlib.RequestID(r.Context())
	query := r.URL.Query().Get("title")
	tasks, err := t.taskSvc.QueryTasks(r.Context(), task.QueryCommand{Title: query})
	if err != nil {
		log.Error("search tasks failed", "error", err.Error(), "query", query)
		w.Header().Set("X-Request-ID", requestID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	readDTOs := make([]dto.ReadTask, 0, len(tasks))
	for _, task := range tasks {
		readDTOs = append(readDTOs, fromDomainTaskToReadTaskDTO(task))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readDTOs); err != nil {
		log.Warn("search tasks: encode dto.ReadTaskDTO failed", "error", err, "query", query)
	}
	log.Info("search tasks completed", "query", query)
}
