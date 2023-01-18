package validatorfunc

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func isStringInSlicesString(str string, slc []string) bool {
	for _, value := range slc {
		if value == str {
			return true
		}
	}
	return false
}

func ValidationPeriod(fl validator.FieldLevel) bool {

	var month []string
	month = append(month, "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec")
	value := fl.Field().String()

	if !strings.Contains(value, " ") {
		return false
	}

	period := strings.Split(value, " ")

	if len(period) > 2 || len(period) < 1 {
		return false
	}

	isPeriodMonthGood := isStringInSlicesString(period[0], month)

	var length = len([]rune(period[1]))
	_, err := strconv.Atoi(period[1])

	if !isPeriodMonthGood || length != 4 || err != nil {
		return false
	}

	return true
}
