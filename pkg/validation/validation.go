package validation

import "github.com/go-playground/validator/v10"

var inputValidator *validator.Validate

func GetInputValidationInstance() *validator.Validate {
	if inputValidator == nil {
		inputValidator = validator.New()
	}
	return inputValidator
}
