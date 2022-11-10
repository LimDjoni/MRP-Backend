package main

import (
	"ajebackend/model/dmo"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"net/http"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dataId = 1
var idNumber = ""
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhlbmRhcmluNkBnbWFpbC5jb20iLCJpZCI6NTEsInVzZXJuYW1lIjoiaGVuZGFyaW42In0.77-X-McTZZUsf3yLKV9QNa0zziFBu922W020Xlz6MuU"

func startSetup() (*gorm.DB, *validator.Validate) {
	loadEnv()
	dbUrlStg := "host=localhost user=postgres password=postgres dbname=deli_aje_development port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, _ := gorm.Open(postgres.Open(dbUrlStg), &gorm.Config{})

	var validate = validator.New()

	// Make Validation for Date
	_ = validate.RegisterValidation("DateValidation", validatorfunc.CheckDateString)

	_ = validate.RegisterValidation("PeriodValidation", validatorfunc.ValidationPeriod)

	return db, validate
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

// User Handler Test

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
				"username": "hendarin7",
				"email":    "hendarin7@gmail.com",
				"password": "hendarin7",
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
		assert.Contains(t, mapUnmarshal, "username", "register user")
		assert.Contains(t, mapUnmarshal, "email", "register user")
		assert.Contains(t, mapUnmarshal, "ID", "register user")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "register user")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "register user")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "register user")
		assert.Contains(t, mapUnmarshal, "username", "register user")
		assert.Contains(t, mapUnmarshal, "password", "register user")
		assert.Contains(t, mapUnmarshal, "email", "register user")
		assert.Contains(t, mapUnmarshal, "is_active", "register user")
	}
}

func TestLoginUser(t *testing.T) {
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
				"data": "Hendri",
				"password": "hendini",
			},
		},
		{
			expectedError: false,
			expectedCode:  200,
			body: fiber.Map{
				"data": "hendarin6",
				"password": "hendarin6",
			},
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"data": "hendarin6",
				"password": "hendarin",
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
			"/api/v1/user/login",
			payload,
		)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, err != nil, "login user")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "login user")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "login user")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body

		if strings.Contains(test.body["data"].(string), "@") {
			assert.Equalf(t, test.body["email"], mapUnmarshal["data"], "login success user")
		} else {
			assert.Equalf(t, test.body["username"], mapUnmarshal["data"], "login success user")
		}
	}
}

// Transaction Handler Test

func TestListDataDN(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/transaction/list/dn",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "limit")
		assert.Contains(t, mapUnmarshal, "page")
		assert.Contains(t, mapUnmarshal, "total_rows")
		assert.Contains(t, mapUnmarshal, "total_pages")
		assert.Contains(t, mapUnmarshal, "data")
	}
}

func TestDetailTransactionDN(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: 38,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/transaction/detail/dn/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list detail dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list detail dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list detail dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "list detail dn")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "list detail dn")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "list detail dn")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_id", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo", "list detail dn")
		assert.Contains(t, mapUnmarshal, "id_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "transaction_type", "list detail dn")
		assert.Contains(t, mapUnmarshal, "shipping_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quantity", "list detail dn")
		assert.Contains(t, mapUnmarshal, "tugboat_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "barge_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "vessel_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "seller", "list detail dn")
		assert.Contains(t, mapUnmarshal, "customer_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "loading_port_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "loading_port_location", "list detail dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_location", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_destination_port", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skb_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skb_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skab_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skab_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "royalty_rate", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_currency", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_ntpn", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_billing_code", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_total", "list detail dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_currency", "list detail dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_ntpn", "list detail dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_billing_code", "list detail dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_total", "list detail dn")
		assert.Contains(t, mapUnmarshal, "lhv_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "lhv_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "surveyor_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "cow_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "cow_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "coa_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "coa_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_tm_ar", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_im_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_ar", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_vm_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_fc_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_ar", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_ar", "list detail dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_adb", "list detail dn")
		assert.Contains(t, mapUnmarshal, "barging_distance", "list detail dn")
		assert.Contains(t, mapUnmarshal, "sales_system", "list detail dn")
		assert.Contains(t, mapUnmarshal, "invoice_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "invoice_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_unit", "list detail dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_total", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_reconciliation_letter", "list detail dn")
		assert.Contains(t, mapUnmarshal, "contract_date", "list detail dn")
		assert.Contains(t, mapUnmarshal, "contract_number", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_buyer_name", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_industry_type", "list detail dn")
		assert.Contains(t, mapUnmarshal, "dmo_category", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skb_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "skab_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "bl_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "royalti_provision_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "royalti_final_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "cow_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "coa_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "invoice_and_contract_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "lhv_document_link", "list detail dn")
		assert.Contains(t, mapUnmarshal, "is_not_claim", "list detail dn")
	}
}

func TestCreateTransactionDN(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		body map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body: fiber.Map{},
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"shipping_date": "2022-02-01",
				"coa_date": "2022-02-01",
				"quantity": 1023.122,
				"tugboat_name": "AJE",
				"dp_royalty_ntpn": "A123SSSS",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"coa_date": "2022-22-01",
				"quantity": 1023.122,
				"tugboat_name": "AJE",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"coa_date": "2022-22-01",
				"quantity": 1023.122,
				"tugboat_name": "AJE",
				"dp_royalty_ntpn": "A123SSSS",
				"dp_royalty_billing_code": "A123SS",
				"payment_dp_royalty_ntpn": "A123SS",
				"payment_dp_royalty_billing_code": "A123SS",
			},
			token: token,
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
			"/api/v1/transaction/create/dn",
			payload,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create data dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create data dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create data dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			dataId = int(mapUnmarshal["ID"].(float64))
			idNumber = mapUnmarshal["id_number"].(string)
		}
		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "create data dn")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create data dn")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create data dn")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_id", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo", "create data dn")
		assert.Contains(t, mapUnmarshal, "id_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "transaction_type", "create data dn")
		assert.Contains(t, mapUnmarshal, "shipping_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "quantity", "create data dn")
		assert.Contains(t, mapUnmarshal, "tugboat_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "barge_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "vessel_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "seller", "create data dn")
		assert.Contains(t, mapUnmarshal, "customer_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_location", "create data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_location", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_destination_port", "create data dn")
		assert.Contains(t, mapUnmarshal, "skb_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "skb_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "skab_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "skab_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "royalty_rate", "create data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_currency", "create data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_ntpn", "create data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_billing_code", "create data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_total", "create data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_currency", "create data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_ntpn", "create data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_billing_code", "create data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_total", "create data dn")
		assert.Contains(t, mapUnmarshal, "lhv_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "lhv_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "surveyor_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "cow_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "cow_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "coa_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "coa_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_tm_ar", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_im_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_ar", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_vm_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_fc_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_ar", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_ar", "create data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_adb", "create data dn")
		assert.Contains(t, mapUnmarshal, "barging_distance", "create data dn")
		assert.Contains(t, mapUnmarshal, "sales_system", "create data dn")
		assert.Contains(t, mapUnmarshal, "invoice_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "invoice_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_unit", "create data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_total", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_reconciliation_letter", "create data dn")
		assert.Contains(t, mapUnmarshal, "contract_date", "create data dn")
		assert.Contains(t, mapUnmarshal, "contract_number", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_buyer_name", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_industry_type", "create data dn")
		assert.Contains(t, mapUnmarshal, "dmo_category", "create data dn")
		assert.Contains(t, mapUnmarshal, "skb_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "skab_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "bl_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "royalti_provision_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "royalti_final_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "cow_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "coa_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "invoice_and_contract_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "lhv_document_link", "create data dn")
		assert.Contains(t, mapUnmarshal, "is_not_claim", "create data dn")
		assert.Contains(t, mapUnmarshal, "is_migration", "create data dn")
	}
}

func TestUpdateTransactionDN(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		body map[string]interface{}
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body: fiber.Map{},
			id: 49,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			body: fiber.Map{
				"ID": 49,
				"CreatedAt": "2022-08-24T11:59:40.094282+07:00",
				"UpdatedAt": "2022-08-24T11:59:40.094282+07:00",
				"DeletedAt": nil,
				"dmo_id": nil,
				"dmo": nil,
				"id_number": "DN-2022-8-0035",
				"transaction_type": "DN",
				"shipping_date": "2022-01-01",
				"quantity": 1023.122,
				"tugboat_name": "AJE",
				"barge_name": "SHIPIPP",
				"vessel_name": "",
				"seller": "",
				"customer_name": "",
				"loading_port_name": "",
				"loading_port_location": "",
				"unloading_port_name": "",
				"unloading_port_location": "",
				"dmo_destination_port": "",
				"skb_date": nil,
				"skb_number": "",
				"skab_date": nil,
				"skab_number": "",
				"bill_of_lading_date": nil,
				"bill_of_lading_number": "",
				"royalty_rate": 0,
				"dp_royalty_currency": "IDR",
				"dp_royalty_date": nil,
				"dp_royalty_ntpn": nil,
				"dp_royalty_billing_code": nil,
				"dp_royalty_total": 0,
				"payment_dp_royalty_currency": "IDR",
				"payment_dp_royalty_date": nil,
				"payment_dp_royalty_ntpn": nil,
				"payment_dp_royalty_billing_code": nil,
				"payment_dp_royalty_total": 0,
				"lhv_date": nil,
				"lhv_number": "",
				"surveyor_name": "",
				"cow_date": nil,
				"cow_number": "",
				"coa_date": "2022-02-01",
				"coa_number": "",
				"quality_tm_ar": 0,
				"quality_im_adb": 0,
				"quality_ash_ar": 0,
				"quality_ash_adb": 0,
				"quality_vm_adb": 0,
				"quality_fc_adb": 0,
				"quality_ts_ar": 0,
				"quality_ts_adb": 0,
				"quality_calories_ar": 0,
				"quality_calories_adb": 0,
				"barging_distance": 0,
				"sales_system": "",
				"invoice_date": nil,
				"invoice_number": "",
				"invoice_price_unit": 0,
				"invoice_price_total": 0,
				"dmo_reconciliation_letter": "",
				"contract_date": nil,
				"contract_number": "",
				"dmo_buyer_name": "",
				"dmo_industry_type": "",
				"skb_document_link": "",
				"skab_document_link": "",
				"bl_document_link": "",
				"royalti_provision_document_link": "",
				"royalti_final_document_link": "",
				"cow_document_link": "",
				"coa_document_link": "",
				"invoice_and_contract_document_link": "",
				"lhv_document_link": "",
				"is_not_claim": false,
			},
			id: dataId,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"ID": 49,
				"CreatedAt": "2022-08-24T11:59:40.094282+07:00",
				"UpdatedAt": "2022-08-24T11:59:40.094282+07:00",
				"DeletedAt": nil,
				"dmo_id": nil,
				"dmo": nil,
				"id_number": "DN-2022-8-0035",
				"transaction_type": "DN",
				"shipping_date": nil,
				"quantity": 1023.122,
				"tugboat_name": "AJE",
				"barge_name": "SHIPIPP",
				"vessel_name": "",
				"seller": "",
				"customer_name": "",
				"loading_port_name": "",
				"loading_port_location": "",
				"unloading_port_name": "",
				"unloading_port_location": "",
				"dmo_destination_port": "",
				"skb_date": "2022",
				"skb_number": "",
				"skab_date": nil,
				"skab_number": "",
				"bill_of_lading_date": nil,
				"bill_of_lading_number": "",
				"royalty_rate": 0,
				"dp_royalty_currency": "IDR",
				"dp_royalty_date": nil,
				"dp_royalty_ntpn": "",
				"dp_royalty_billing_code": "",
				"dp_royalty_total": 0,
				"payment_dp_royalty_currency": "IDR",
				"payment_dp_royalty_date": nil,
				"payment_dp_royalty_ntpn": "",
				"payment_dp_royalty_billing_code": "",
				"payment_dp_royalty_total": 0,
				"lhv_date": nil,
				"lhv_number": "",
				"surveyor_name": "",
				"cow_date": nil,
				"cow_number": "",
				"coa_date": "2022-02-01",
				"coa_number": "",
				"quality_tm_ar": 0,
				"quality_im_adb": 0,
				"quality_ash_ar": 0,
				"quality_ash_adb": 0,
				"quality_vm_adb": 0,
				"quality_fc_adb": 0,
				"quality_ts_ar": 0,
				"quality_ts_adb": 0,
				"quality_calories_ar": 0,
				"quality_calories_adb": 0,
				"barging_distance": 0,
				"sales_system": "",
				"invoice_date": nil,
				"invoice_number": "",
				"invoice_price_unit": 0,
				"invoice_price_total": 0,
				"dmo_reconciliation_letter": "",
				"contract_date": nil,
				"contract_number": "",
				"dmo_buyer_name": "",
				"dmo_industry_type": "",
				"skb_document_link": "",
				"skab_document_link": "",
				"bl_document_link": "",
				"royalti_provision_document_link": "",
				"royalti_final_document_link": "",
				"cow_document_link": "",
				"coa_document_link": "",
				"invoice_and_contract_document_link": "",
				"lhv_document_link": "",
				"is_not_claim": false,
			},
			id: 49,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			body: fiber.Map{
				"ID": 49,
				"CreatedAt": "2022-08-24T11:59:40.094282+07:00",
				"UpdatedAt": "2022-08-24T11:59:40.094282+07:00",
				"DeletedAt": nil,
				"dmo_id": nil,
				"dmo": nil,
				"id_number": "DN-2022-8-0035",
				"transaction_type": "DN",
				"shipping_date": "2022-01-01",
				"quantity": 1023.122,
				"tugboat_name": "AJE",
				"barge_name": "SHIPIPP",
				"vessel_name": "",
				"seller": "",
				"customer_name": "",
				"loading_port_name": "",
				"loading_port_location": "",
				"unloading_port_name": "",
				"unloading_port_location": "",
				"dmo_destination_port": "",
				"skb_date": nil,
				"skb_number": "",
				"skab_date": nil,
				"skab_number": "",
				"bill_of_lading_date": nil,
				"bill_of_lading_number": "",
				"royalty_rate": 0,
				"dp_royalty_currency": "IDR",
				"dp_royalty_date": nil,
				"dp_royalty_ntpn": nil,
				"dp_royalty_billing_code": nil,
				"dp_royalty_total": 0,
				"payment_dp_royalty_currency": "IDR",
				"payment_dp_royalty_date": nil,
				"payment_dp_royalty_ntpn": nil,
				"payment_dp_royalty_billing_code": nil,
				"payment_dp_royalty_total": 0,
				"lhv_date": nil,
				"lhv_number": "",
				"surveyor_name": "",
				"cow_date": nil,
				"cow_number": "",
				"coa_date": "2022-02-01",
				"coa_number": "",
				"quality_tm_ar": 0,
				"quality_im_adb": 0,
				"quality_ash_ar": 0,
				"quality_ash_adb": 0,
				"quality_vm_adb": 0,
				"quality_fc_adb": 0,
				"quality_ts_ar": 0,
				"quality_ts_adb": 0,
				"quality_calories_ar": 0,
				"quality_calories_adb": 0,
				"barging_distance": 0,
				"sales_system": "",
				"invoice_date": nil,
				"invoice_number": "",
				"invoice_price_unit": 0,
				"invoice_price_total": 0,
				"dmo_reconciliation_letter": "",
				"contract_date": nil,
				"contract_number": "",
				"dmo_buyer_name": "",
				"dmo_industry_type": "",
				"skb_document_link": "",
				"skab_document_link": "",
				"bl_document_link": "",
				"royalti_provision_document_link": "",
				"royalti_final_document_link": "",
				"cow_document_link": "",
				"coa_document_link": "",
				"invoice_and_contract_document_link": "",
				"lhv_document_link": "",
				"is_not_claim": false,
			},
			id: 904,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/transaction/update/dn/%v", test.id)
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		req, _ := http.NewRequest(
			"PUT",
			url,
			payload,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "update data dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update data dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update data dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update data dn")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update data dn")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update data dn")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_id", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo", "update data dn")
		assert.Contains(t, mapUnmarshal, "id_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "transaction_type", "update data dn")
		assert.Contains(t, mapUnmarshal, "shipping_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "quantity", "update data dn")
		assert.Contains(t, mapUnmarshal, "tugboat_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "barge_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "vessel_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "seller", "update data dn")
		assert.Contains(t, mapUnmarshal, "customer_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_location", "update data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_location", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_destination_port", "update data dn")
		assert.Contains(t, mapUnmarshal, "skb_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "skb_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "skab_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "skab_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "royalty_rate", "update data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_currency", "update data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_ntpn", "update data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_billing_code", "update data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_total", "update data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_currency", "update data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_ntpn", "update data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_billing_code", "update data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_total", "update data dn")
		assert.Contains(t, mapUnmarshal, "lhv_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "lhv_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "surveyor_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "cow_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "cow_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "coa_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "coa_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_tm_ar", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_im_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_ar", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_vm_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_fc_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_ar", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_ar", "update data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_adb", "update data dn")
		assert.Contains(t, mapUnmarshal, "barging_distance", "update data dn")
		assert.Contains(t, mapUnmarshal, "sales_system", "update data dn")
		assert.Contains(t, mapUnmarshal, "invoice_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "invoice_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_unit", "update data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_total", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_reconciliation_letter", "update data dn")
		assert.Contains(t, mapUnmarshal, "contract_date", "update data dn")
		assert.Contains(t, mapUnmarshal, "contract_number", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_buyer_name", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_industry_type", "update data dn")
		assert.Contains(t, mapUnmarshal, "dmo_category", "update data dn")
		assert.Contains(t, mapUnmarshal, "skb_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "skab_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "bl_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "royalti_provision_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "royalti_final_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "cow_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "coa_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "invoice_and_contract_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "lhv_document_link", "update data dn")
		assert.Contains(t, mapUnmarshal, "is_not_claim", "update data dn")
		assert.Contains(t, mapUnmarshal, "is_migration", "update data dn")
	}
}

func TestUpdateDocumentTransactionDN(t *testing.T) {
	openDocumentPdf , errOpenDocumentPdf := os.Open("upload_test/output.pdf")
	openDocumentPng , errOpenDocumentPng := os.Open("upload_test/output.png")

	assert.Nilf(t, errOpenDocumentPdf, "update document data dn")
	assert.Nilf(t, errOpenDocumentPng, "update document data dn")

	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		file string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 49,
			token: "asdawfaeac",
			file: "upload_test/output.pdf",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: dataId,
			token: token,
			file: "upload_test/output.pdf",
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: 49,
			token: token,
			file: "upload_test/output.png",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 904,
			token: token,
			file: "upload_test/output.pdf",
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	defer openDocumentPdf.Close()
	defer openDocumentPng.Close()

	for _, test := range tests {
		bodyRequest := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyRequest)

		fw, err := writer.CreateFormFile("document", test.file)
		if err != nil {
			assert.Nilf(t, err, "update document data dn")
		}
		file, err := os.Open("upload_test/output.pdf")
		if err != nil {
			assert.Nilf(t, err, "update document data dn")
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			assert.Nilf(t, err, "update document data dn")
		}
		writer.Close()

		urlApi := fmt.Sprintf("/api/v1/transaction/update/document/dn/%v/lhv", test.id)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			bytes.NewReader(bodyRequest.Bytes()),
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update document data dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update document data dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update document data dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update document data dn")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update document data dn")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update document data dn")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_id", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo", "update document data dn")
		assert.Contains(t, mapUnmarshal, "id_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "transaction_type", "update document data dn")
		assert.Contains(t, mapUnmarshal, "shipping_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quantity", "update document data dn")
		assert.Contains(t, mapUnmarshal, "tugboat_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "barge_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "vessel_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "seller", "update document data dn")
		assert.Contains(t, mapUnmarshal, "customer_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "loading_port_location", "update document data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "unloading_port_location", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_destination_port", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skb_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skb_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skab_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skab_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "bill_of_lading_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "royalty_rate", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_currency", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_ntpn", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_billing_code", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dp_royalty_total", "update document data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_currency", "update document data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_ntpn", "update document data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_billing_code", "update document data dn")
		assert.Contains(t, mapUnmarshal, "payment_dp_royalty_total", "update document data dn")
		assert.Contains(t, mapUnmarshal, "lhv_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "lhv_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "surveyor_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "cow_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "cow_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "coa_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "coa_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_tm_ar", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_im_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_ar", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_ash_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_vm_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_fc_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_ar", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_ts_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_ar", "update document data dn")
		assert.Contains(t, mapUnmarshal, "quality_calories_adb", "update document data dn")
		assert.Contains(t, mapUnmarshal, "barging_distance", "update document data dn")
		assert.Contains(t, mapUnmarshal, "sales_system", "update document data dn")
		assert.Contains(t, mapUnmarshal, "invoice_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "invoice_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_unit", "update document data dn")
		assert.Contains(t, mapUnmarshal, "invoice_price_total", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_reconciliation_letter", "update document data dn")
		assert.Contains(t, mapUnmarshal, "contract_date", "update document data dn")
		assert.Contains(t, mapUnmarshal, "contract_number", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_buyer_name", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_industry_type", "update document data dn")
		assert.Contains(t, mapUnmarshal, "dmo_category", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skb_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "skab_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "bl_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "royalti_provision_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "royalti_final_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "cow_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "coa_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "invoice_and_contract_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "lhv_document_link", "update document data dn")
		assert.Contains(t, mapUnmarshal, "is_not_claim", "update document data dn")
		assert.Contains(t, mapUnmarshal, "is_migration", "update document data dn")
	}
}

func TestDeleteTransactionDN(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: dataId,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/transaction/delete/dn/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete data dn")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete data dn")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete data dn")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete data dn")
	}
}

// Minerba Handler Test
var idMinerba = 0
var periodMinerba = ""
func TestListMinerba(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/minerba/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "limit")
		assert.Contains(t, mapUnmarshal, "page")
		assert.Contains(t, mapUnmarshal, "total_rows")
		assert.Contains(t, mapUnmarshal, "total_pages")
		assert.Contains(t, mapUnmarshal, "data")
	}
}

func TestListDataDNWithoutMinerba(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token:         "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token:         "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token:         token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/minerba/list/transaction",
			nil,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data dn without minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data dn without minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data dn without minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "list")
	}
}

func TestCreateMinerba(t *testing.T) {
	var listDn []int
	listDn = append(listDn, 7, 8, 9)

	var errorListDn []int
	errorListDn = append(errorListDn, 1, 2, 3)

	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"period":       "Jun 2022",
				"list_data_dn": listDn,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period":       "Jun 2022",
				"list_data_dn": listDn,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period":       "Dec 2022",
				"list_data_dn": listDn,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period":       "Mei 2022",
				"list_data_dn": errorListDn,
			},
			token: token,
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
			"/api/v1/minerba/create",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create data minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create data minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create data minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idMinerba = int(mapUnmarshal["ID"].(float64))
			periodMinerba = mapUnmarshal["id_number"].(string)
		}

		assert.Contains(t, mapUnmarshal, "ID", "create data minerba")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create data minerba")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create data minerba")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create data minerba")
		assert.Contains(t, mapUnmarshal, "period", "create data minerba")
		assert.Contains(t, mapUnmarshal, "id_number", "create data minerba")
		assert.Contains(t, mapUnmarshal, "quantity", "create data minerba")
		assert.Contains(t, mapUnmarshal, "sp3medn_document_link", "create data minerba")
		assert.Contains(t, mapUnmarshal, "recap_dmo_document_link", "create data minerba")
		assert.Contains(t, mapUnmarshal, "detail_dmo_document_link", "create data minerba")
		assert.Contains(t, mapUnmarshal, "sp3meln_document_link", "create data minerba")
		assert.Contains(t, mapUnmarshal, "insw_export_document_link", "create data minerba")
	}
}

func TestUpdateMinerba(t *testing.T) {
	var listDn []int
	listDn = append(listDn, 71, 72, 73, 74)

	var errorListDn []int
	errorListDn = append(errorListDn, 1, 2, 3)

	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			body: fiber.Map{
				"list_data_dn": listDn,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			body: fiber.Map{
				"list_data_dn": errorListDn,
			},
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		urlString := fmt.Sprintf("/api/v1/minerba/update/%v", idMinerba)
		req, _ := http.NewRequest(
			"PUT",
			urlString,
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "update data minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update data minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update data minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idMinerba = int(mapUnmarshal["ID"].(float64))
		}

		assert.Contains(t, mapUnmarshal, "ID", "update data minerba")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update data minerba")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update data minerba")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update data minerba")
		assert.Contains(t, mapUnmarshal, "period", "update data minerba")
		assert.Contains(t, mapUnmarshal, "id_number", "update data minerba")
		assert.Contains(t, mapUnmarshal, "quantity", "update data minerba")
		assert.Contains(t, mapUnmarshal, "sp3medn_document_link", "update data minerba")
		assert.Contains(t, mapUnmarshal, "recap_dmo_document_link", "update data minerba")
		assert.Contains(t, mapUnmarshal, "detail_dmo_document_link", "update data minerba")
		assert.Contains(t, mapUnmarshal, "sp3meln_document_link", "update data minerba")
		assert.Contains(t, mapUnmarshal, "insw_export_document_link", "update data minerba")
	}
}

func TestDetailMinerba(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idMinerba,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/minerba/detail/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list detail minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list detail minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list detail minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "detail", "list detail minerba")
		assert.Contains(t, mapUnmarshal, "list", "list detail minerba")
	}
}

func TestCheckValidPeriodMinerba(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": periodMinerba,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Mar 20",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  200,
			body: fiber.Map{
				"period": "Dec 2050",
			},
			token: token,
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
			"/api/v1/minerba/check",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "check data minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "check data minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "check data minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		assert.Contains(t, mapUnmarshal, "message", "check data minerba")
	}
}

func TestUpdateDocumentMinerba(t *testing.T) {

	bodyString := make(map[string][]map[string]interface{})

	bodyString["data"] = []map[string]interface{}{
		{
			"Key": "LM-2022-03-0008/sp3medn.xlsx",
			"key": "LM-2022-03-0008/sp3medn.xlsx",
			"ETag": "23acbb2e206c924f5b438162ffa0b425",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/LM-2022-03-0008/sp3medn.xlsx",
		},
		{
			"Key": "LM-2022-03-0008/recapdmo.xlsx",
			"key": "LM-2022-03-0008/recapdmo.xlsx",
			"ETag": "610288c6357a94aba8d5c4e04ee588e3",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/LM-2022-03-0008/recapdmo.xlsx",
		},
		{
			"Key": "LM-2022-03-0008/detaildmo.xlsx",
			"key": "LM-2022-03-0008/detaildmo.xlsx",
			"ETag": "fbb861e4a88950eb7d38845219f99f43",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/LM-2022-03-0008/detaildmo.xlsx",
		},
	}

	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		body map[string][]map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 49,
			token: "asdawfaeac",
			body: bodyString,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idMinerba,
			token: token,
			body: bodyString,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 904,
			token: token,
			body: bodyString,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		urlApi := fmt.Sprintf("/api/v1/minerba/update/document/%v", test.id)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			payload,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update data document minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update data document minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update data document minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "period", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "id_number", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "quantity", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "sp3medn_document_link", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "recap_dmo_document_link", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "detail_dmo_document_link", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "sp3meln_document_link", "update data document minerba")
		assert.Contains(t, mapUnmarshal, "insw_export_document_link", "update data document minerba")
	}
}

func TestDeleteMinerba(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idMinerba,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/minerba/delete/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete data minerba")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete data minerba")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete data minerba")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete data minerba")
	}
}

// Dmo Handler Test

var idDmo = 0
func TestListDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/dmo/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "limit")
		assert.Contains(t, mapUnmarshal, "page")
		assert.Contains(t, mapUnmarshal, "total_rows")
		assert.Contains(t, mapUnmarshal, "total_pages")
		assert.Contains(t, mapUnmarshal, "data")
	}
}

func TestListDataDNWithoutDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token:         "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token:         "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token:         token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/dmo/list/transaction",
			nil,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data dn without dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data dn without dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data dn without dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "barge_transaction", "list data dn without dmo")
		assert.Contains(t, mapUnmarshal, "vessel_transaction", "list data dn without dmo")
	}
}

func TestCreateDmo(t *testing.T) {
	var traderList []int

	var endUser int

	endUser = 77

	traderList = append(traderList, 78, 79)

	var vesselAdjustment []dmo.VesselAdjustmentInput
	
	vesselAdjustment = append(vesselAdjustment, dmo.VesselAdjustmentInput{
		VesselName: "MV. PACIFIC BULK",
		Quantity:   7504.086,
		Adjustment: -10,
	})

	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"period": "Dec 2022",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": vesselAdjustment,
				"transaction_barge": []int{150,151},
				"transaction_vessel": []int{152},
				"is_document_custom": false,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Dec 2022",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": vesselAdjustment,
				"transaction_barge": []int{150,151},
				"transaction_vessel": []int{152},
				"is_document_custom": false,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Dec 2023",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": vesselAdjustment,
				"transaction_barge": []int{150,151},
				"transaction_vessel": []int{152},
				"is_document_custom": false,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Dec 2023",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": vesselAdjustment,
				"transaction_barge": []int{152},
				"transaction_vessel": []int{152},
				"is_document_custom": false,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Jun 2022",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": vesselAdjustment,
				"transaction_barge": []int{},
				"transaction_vessel": []int{},
				"is_document_custom": false,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"period": "Dec 2022",
				"trader": traderList,
				"end_user":  endUser,
				"vessel_adjustment": []dmo.VesselAdjustmentInput{},
				"transaction_barge": []int{},
				"transaction_vessel": []int{153},
				"is_document_custom": false,
			},
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyRequest := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyRequest)

		period, err := writer.CreateFormField("period")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		periodMarshal, _ := json.Marshal(test.body["period"])

		period.Write(periodMarshal)

		trader, err := writer.CreateFormField("trader")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		traderMarshal, _ := json.Marshal(test.body["trader"])

		trader.Write(traderMarshal)

		endUserForm, err := writer.CreateFormField("end_user")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		endUserMarshal, _ := json.Marshal(test.body["end_user"])

		endUserForm.Write(endUserMarshal)

		vesselAdjustmentForm, err := writer.CreateFormField("vessel_adjustment")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		vesselAdjustmentFormMarshal, _ := json.Marshal(test.body["vessel_adjustment"])

		vesselAdjustmentForm.Write(vesselAdjustmentFormMarshal)

		transactionBargeForm, err := writer.CreateFormField("transaction_barge")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		transactionBargeFormMarshal, _ := json.Marshal(test.body["transaction_barge"])

		transactionBargeForm.Write(transactionBargeFormMarshal)

		transactionVesselForm, err := writer.CreateFormField("transaction_vessel")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		transactionVesselFormMarshal, _ := json.Marshal(test.body["transaction_vessel"])

		transactionVesselForm.Write(transactionVesselFormMarshal)


		isDocumentCustom, err := writer.CreateFormField("is_document_custom")
		if err != nil {
			assert.Nilf(t, err, "create data dmo")
		}
		isDocumentCustomMarshal, _ := json.Marshal(test.body["is_document_custom"])

		isDocumentCustom.Write(isDocumentCustomMarshal)

		writer.Close()

		req, _ := http.NewRequest(
			"POST",
			"/api/v1/dmo/create",
			bytes.NewReader(bodyRequest.Bytes()),
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create data dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create data dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create data dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idDmo = int(mapUnmarshal["ID"].(float64))
		}

		assert.Contains(t, mapUnmarshal, "ID", "create data dmo")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create data dmo")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create data dmo")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create data dmo")
		assert.Contains(t, mapUnmarshal, "period", "create data dmo")
		assert.Contains(t, mapUnmarshal, "id_number", "create data dmo")
		assert.Contains(t, mapUnmarshal, "type", "create data dmo")
		assert.Contains(t, mapUnmarshal, "barge_total_quantity", "create data dmo")
		assert.Contains(t, mapUnmarshal, "barge_adjustment", "create data dmo")
		assert.Contains(t, mapUnmarshal, "barge_grand_total_quantity", "create data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_total_quantity", "create data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_adjustment", "create data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_grand_total_quantity", "create data dmo")
		assert.Contains(t, mapUnmarshal, "reconciliation_letter_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "signed_reconciliation_letter_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_downloaded", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_signed", "create data dmo")
		assert.Contains(t, mapUnmarshal, "bast_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "signed_bast_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_downloaded", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_signed", "create data dmo")
		assert.Contains(t, mapUnmarshal, "statement_letter_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "signed_statement_letter_document_link", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_downloaded", "create data dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_signed", "create data dmo")
	}
}

func TestDetailDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/dmo/detail/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list detail dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list detail dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list detail dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "detail", "list detail dmo")
		assert.Contains(t, mapUnmarshal, "list", "list detail dmo")
	}
}

func TestUpdateDocumentDmo(t *testing.T) {

	bodyString := make(map[string][]map[string]interface{})

	bodyString["data"] = []map[string]interface{}{
		{
			"Key": "DD-2022-03-0008/bast.pdf",
			"key": "DD-2022-03-0008/bast.pdf",
			"ETag": "23acbb2e206c924f5b438162ffa0b425",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DD-2022-03-0008/bast.pdf",
		},
		{
			"Key": "DD-2022-03-0008/berita_acara.pdf",
			"key": "DD-2022-03-0008/berita_acara.pdf",
			"ETag": "610288c6357a94aba8d5c4e04ee588e3",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DD-2022-03-0008/berita_acara.pdf",
		},
		{
			"Key": "DD-2022-03-0008/surat_pernyataan.pdf",
			"key": "DD-2022-03-0008/surat_pernyataan.pdf",
			"ETag": "fbb861e4a88950eb7d38845219f99f43",
			"Bucket": "deli-aje",
			"Location": "https://deli-aje.s3.ap-southeast-1.amazonaws.com/DD-2022-03-0008/surat_pernyataan.pdf",
		},
	}

	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		body map[string][]map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 49,
			token: "asdawfaeac",
			body: bodyString,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			body: bodyString,
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idDmo,
			token: token,
			body: bodyString,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 904,
			token: token,
			body: bodyString,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		urlApi := fmt.Sprintf("/api/v1/dmo/update/document/%v", test.id)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			payload,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update document data dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update document data dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update document data dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "period", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "id_number", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "type", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "barge_total_quantity", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "barge_adjustment", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "barge_grand_total_quantity", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_total_quantity", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_adjustment", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "vessel_grand_total_quantity", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "reconciliation_letter_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "signed_reconciliation_letter_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_downloaded", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_signed", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "bast_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "signed_bast_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_downloaded", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_signed", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "statement_letter_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "signed_statement_letter_document_link", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_downloaded", "update document data dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_signed", "update document data dmo")
	}
}

func TestUpdateIsDownloadedDocumentDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		typeDocument string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: idDmo,
			token: "asdawfaeac",
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "statement_letter",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "reconciliation_letter",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idDmo,
			token: token,
			typeDocument: "link_onl",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1000,
			token: token,
			typeDocument: "bast",
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		urlApi := fmt.Sprintf("/api/v1/dmo/update/document/downloaded/%v/%v", test.id, test.typeDocument)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update document data downloaded dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update document data downloaded dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update document data downloaded dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "period", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "id_number", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "type", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "barge_total_quantity", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "barge_adjustment", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "barge_grand_total_quantity", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "vessel_total_quantity", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "vessel_adjustment", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "vessel_grand_total_quantity", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "reconciliation_letter_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "signed_reconciliation_letter_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_downloaded", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_signed", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "bast_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "signed_bast_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_downloaded", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_signed", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "statement_letter_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "signed_statement_letter_document_link", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_downloaded", "update document data downloaded dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_signed", "update document data downloaded dmo")
	}
}

func TestUpdateTrueIsSignedDocumentDmo(t *testing.T) {
	openDocumentPdf , errOpenDocumentPdf := os.Open("upload_test/output.pdf")
	defer openDocumentPdf.Close()
	assert.Nilf(t, errOpenDocumentPdf, "update document data signed dmo")

	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		typeDocument string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: idDmo,
			token: "asdawfaeac",
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "statement_letter",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "reconciliation_letter",
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idDmo,
			token: token,
			typeDocument: "link_onl",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1000,
			token: token,
			typeDocument: "bast",
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyRequest := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyRequest)

		fw, err := writer.CreateFormFile("document", "upload_test/output.pdf")
		if err != nil {
			assert.Nilf(t, err, "update document data signed dmo")
		}
		file, err := os.Open("upload_test/output.pdf")
		if err != nil {
			assert.Nilf(t, err, "update document data signed dmo")
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			assert.Nilf(t, err, "update document data signed dmo")
		}
		writer.Close()

		urlApi := fmt.Sprintf("/api/v1/dmo/update/document/signed/%v/%v", test.id, test.typeDocument)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			bytes.NewReader(bodyRequest.Bytes()),
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update document data signed dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update document data signed dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update document data signed dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "period", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "id_number", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "type", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_adjustment", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_grand_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_adjustment", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_grand_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "reconciliation_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_reconciliation_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_signed", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "bast_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_bast_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_signed", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "statement_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_statement_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_signed", "update document data signed dmo")
	}
}

func TestUpdateFalseIsSignedDocumentDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
		typeDocument string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: idDmo,
			token: "asdawfaeac",
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "bast",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "statement_letter",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
			typeDocument: "reconciliation_letter",
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idDmo,
			token: token,
			typeDocument: "link_onl",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1000,
			token: token,
			typeDocument: "bast",
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {

		urlApi := fmt.Sprintf("/api/v1/dmo/update/document/not_signed/%v/%v", test.id, test.typeDocument)

		req, _ := http.NewRequest(
			"PUT",
			urlApi,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, errTest := app.Test(req, -1)


		assert.Equalf(t, test.expectedError, errTest != nil, "update document data signed dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update document data signed dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update document data signed dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "period", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "id_number", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "type", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_adjustment", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "barge_grand_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_adjustment", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "vessel_grand_total_quantity", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "reconciliation_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_reconciliation_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_reconciliation_letter_signed", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "bast_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_bast_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_bast_document_signed", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "statement_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "signed_statement_letter_document_link", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_downloaded", "update document data signed dmo")
		assert.Contains(t, mapUnmarshal, "is_statement_letter_signed", "update document data signed dmo")
	}
}

func TestMasterCompanyTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token:         "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token:         "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token:         token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/dmo/master",
			nil,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data master company & trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data master company & trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data master company & trader")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "companies", "list data master company & trader")
		assert.Contains(t, mapUnmarshal, "traders", "list data master company & trader")
	}
}
var idCompany = 0

// Company
func TestListCompany(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/company/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data company")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data company")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data company")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal,"companies")
	}
}

func TestCreateCompany(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"company_name": "PT. Integrata",
				"address":      "Narata Koplo, Japan 1231412",
				"province":     "Koplo",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"company_name": "PT. Maju Mundur",
				"address":      "Rasakana kan Ciledug",
				"province":     "",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"company_name": "",
				"address":      "Namgong Plateu",
				"province":     "Namgong",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
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
			"/api/v1/company/create",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create company")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create company")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create company")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idCompany = int(mapUnmarshal["ID"].(float64))
		}

		assert.Contains(t, mapUnmarshal, "ID", "create company")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create company")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create company")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create company")
		assert.Contains(t, mapUnmarshal, "company_name", "create company")
		assert.Contains(t, mapUnmarshal, "address", "create company")
		assert.Contains(t, mapUnmarshal, "province", "create company")
		assert.Contains(t, mapUnmarshal, "phone_number", "create company")
		assert.Contains(t, mapUnmarshal, "fax_number", "create company")
	}
}

func TestUpdateCompany(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		id	int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idCompany,
			body: fiber.Map{
				"company_name": "PT. Integrata",
				"address":      "Narata Koplo, Japan 1231412",
				"province":     "Koplo",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idCompany,
			body: fiber.Map{
				"company_name": "PT. Maju Mundur",
				"address":      "Rasakana kan Ciledug",
				"province":     "",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idCompany,
			body: fiber.Map{
				"company_name": "",
				"address":      "Namgong Plateu",
				"province":     "Namgong",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1002,
			body: fiber.Map{
				"company_name": "PT JUMAN",
				"address":      "Namgong Plateu",
				"province":     "Namgong",
				"phone_number": "",
				"fax_number":   "",
			},
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		url := fmt.Sprintf("/api/v1/company/update/%v", test.id)
		req, _ := http.NewRequest(
			"PUT",
			url,
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "update company")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update company")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update company")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		assert.Contains(t, mapUnmarshal, "ID", "update company")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update company")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update company")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update company")
		assert.Contains(t, mapUnmarshal, "company_name", "update company")
		assert.Contains(t, mapUnmarshal, "address", "update company")
		assert.Contains(t, mapUnmarshal, "province", "update company")
		assert.Contains(t, mapUnmarshal, "phone_number", "update company")
		assert.Contains(t, mapUnmarshal, "fax_number", "update company")
	}
}

func TestDetailCompany(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idCompany,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/company/detail/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "detail company")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "detail company")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "detail company")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "company", "detail company")
		assert.Contains(t, mapUnmarshal, "list_traders", "detail company")
	}
}

var idTrader = 0
// Trader
func TestListTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/trader/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data trader")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal,"traders")
	}
}

func TestCreateTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"trader_name": "Budi Arya",
				"position": "Procu",
				"company_id": idCompany,
				"email": nil,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			body: fiber.Map{
				"trader_name": "Budi Arya",
				"position": "Procu",
				"company_id": 1000,
				"email": nil,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"trader_name": "Budi Arya",
				"position": "Procu",
				"company_id": idCompany,
				"email": "abcsesd",
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"trader_name": "",
				"position": "",
				"company_id": idCompany,
				"email": nil,
			},
			token: token,
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
			"/api/v1/trader/create",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create trader")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idTrader = int(mapUnmarshal["ID"].(float64))
		}

		assert.Contains(t, mapUnmarshal, "ID", "create trader")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create trader")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create trader")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create trader")
		assert.Contains(t, mapUnmarshal, "trader_name", "create trader")
		assert.Contains(t, mapUnmarshal, "position", "create trader")
		assert.Contains(t, mapUnmarshal, "email", "create trader")
		assert.Contains(t, mapUnmarshal, "company_id", "create trader")
		assert.Contains(t, mapUnmarshal, "company", "create trader")
	}
}

func TestUpdateTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		id	int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idTrader,
			body: fiber.Map{
				"trader_name": "Budi Arya",
				"position": "Procu",
				"company_id": idCompany,
				"email": nil,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: idCompany,
			body: fiber.Map{
				"trader_name": "Budi Arya",
				"position": "Procu",
				"company_id": 1000,
				"email": nil,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			id: idCompany,
			body: fiber.Map{
				"trader_name": "",
				"position": "",
				"company_id": idCompany,
				"email": nil,
			},
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		url := fmt.Sprintf("/api/v1/trader/update/%v", test.id)
		req, _ := http.NewRequest(
			"PUT",
			url,
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "update trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update trader")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		assert.Contains(t, mapUnmarshal, "ID", "update trader")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update trader")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update trader")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update trader")
		assert.Contains(t, mapUnmarshal, "trader_name", "update trader")
		assert.Contains(t, mapUnmarshal, "position", "create trader")
		assert.Contains(t, mapUnmarshal, "email", "create trader")
		assert.Contains(t, mapUnmarshal, "company_id", "create trader")
		assert.Contains(t, mapUnmarshal, "company", "create trader")
	}
}

func TestDetailTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idTrader,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/trader/detail/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "detail trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "detail trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "detail trader")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "detail trader")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "detail trader")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "detail trader")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "detail trader")
		assert.Contains(t, mapUnmarshal, "trader_name", "detail trader")
		assert.Contains(t, mapUnmarshal, "position", "detail trader")
		assert.Contains(t, mapUnmarshal, "email", "detail trader")
		assert.Contains(t, mapUnmarshal, "company_id", "detail trader")
		assert.Contains(t, mapUnmarshal, "company", "detail trader")
	}
}

func TestDeleteTrader(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idTrader,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1000,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/trader/delete/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete trader")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete trader")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete trader")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete trader")
	}
}

// Delete Company after Delete Trader
func TestDeleteCompany(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idCompany,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1000,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/company/delete/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete company")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete company")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete company")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete company")
	}
}

// Notification

func TestCreateNotification(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"status" : "success create",
				"type" : "dmo",
				"period" : "Jan 2022",
			},
			token: token,
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
			"/api/v1/notification/create",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create notification")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create notification")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create notification")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		assert.Contains(t, mapUnmarshal, "ID", "create notification")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create notification")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create notification")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create notification")
		assert.Contains(t, mapUnmarshal, "status", "create notification")
		assert.Contains(t, mapUnmarshal, "type", "create notification")
		assert.Contains(t, mapUnmarshal, "period", "create notification")
		assert.Contains(t, mapUnmarshal, "user_id", "create notification")
	}
}

func TestGetNotification(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/notification/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data notification")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data notification")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data notification")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal,"list", "list data notification")
	}
}

func TestUpdateNotification(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"PUT",
			"/api/v1/notification/update",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list update data notification")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list update data notification")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list update data notification")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal,"list", "list update data notification")
	}
}

// Delete Dmo

func TestDeleteDmo(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idDmo,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/dmo/delete/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete data dmo")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete data dmo")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete data dmo")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete data dmo")
	}
}

// Report
func TestGetReportDetail(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/transaction/report/detail",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "report transaction detail")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "report transaction detail")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "report transaction detail")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal,"electricity", "report transaction detail")
		assert.Contains(t, mapUnmarshal,"non_electricity", "report transaction detail")
		assert.Contains(t, mapUnmarshal,"not_claimable", "report transaction detail")
	}
}

func TestGetReportRecap(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token:         "",
			body:          fiber.Map{},
		},
		{
			expectedError: false,
			expectedCode:  401,
			token:         "afwifiwgjwigjianveri",
			body:          fiber.Map{},
		},
		{
			expectedError: false,
			expectedCode:  200,
			token:         token,
			body: fiber.Map{
				"production_plan":                  4000000,
				"percentage_production_obligation": 25,
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
			"/api/v1/transaction/report/recap",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "report transaction recap")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "report transaction recap")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "report transaction recap")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "electricity_total", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "non_electricity_total", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "total", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "rate_calories", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "production_plan", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "production_obligation", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "percentage_production_obligation", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "prorate_production_plan", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "fulfillment_of_production_plan", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "fulfillment_of_production_realization", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "fulfillment_percentage_production_obligation", "report transaction recap")
		assert.Contains(t, mapUnmarshal, "year", "report transaction recap")
	}
}

var idProduction = 0
// Production
func TestCreateProduction(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  201,
			body: fiber.Map{
				"production_date": "2022-10-30",
				"quantity": 2524.242,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"production_date": "2022-10-30",
				"quantity": 0,
			},
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"production_date": "2022-15-10",
				"quantity": 125123,
			},
			token: token,
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
			"/api/v1/production/create",
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "create production")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "create production")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "create production")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		if res.StatusCode == 201 {
			idProduction = int(mapUnmarshal["ID"].(float64))
		}

		assert.Contains(t, mapUnmarshal, "ID", "create production")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "create production")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "create production")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "create production")
		assert.Contains(t, mapUnmarshal, "production_date", "create production")
		assert.Contains(t, mapUnmarshal, "quantity", "create production")
	}
}

func TestUpdateProduction(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token         string
		id	int
		body          map[string]interface{}
	}{
		{
			expectedError: false,
			expectedCode:  401,
			body:          fiber.Map{},
			id: idProduction,
			token:         "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			body: fiber.Map{
				"production_date": "2022-10-30",
				"quantity": 2562,
			},
			id: idProduction,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			body: fiber.Map{
				"production_date": "2022-10-30",
				"quantity": 3535,
			},
			id: 1000,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  400,
			body: fiber.Map{
				"production_date": "2022-15-10",
				"quantity": 125123,
			},
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		bodyJson, err := json.Marshal(test.body)
		var payload = bytes.NewBufferString(string(bodyJson))
		urlLink := fmt.Sprintf("/api/v1/production/update/%v", test.id)
		req, _ := http.NewRequest(
			"PUT",
			urlLink,
			payload,
		)

		req.Header.Add("Authorization", "Bearer "+test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "update production")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "update production")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "update production")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)

		if res.StatusCode >= 400 {
			continue
		}

		assert.Contains(t, mapUnmarshal, "ID", "update production")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "update production")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "update production")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "update production")
		assert.Contains(t, mapUnmarshal, "production_date", "update production")
		assert.Contains(t, mapUnmarshal, "quantity", "update production")
	}
}

func TestListProduction(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
	}{
		{
			expectedError: false,
			expectedCode:  401,
			token: "",
		},
		{
			expectedError: false,
			expectedCode:  401,
			token: "afwifiwgjwigjianveri",
		},
		{
			expectedError: false,
			expectedCode:  200,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			"/api/v1/production/list",
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "list data production")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "list data production")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "list data production")

		mapUnmarshal  := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "limit", "list data production")
		assert.Contains(t, mapUnmarshal, "page", "list data production")
		assert.Contains(t, mapUnmarshal, "total_rows", "list data production")
		assert.Contains(t, mapUnmarshal, "total_pages", "list data production")
		assert.Contains(t, mapUnmarshal, "data", "list data production")
	}
}

func TestDetailProduction(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idProduction,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/production/detail/%v", test.id)
		req, _ := http.NewRequest(
			"GET",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "detail production")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "detail production")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "detail production")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "ID", "detail production")
		assert.Contains(t, mapUnmarshal, "CreatedAt", "detail production")
		assert.Contains(t, mapUnmarshal, "UpdatedAt", "detail production")
		assert.Contains(t, mapUnmarshal, "DeletedAt", "detail production")
		assert.Contains(t, mapUnmarshal, "production_date", "detail production")
		assert.Contains(t, mapUnmarshal, "quantity", "detail production")
	}
}

func TestDeleteProduction(t *testing.T) {
	tests := []struct {
		expectedError bool
		expectedCode  int
		token string
		id int
	}{
		{
			expectedError: false,
			expectedCode:  401,
			id: 1,
			token: "asdawfaeac",
		},
		{
			expectedError: false,
			expectedCode:  404,
			id: 1050,
			token: token,
		},
		{
			expectedError: false,
			expectedCode:  200,
			id: idProduction,
			token: token,
		},
	}

	db, validate := startSetup()
	app := fiber.New()
	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	for _, test := range tests {
		url := fmt.Sprintf("/api/v1/production/delete/%v", test.id)
		req, _ := http.NewRequest(
			"DELETE",
			url,
			nil,
		)

		req.Header.Add("Authorization", "Bearer " + test.token)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, "delete data production")
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, "delete data production")

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Ensure that the body was read correctly
		assert.Nilf(t, err, "delete data production")

		mapUnmarshal := make(map[string]interface{})

		errUnmarshal := json.Unmarshal(body, &mapUnmarshal)

		fmt.Println(errUnmarshal)
		if res.StatusCode >= 400 {
			continue
		}

		//// Verify, that the reponse body equals the expected body
		assert.Contains(t, mapUnmarshal, "message", "delete data production")
	}
}
