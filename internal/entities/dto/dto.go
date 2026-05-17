// Package dto describe a input/output DTO's for users
package dto

// DTO interface for all DTO's (used for DtoValidator for validate DTO's)
type DTO interface {
	Type() string
}
