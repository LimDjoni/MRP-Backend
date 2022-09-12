package helper

import "strconv"

func MonthStringToNumberString(month string) string{
	var staticMonth []string
	staticMonth = append(staticMonth, "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec")
	var number int
	for i, value := range staticMonth {
		if value == month {
			number = i
		}
	}

	number += 1

	numberString  := strconv.Itoa(number)

	if len(numberString) == 1 {
		numberString = "0" +numberString
	}

	return  numberString
}
