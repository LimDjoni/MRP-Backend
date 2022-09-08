package helper

import "strconv"

func MonthStringToNumberString(month string) string{
	var staticMonth []string
	staticMonth = append(staticMonth, "Januari","Februari","Maret","April","Mei","Juni","Juli","Agustus","September","Oktober","November","Desember")
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
