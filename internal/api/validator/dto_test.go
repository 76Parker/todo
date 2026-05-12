package validator

import (
	"testing"
	"todo/internal/entities/dto"
)

func TestDTO_Validate(t *testing.T) {

	v := NewDtoValidator()
	tests := []struct {
		name    string
		input   dto.CreateTask
		wantErr bool
	}{
		{
			name:    "valid dto",
			input:   dto.CreateTask{Title: "test"},
			wantErr: false,
		},
		{
			name:    "valid dto 2",
			input:   dto.CreateTask{Title: "test", Description: new("description")},
			wantErr: false,
		},
		{
			name:    "valid dto 3",
			input:   dto.CreateTask{Title: "test", Description: new("description"), Category: new("category")},
			wantErr: false,
		},
		{
			name:    "valid dto 4",
			input:   dto.CreateTask{Title: "test", Tags: new([]string{"tag1", "tag2"}), Description: new("description"), Category: new("category")},
			wantErr: false,
		},
		{
			name:    "invalid dto",
			input:   dto.CreateTask{Title: ""},
			wantErr: true,
		},
		{
			name:    "invalid dto 2",
			input:   dto.CreateTask{Title: "title", Description: new("")},
			wantErr: true,
		},
		{
			name:    "invalid dto 3",
			input:   dto.CreateTask{Title: "title", Description: new("description"), Category: new("")},
			wantErr: true,
		},
		{
			name:    "invalid dto 4",
			input:   dto.CreateTask{Title: "title", Tags: new([]string{""})},
			wantErr: true,
		},
		{
			name:    "invalid dto 5",
			input:   dto.CreateTask{Title: "title", Tags: new([]string{"fdsfsdfsdfsdfdsfsfsdfddsfdsfsdfsdfsdfsdfsfsfsfsfsfsdfsdfsdfsfsfsfdsfdsfdsfdsfsdfsdf"})},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
