package dto

type CreateTask struct {
	Title       string    `json:"title" validate:"min=1,max=200"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	Category    *string   `json:"category,omitempty" validate:"omitempty,min=1,max=500"`
	Tags        *[]string `json:"tags,omitempty" validate:"omitempty,min=1,max=500,dive,required,min=1,max=30"`
}

func (v CreateTask) Type() string {
	return "CreateTask"
}
