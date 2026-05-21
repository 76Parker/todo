package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
)

const allowedPortsTag = "allowed_ports"

// AllowedPortsError custom error for guard privileged ports
type AllowedPortsError struct {
	Field string
	Value int
}

func (e AllowedPortsError) Error() string {
	return fmt.Sprintf("field %s: port %d is out of allowed range (1023-65535)", e.Field, e.Value)
}

// PlaygroundValidator it's a go-playground implementations of Validator
type PlaygroundValidator struct {
	v *validator.Validate
}

// NewValidator construct Validator
func NewValidator() (Validator, error) {
	v := &PlaygroundValidator{}
	v.v = validator.New()
	if err := v.addPortRangeValidation(); err != nil {
		return nil, err
	}
	return v, nil
}

// Validate validating Config
func (v *PlaygroundValidator) Validate(cfg *Config) error {
	var errs []error
	if err := v.v.Struct(cfg); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range ve {
				if e.Tag() == allowedPortsTag {
					f := e.Field()
					v, _ := strconv.Atoi(e.Value().(string))
					errs = append(errs, AllowedPortsError{f, v})
					continue
				}
				errs = append(errs, fmt.Errorf("config validation error: field: %s validation tag: %s",e.Field(), e.Tag()))
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// addPortRangeValidation валидирует поля с тегом allowed_ports на разрешенный диапазон портов
func (v *PlaygroundValidator) addPortRangeValidation() error {
	if err := v.v.RegisterValidation(allowedPortsTag, func(fl validator.FieldLevel) bool {
		port := fl.Field().Int()
		if port < 1023 || port > 65535 {
			return false
		}
		return true
	}); err != nil {
		return err
	}
	return nil
}
