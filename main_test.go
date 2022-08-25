package main

import (
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	//"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func startSetup() (*gorm.DB, *validator.Validate) {
	dbUrlStg := "host=localhost user=postgres password=postgres dbname=deli_aje_development port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, _ := gorm.Open(postgres.Open(dbUrlStg), &gorm.Config{})

	var validate = validator.New()

	// Make Validation for Gender
	_ = validate.RegisterValidation("DateValidation", validatorfunc.CheckDateString)

	return db, validate
}

func TestRegisterUser(t *testing.T) {
	//Data for each test on the route
	tests := []struct {
		body          map[string]interface{}
		expectedError bool
		expectedCode  int
	}{
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"username": "Hendri",
				"email":    "hendini@gmail.com",
				"password": "hendini",
			},
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"username": "hendarin6",
				"email":    "hendarin6@gmail.com",
				"password": "hendarin6",
			},
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		req, _ := http.NewRequest(
			"POST",
			"/api/v1/user/register",
			payload,
		)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, err != nil, "register user")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "register user")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "register user")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			db.Unscoped().Where("email = ?", test.body["email"]).Delete(&user.User{})
		}
		//// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.body["username"], mapUnmarshal["username"], "register success user")
		assert.Equalf(t, test.body["email"], mapUnmarshal["email"], "register success user")

	}
}

func TestLoginUser(t *testing.T) {

}
