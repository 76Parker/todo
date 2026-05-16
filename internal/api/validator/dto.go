package validator

import (
	"errors"
	"fmt"
	"strings"
	"todo/internal/entities/dto"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var (
	requiredTag = "required"
	minTag      = "min"
	maxTag      = "max"
	notBlankTag = "notblank"
)

type DTO struct {
	v *validator.Validate
}

func MustDtoValidator() *DTO {
	v := validator.New()
	if err := v.RegisterValidation(notBlankTag, validators.NotBlank); err != nil {
		panic(fmt.Sprintf("panic in validator.MustDtoValidator: %s", err.Error()))
	}
	return &DTO{
		v: v,
	}
}
func (v *DTO) Validate(dto dto.DTO) error {
	if err := v.v.Struct(dto); err != nil {
		if errs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range errs {
				field := strings.ToLower(e.Field())
				switch e.Tag() {
				case requiredTag:
					return fmt.Errorf("%s is required", field)
				case minTag:
					return fmt.Errorf("%s is too small", field)
				case maxTag:
					return fmt.Errorf("%s is too large", field)
				case notBlankTag:
					return fmt.Errorf("%s is blank", field)
				default:
					return fmt.Errorf("%s", field)
				}
			}
		}
	}
	return nil
}
