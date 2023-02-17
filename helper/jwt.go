package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey = []byte(GetEnvWithKey("JWT_SECRET_KEY"))

func GenerateToken(id uint, username string, email string) (string, error) {

	claim := jwt.MapClaims{}
	claim["id"] = id
	claim["username"] = username
	claim["email"] = email
	claim["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(GetEnvWithKey("JWT_SECRET_KEY")))

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func ValidateToken(token string) (*jwt.Token, error) {

	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("not authorize")
		}

		return jwtSecretKey, nil
	})

	if err != nil {
		return tokenParse, err
	}

	return tokenParse, nil
}
