package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo/internal/api/middleware"
	"todo/internal/entities/domain"
	"todo/internal/usecase/task"

	"github.com/76Parker/golib/loglib"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskHandler_Create(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		wantServiceCall  bool
		wantValidateCall bool
		serviceReturn    domain.Task
		wantCmd          task.CreateCommand
		validatorErr     error
		serviceErr       error
		expectedCode     int
	}{
		{
			name:             "StatusCreated_1",
			body:             `{"title":"test", "status":"open"}`,
			serviceReturn:    domain.NewTask("test", "", "", nil, "open"),
			wantCmd:          task.CreateCommand{Title: "test", Status: "open"},
			wantServiceCall:  true,
			wantValidateCall: true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusCreated,
		},
		{
			name:             "StatusCreated_2",
			body:             `{"title":"test", "status": "open", "description":"test", "category":"test"}`,
			serviceReturn:    domain.NewTask("test", "test", "test", nil, "open"),
			wantCmd:          task.CreateCommand{Title: "test", Description: "test", Category: "test", Tags: nil, Status: "open"},
			wantValidateCall: true,
			wantServiceCall:  true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusCreated,
		},
		{
			name:             "StatusCreated_3",
			body:             `{"title":"test", "status": "done", "tags": ["test", "test"]}`,
			serviceReturn:    domain.NewTask("test", "", "", []string{"test", "test"}, "done"),
			wantCmd:          task.CreateCommand{Title: "test", Tags: []string{"test", "test"}, Status: "done"},
			wantValidateCall: true,
			wantServiceCall:  true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusCreated,
		},
		{
			name:             "BadRequest_1",
			body:             `{"title":"tes`,
			wantValidateCall: false,
			wantServiceCall:  false,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
		{
			name:             "BadRequest_2",
			body:             `{"title":"  ", "description":"test", "category":"test"}`,
			wantValidateCall: true,
			wantServiceCall:  false,
			validatorErr:     errors.New("validation error"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			taskSvc := NewMockTaskService(ctrl)

			if tt.wantServiceCall {
				taskSvc.EXPECT().Create(gomock.Any(), tt.wantCmd).Return(tt.serviceReturn, nil)
			}

			mockLogger := loglib.NewMockLogger()
			taskHandler := NewTaskHandler(taskSvc)

			req := httptest.NewRequest(http.MethodPost, "/v1/tasks", strings.NewReader(tt.body))
			ctx := req.Context()
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			handler := middleware.RequestID(mockLogger)(
				http.HandlerFunc(taskHandler.Create),
			)
			handler.ServeHTTP(rec, req)
			resp := rec.Result()
			defer func() {
				_ = req.Body.Close()
			}()
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestTaskHandler_Read(t *testing.T) {
	tests := []struct {
		name             string
		taskID           string
		wantServiceCall  bool
		wantValidateCall bool
		serviceErr       error
		validatorErr     error
		expectedCode     int
		serviceReturn    domain.Task
	}{
		{
			name:             "StatusOK_1",
			taskID:           "10",
			wantServiceCall:  true,
			wantValidateCall: true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusOK,
			serviceReturn:    domain.ReconstituteTask(10, "", "", "", "", []string{}, time.Time{}),
		},
		{
			name:             "StatusOK_2",
			taskID:           "1",
			serviceReturn:    domain.ReconstituteTask(1, "", "", "", "", []string{}, time.Time{}),
			wantServiceCall:  true,
			wantValidateCall: true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusOK,
		}, {
			name:             "BadRequest_1",
			taskID:           "bad_request_id",
			wantValidateCall: true,
			wantServiceCall:  false,
			validatorErr:     errors.New("validation error"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
		{
			name:             "BadRequest_2",
			taskID:           "",
			wantValidateCall: true,
			wantServiceCall:  false,
			validatorErr:     errors.New("validation error"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			taskSvc := NewMockTaskService(ctrl)
			if tt.wantServiceCall {
				taskSvc.EXPECT().ReadByID(gomock.Any(), gomock.Any()).Return(tt.serviceReturn, nil)
			}
			mockLogger := loglib.NewMockLogger()
			taskHandler := NewTaskHandler(taskSvc)
			req := httptest.NewRequest(http.MethodGet, "/v1/tasks/"+tt.taskID, nil)
			req.SetPathValue("id", tt.taskID)
			rec := httptest.NewRecorder()
			handler := middleware.RequestID(mockLogger)(
				http.HandlerFunc(taskHandler.Read),
			)
			handler.ServeHTTP(rec, req)
			resp := rec.Result()
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name            string
		taskID          string
		wantValidateCall bool
		wantServiceCall  bool
		body            string
		serviceReturn   domain.Task
		updatedField    string
		validatorErr     error
		serviceErr       error
		expectedCode     int
	}{
		{
			name: "UpdateTask_Success",
			taskID:          "10",
			body:            `{"title": "Updated Title"}`,
			serviceReturn:   domain.NewTask("Updated Title", "test", "test", []string{"123"},domain.Status("open")),
			updatedField:    "title",
			wantValidateCall: true,
			wantServiceCall: true,
			validatorErr:     nil,
			serviceErr:       nil,
			expectedCode:     http.StatusOK,
		},
		{
			name: "UpdateTask_ValidateError",
			taskID:          "10",
			body:            `{"title":""}`,
			wantValidateCall: true,
			wantServiceCall: false,
			validatorErr:     errors.New("title is required"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
		{
			name: "UpdateTask_NotFound",
			taskID:          "11",
			body:            `{"title": "Updated Title"}`,
			wantValidateCall: true,
			wantServiceCall: true,
			validatorErr:     nil,
			serviceErr:       domain.ErrTaskNotFound,
			expectedCode:     http.StatusNotFound,
		},
		{
			name: "UpdateTask_InvalidJSON",
			taskID:          "10",
			body:            `{"title": "Updated Tit}`,
			wantValidateCall: false,
			wantServiceCall:  false,
			validatorErr:     errors.New("invalid JSON"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},
		{
			name: "UpdateTask_InvalidStatus",
			taskID:          "10",
			body:            `{"status": "fsdfsdf"}`,
			wantValidateCall: true,
			wantServiceCall:  true,
			validatorErr:     nil,
			serviceErr:       domain.ErrInvalidStatus,
			expectedCode:     http.StatusBadRequest,
		}, {
			name: "UpdateTask_BlankTitle",
			taskID:          "10",
			body:            `{"title": ""}`,
			wantValidateCall: true,
			wantServiceCall:  false,
			validatorErr:     errors.New("title is required"),
			serviceErr:       nil,
			expectedCode:     http.StatusBadRequest,
		},{
			name: "UpdateTask_ClearDescription",
			taskID:          "10",
			body:            `{"description": ""}`,
			wantValidateCall: true,
			wantServiceCall:  true,
			validatorErr:    nil,
			serviceErr:       nil,
			expectedCode:     http.StatusOK,
		}, {
			name: "UpdateTask_BlankStatus",
			taskID:          "10",
			body:            `{"status": ""}`,
			wantValidateCall: true,
			wantServiceCall:  true,
			validatorErr:   nil,
			serviceErr:       domain.ErrInvalidStatus,
			expectedCode:     http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := NewMockTaskService(ctrl)
			if tt.wantServiceCall {
				if tt.serviceErr != nil {
					mockSvc.EXPECT().UpdateByID(gomock.Any(), gomock.Any()).Return(domain.Task{},tt.serviceErr)
				} else {
					mockSvc.EXPECT().UpdateByID(gomock.Any(), gomock.Any()).Return(domain.Task{}, nil)
				}
			}

			mockLogger := loglib.NewMockLogger()
			taskHandler := NewTaskHandler(mockSvc)
			wrappedTaskHandler := middleware.RequestID(mockLogger) (
				http.HandlerFunc(taskHandler.Update),
			)
			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, strings.NewReader(tt.body))
			req.SetPathValue("id", tt.taskID)
			rec := httptest.NewRecorder()
			wrappedTaskHandler.ServeHTTP(rec,req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
