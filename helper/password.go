package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func GeneratePasswordHash(password string) (string, error) {

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return string(newPasswordHash), err
	}

	return string(newPasswordHash), nil
}
