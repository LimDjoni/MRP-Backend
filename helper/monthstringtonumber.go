package helper

import "strconv"

func MonthStringToNumberString(month string) string {
	var staticMonth []string
	staticMonth = append(staticMonth, "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des")
	var number int
	for i, value := range staticMonth {
		if value == month {
			number = i
		}
	}

	number += 1

	numberString := strconv.Itoa(number)

	if len(numberString) == 1 {
		numberString = "0" + numberString
	}

	return numberString
}
