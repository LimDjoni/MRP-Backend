package handler

import (
	"ajebackend/model/company"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"strconv"
)

type companyHandler struct {
	userService user.Service
	companyService company.Service
	traderService trader.Service
	logsService logs.Service
	v *validator.Validate
}

func NewCompanyHandler(userService user.Service, companyService company.Service, traderService trader.Service, logsService logs.Service, v *validator.Validate) *companyHandler {
	return &companyHandler{
		userService,
		companyService,
		traderService,
		logsService,
		v,
	}
}

func (h *companyHandler) ListCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	_, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil {
		return c.Status(401).JSON(responseUnauthorized)
	}

	listCompany, listCompanyErr := h.companyService.ListCompany()

	if listCompanyErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listCompanyErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"companies": listCompany,
	})
}

func (h *companyHandler) CreateCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	_, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateCompany := new(company.InputCreateUpdateCompany)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateCompany); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateCompany)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCompany
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
			"message": "failed to create company",
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createCompany, createCompanyErr := h.companyService.CreateCompany(*inputCreateCompany)

	if createCompanyErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCompany
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": createCompanyErr,
			"message": "failed to create company",
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": createCompanyErr,
		})
	}

	return c.Status(201).JSON(createCompany)
}

func (h *companyHandler) UpdateCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	_, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	inputCompanyUpdate := new(company.InputCreateUpdateCompany)

	// Binds the request body to the Person struct
	if errParsing := c.BodyParser(inputCompanyUpdate); errParsing != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errParsing.Error(),
		})
	}

	errors := h.v.Struct(*inputCompanyUpdate)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateCompany, updateCompanyErr := h.companyService.UpdateCompany(*inputCompanyUpdate, idInt)

	if updateCompanyErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = inputCompanyUpdate
		inputMap["user_id"] = claims["id"]
		inputMap["company_id"] = id

		inputJson , _ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateCompanyErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if  updateCompanyErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update company",
			"error": updateCompanyErr.Error(),
		})
	}

	return c.Status(200).JSON(updateCompany)
}

func (h *companyHandler) DeleteCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	_, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, findCompanyErr := h.companyService.DetailCompany(idInt)

	if findCompanyErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": findCompanyErr.Error(),
		})
	}

	_, deleteCompanyErr := h.companyService.DeleteCompany(idInt)

	if deleteCompanyErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["company_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteCompanyErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if  deleteCompanyErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete company",
			"error": deleteCompanyErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete company",
	})
}

func (h *companyHandler) DetailCompany(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	_, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	var outputDetail trader.OutputCompanyDetail
	detailCompany, detailCompanyErr := h.companyService.DetailCompany(idInt)

	if detailCompanyErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailCompanyErr.Error(),
		})
	}

	outputDetail.Company = detailCompany

	listTraderWithCompany, listTraderWithCompanyErr := h.traderService.ListTraderWithCompanyId(idInt)

	if listTraderWithCompanyErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listTraderWithCompanyErr.Error(),
		})
	}

	outputDetail.ListTraders = listTraderWithCompany

	return c.Status(200).JSON(outputDetail)
}
