package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {

	tests := []struct {
		name    string
		input   Task
		wantErr bool
	}{
		{
			name: "ValidTask_1",
			input: Task{
				title:       " title_test ",
				description: " description_test ",
				category:    "category_test",
				status:      statusOpen,
				tags:        []string{" tag1", "tag2 "},
			},
			wantErr: false,
		},
		{
			name: "ValidTask_2",
			input: Task{
				title: "title_test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testTask := NewTask(tt.input.title, tt.input.description, tt.input.category, tt.input.tags, tt.input.status)

			assert.Equal(t, normalizeString(tt.input.title), testTask.title)
			assert.Equal(t, normalizeString(tt.input.description), testTask.description)
			assert.Equal(t, normalizeString(tt.input.category), testTask.category)
			assert.Equal(t, tt.input.status, testTask.status)

			inputNormalizedTags := make([]string, len(testTask.tags))
			for i, tag := range tt.input.tags {
				inputNormalizedTags[i] = normalizeString(tag)
			}
			assert.Equal(t, inputNormalizedTags, testTask.tags)
		})
	}
}

func TestReconstituteTask(t *testing.T) {
	tests := []struct {
		name    string
		input   Task
		wantErr bool
	}{
		{
			name: "ValidReconstitutedTask_1",
			input: Task{
				id:          2,
				title:       " title_test ",
				description: " description_test ",
				status:      " status_test",
				createdAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "ValidReconstitutedTask_2",
			input: Task{
				id:          2,
				title:       "",
				description: " description_test ",
				status:      "open",
				category:    "category_test",
				tags:        []string{" tag1", "tag2 "},
				createdAt:   time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testTask := ReconstituteTask(tt.input.id, tt.input.title, tt.input.description, string(tt.input.Status()), tt.input.category, tt.input.tags, tt.input.createdAt)
			assert.Equal(t, tt.input.status, testTask.Status())
			assert.Equal(t, tt.input.category, testTask.Category())
			assert.Equal(t, tt.input.tags, testTask.Tags())
			assert.Equal(t, tt.input.createdAt, testTask.CreatedAt())
			assert.Equal(t, tt.input.id, testTask.ID())
		})
	}
}

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "ValidStatus_1",
			input:   "open",
			wantErr: false,
		},
		{
			name:    "ValidStatus_2",
			input:   "done",
			wantErr: false,
		},
		{
			name:    "ValidStatus_3",
			input:   "closed",
			wantErr: false,
		},
		{
			name:    "ValidStatus_4",
			input:   "in_progress",
			wantErr: false,
		},
		{
			name:    "InvalidStatus_1",
			input:   "",
			wantErr: true,
		},
		{
			name:    "InvalidStatus_2",
			input:   "pending",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewStatus(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, string(status))
			}
		})
	}
}
