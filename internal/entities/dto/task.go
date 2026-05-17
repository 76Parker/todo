package dto

import "time"

// CreateTask dto for create task request
type CreateTask struct {
	Title       string    `json:"title" validate:"min=1,max=200,required,notblank"`
	Status      string    `json:"status" validate:"min=1,max=20,required,notblank"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1,max=500,notblank"`
	Category    *string   `json:"category,omitempty" validate:"omitempty,min=1,max=500,notblank"`
	Tags        *[]string `json:"tags,omitempty" validate:"omitempty,min=1,max=500,dive,required,min=1,max=30,notblank"`
}

// ReadTask dto for read task
type ReadTask struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

// Type method for implement DTO interface
func (v CreateTask) Type() string {
	return "CreateTask"
}

// Type method for implement DTO interface
func (v ReadTask) Type() string {
	return "ReadTask"
}
