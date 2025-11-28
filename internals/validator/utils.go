package validator

import "github.com/go-playground/validator/v10"

var _validator *validator.Validate

func init() {
	_validator = validator.New()
}
