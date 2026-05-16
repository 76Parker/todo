package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
