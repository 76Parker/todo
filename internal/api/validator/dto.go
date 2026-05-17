// Package validator provide validate methods for user inputs
package validator

import (
	"errors"
	"fmt"
	"strconv"
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

// DTO it's a validator for dto.DTO
type DTO struct {
	v *validator.Validate
}

// MustDtoValidator it's a Must-constructor for DTO. MAY PANIC IF REGISTER VALIDATION HAVE ERROR
func MustDtoValidator() *DTO {
	v := validator.New()
	if err := v.RegisterValidation(notBlankTag, validators.NotBlank); err != nil {
		panic(fmt.Sprintf("panic in validator.MustDtoValidator: %s", err.Error()))
	}
	return &DTO{
		v: v,
	}
}

// Validate validating dto.DTO types
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

// ExtractAndValidateTaskID validate taskID path param and convert into int64
func ExtractAndValidateTaskID(taskID string) (int64, error) {
	if taskID == "" {
		return 0, fmt.Errorf("task_id is empty")
	}
	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid task_id")
	}
	if id <= 0 {
		return 0, fmt.Errorf("invalid task_id")
	}
	return id, nil
}
