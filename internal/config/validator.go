package config

import (
	"github.com/go-playground/validator/v10"
)

// NewValidator initializes playground validator
func NewValidator() *validator.Validate {
	return validator.New()
}
