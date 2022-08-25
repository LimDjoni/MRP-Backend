package helper

import "strconv"

func CreateIdNumber(number int) string {
	numberString := strconv.Itoa(number)

	lengthNumberNeededMinimum := 4

	nowLength := len(numberString)

	neededAddLength := lengthNumberNeededMinimum - nowLength

	for i := 0; i < neededAddLength; i++ {
		numberString = "0" + numberString
	}

	return numberString
}
