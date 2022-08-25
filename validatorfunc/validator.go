package validatorfunc

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func CheckDateString(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	layout := "2006-01-02"
	_, err := time.Parse(layout, value)

	if err != nil {
		return false
	}
	return true
}
