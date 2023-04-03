package helper

func MonthLongToNumber(month string) int {
	var staticMonth []string
	staticMonth = append(staticMonth, "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec")
	var number int
	for i, value := range staticMonth {
		if value == month {
			number = i
		}
	}

	number += 1

	return number
}
