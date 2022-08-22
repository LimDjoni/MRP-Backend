package validatorfunc

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       interface{}
}

func ValidateStruct(value interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	for _, fieldErr := range value.(validator.ValidationErrors) {
		var element ErrorResponse
		element.FailedField = fieldErr.StructNamespace()
		element.Tag = fieldErr.Tag()
		element.Value = fieldErr.Value()
		errors = append(errors, &element)
	}

	return errors
}
