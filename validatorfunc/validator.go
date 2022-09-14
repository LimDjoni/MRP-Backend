package validatorfunc

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func CheckDateString(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	isLayoutDateErr := false
	isLayoutTimeErr := false
	layout := "2006-01-02"
	layoutTime := "2006-01-02T00:00:00Z"
	_, err := time.Parse(layout, value)
	_, errTime := time.Parse(layoutTime, value)

	if err != nil {
		isLayoutDateErr = true
	}

	if errTime != nil {
		isLayoutTimeErr = true
	}

	if isLayoutDateErr == false || isLayoutTimeErr == false {
		return true
	}

	return false
}
