package apierr

import "fmt"

type DtoValidatorErr struct {
	DtoName string `json:"-"`
	Err     error  `json:"error"`
	Field   string `json:"field"`
}

func (v DtoValidatorErr) Error() string {
	return fmt.Sprintf("validation error: %s - field `%s`", v.Err.Error(), v.Field)
}

func NewValidatorErr(err error, field, dtoName string) error {
	return DtoValidatorErr{
		Err:     err,
		Field:   field,
		DtoName: dtoName,
	}
}
