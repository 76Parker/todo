package task

import (
	"context"
	"errors"
	"testing"
	"time"
	"todo/internal/entities/domain"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func testReconstituteTaskByID(id int64) domain.Task {
	return domain.ReconstituteTask(id, "", "", "", "", []string{}, time.Time{})
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name         string
		cmd          CreateCommand
		wantErr      bool
		wantRepoCall bool
	}{
		{
			name: "ValidCreateCommand_1",
			cmd: CreateCommand{
				Title:  "test",
				Status: "done",
			},
			wantErr:      false,
			wantRepoCall: true,
		},
		{
			name: "ValidCreateCommand_2",
			cmd: CreateCommand{
				Title:       "test",
				Description: "test",
				Status:      "in_progress",
				Category:    "test",
				Tags:        []string{"test"},
			},
			wantErr:      false,
			wantRepoCall: true,
		},
		{
			name: "InvalidCreateCommand_1",
			cmd: CreateCommand{
				Title:       "test",
				Description: "test",
				Status:      "invalid_status",
			},
			wantErr:      true,
			wantRepoCall: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepository(ctrl)
			if tt.wantErr && tt.wantRepoCall {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.Task{}, errors.New("error"))
			} else if !tt.wantErr && tt.wantRepoCall {
				status, _ := domain.NewStatus(tt.cmd.Status)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.NewTask(tt.cmd.Title, tt.cmd.Description, tt.cmd.Category, tt.cmd.Tags, status), nil)
			}
			service := NewService(mockRepo)
			gotTask, err := service.Create(context.Background(), tt.cmd)
			if err != nil && tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.cmd.Title, gotTask.Title())
				assert.Equal(t, tt.cmd.Description, gotTask.Description())
				assert.Equal(t, tt.cmd.Category, gotTask.Category())
				assert.Equal(t, tt.cmd.Tags, gotTask.Tags())
			}
		})
	}
}

func TestService_ReadByID(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		wantRepoCall bool
		wantErr      bool
		repoErr      error
		repoReturn   domain.Task
	}{
		{
			name:         "ValidReadByID_1",
			id:           1,
			wantRepoCall: true,
			wantErr:      false,
			repoReturn:   testReconstituteTaskByID(int64(1)),
		},
		{
			name:         "InvalidReadByID_1",
			id:           0,
			wantRepoCall: true,
			wantErr:      true,
			repoErr:      errors.New("task not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockRepository(ctrl)
			if tt.wantErr && tt.wantRepoCall {
				mockRepo.EXPECT().ReadByID(gomock.Any(), gomock.Any()).Return(domain.Task{}, tt.repoErr)
			} else if tt.wantRepoCall {
				mockRepo.EXPECT().ReadByID(gomock.Any(), gomock.Any()).Return(tt.repoReturn, tt.repoErr)
			}
			service := NewService(mockRepo)
			gotTask, err := service.ReadByID(context.Background(), tt.id)
			if err != nil && tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, gotTask.ID())
			}
		})
	}
}
