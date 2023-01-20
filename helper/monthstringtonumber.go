package helper

func MonthLongToNumber(month string) int {
	var staticMonth []string
	staticMonth = append(staticMonth, "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December")
	var number int
	for i, value := range staticMonth {
		if value == month {
			number = i
		}
	}

	number += 1

	return number
}
