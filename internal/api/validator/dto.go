package validator

import (
	"errors"
	"todo/internal/api/apierr"
	"todo/internal/entities/dto"

	"github.com/go-playground/validator"
	"github.com/go-playground/validator/non-standard/validators"
)

var (
	requiredTag = "required"
	minTag      = "min"
	maxTag      = "max"
	notBlankTag = "notblank"

	errRequired = errors.New("required field is missing")
	errMin      = errors.New("field is too small")
	errMax      = errors.New("field is too large")
	errBlank    = errors.New("field is blank")
)

type DTO struct {
	v *validator.Validate
}

func NewDtoValidator() (*DTO, error) {
	v := validator.New()
	if err := v.RegisterValidation(notBlankTag, validators.NotBlank); err != nil {
		return nil, err
	}
	return &DTO{
		v: v,
	}, nil
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
				case notBlankTag:
					return apierr.NewValidatorErr(errBlank, e.Field(), dtoName)
				default:
					return apierr.NewValidatorErr(err, e.Field(), dtoName)
				}
			}
		}
	}
	return nil
}
