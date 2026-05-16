package validator

import (
	"testing"
	"todo/internal/entities/dto"

	"github.com/stretchr/testify/assert"
)

func TestDTO_Validate(t *testing.T) {

	v := MustDtoValidator()
	tests := []struct {
		name    string
		input   dto.CreateTask
		wantErr bool
	}{
		{
			name:    "ValidDTO_1",
			input:   dto.CreateTask{Title: "test"},
			wantErr: false,
		},
		{
			name:    "ValidDTO_2",
			input:   dto.CreateTask{Title: "test", Description: new("description")},
			wantErr: false,
		},
		{
			name:    "ValidDTO_3",
			input:   dto.CreateTask{Title: "test", Description: new("description"), Category: new("category")},
			wantErr: false,
		},
		{
			name:    "ValidDTO_4",
			input:   dto.CreateTask{Title: "test", Tags: new([]string{"tag1", "tag2"}), Description: new("description"), Category: new("category")},
			wantErr: false,
		},
		{
			name:    "ValidDTO_5",
			input:   dto.CreateTask{Title: ""},
			wantErr: true,
		},
		{
			name:    "InvalidDTO_1",
			input:   dto.CreateTask{Title: "title", Description: new("")},
			wantErr: true,
		},
		{
			name:    "InvalidDTO_2",
			input:   dto.CreateTask{Title: "title", Description: new("description"), Category: new("")},
			wantErr: true,
		},
		{
			name:    "InvalidDTO_3",
			input:   dto.CreateTask{Title: "title", Tags: new([]string{""})},
			wantErr: true,
		},
		{
			name:    "InvalidDTO_4",
			input:   dto.CreateTask{Title: "title", Tags: new([]string{"fdsfsdfsdfsdfdsfsfsdfddsfdsfsdfsdfsdfsdfsfsfsfsfsfsdfsdfsdfsfsfsfdsfdsfdsfdsfsdfsdf"})},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if tt.wantErr {
				t.Log(err.Error())
				assert.Error(t, err)
			}
		})
	}
}
