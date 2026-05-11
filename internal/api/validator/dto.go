package validator

import (
	"todo/internal/entities/dto"

	"github.com/go-playground/validator"
)

type DTO struct {
	v *validator.Validate
}

func (v *DTO) Validate(dto dto.DTO) error {
	return nil
}
