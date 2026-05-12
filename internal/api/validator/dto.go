package validator

import (
	"errors"
	"todo/internal/api/apierr"
	"todo/internal/entities/dto"

	"github.com/go-playground/validator"
)

var (
	requiredTag = "required"
	minTag      = "min"
	maxTag      = "max"

	errRequired = errors.New("required field is missing")
	errMin      = errors.New("field is too small")
	errMax      = errors.New("field is too large")
)

type DTO struct {
	v *validator.Validate
}

func NewDtoValidator() *DTO {
	return &DTO{
		v: validator.New(),
	}
}
func (v *DTO) Validate(dto dto.DTO) error {
	dtoName := dto.Type()
	if err := v.v.Struct(dto); err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range errs {
				switch e.Tag() {
				case requiredTag:
					return apierr.NewValidatorErr(errRequired, e.Field(), dtoName)
				case minTag:
					return apierr.NewValidatorErr(errMin, e.Field(), dtoName)
				case maxTag:
					return apierr.NewValidatorErr(errMax, e.Field(), dtoName)
				}
			}
		}
	}
	return nil
}
