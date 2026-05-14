package dto

type CreateTask struct {
	Title       string    `json:"title" validate:"min=1,max=200,required,notblank"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1,max=500,notblank"`
	Category    *string   `json:"category,omitempty" validate:"omitempty,min=1,max=500,notblank"`
	Tags        *[]string `json:"tags,omitempty" validate:"omitempty,min=1,max=500,dive,required,min=1,max=30,notblank"`
}

func (v CreateTask) Type() string {
	return "CreateTask"
}
