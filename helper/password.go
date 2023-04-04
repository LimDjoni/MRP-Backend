package helper

import (
	"math/rand"
	"strings"
	"time"

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

func CreateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	length := 8

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"
	return str
}
