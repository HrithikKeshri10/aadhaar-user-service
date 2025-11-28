package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Payload validates a struct and returns validation errors
func Payload(s interface{}) []ValidationError {
	var errors []ValidationError

	err := _validator.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", e.Field())
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", e.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
			case "len":
				msg = fmt.Sprintf("%s must be exactly %s characters", e.Field(), e.Param())
			case "oneof":
				msg = fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
			case "numeric":
				msg = fmt.Sprintf("%s must contain only numbers", e.Field())
			default:
				msg = fmt.Sprintf("%s is invalid", e.Field())
			}
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: msg,
			})
		}
	}

	return errors
}

// ValidatePagination validates and normalizes pagination parameters
func ValidatePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
