package handler

import (
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/production"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"strconv"
)

type productionHandler struct {
	userService user.Service
	historyService history.Service
	productionService production.Service
	logsService logs.Service
	v *validator.Validate
}

func NewProductionHandler(userService user.Service, historyService history.Service, productionService production.Service, logsService logs.Service, v *validator.Validate) *productionHandler {
	return &productionHandler{
		userService,
		historyService,
		productionService,
		logsService,
		v,
	}
}

func (h *productionHandler) ListProduction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page == "" {
		pageNumber = 1
	}

	var filterProduction production.FilterListProduction

	filterProduction.ProductionDateStart = c.Query("production_date_start")
	filterProduction.ProductionDateEnd = c.Query("production_date_end")
	filterProduction.Field = c.Query("field")
	filterProduction.Sort = c.Query("sort")

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		filterProduction.Quantity = 0
	} else {
		filterProduction.Quantity = quantity
	}

	listProduction, listProductionErr := h.productionService.GetListProduction(pageNumber, filterProduction)

	if listProductionErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listProductionErr.Error(),
		})
	}

	return c.Status(200).JSON(listProduction)
}

func (h *productionHandler) CreateProduction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	productionInput := new(production.InputCreateProduction)

	if err := c.BodyParser(productionInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*productionInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createProduction, createProductionErr := h.historyService.CreateProduction(*productionInput, uint(claims["id"].(float64)))

	if createProductionErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createProductionErr.Error(),
		})
	}

	return c.Status(201).JSON(createProduction)
}

func (h *productionHandler) UpdateProduction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	productionUpdateInput := new(production.InputCreateProduction)

	if err := c.BodyParser(productionUpdateInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*productionUpdateInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	updateProduction, updateProductionErr := h.historyService.UpdateProduction(*productionUpdateInput, idInt, uint(claims["id"].(float64)))

	if updateProductionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = productionUpdateInput
		inputMap["user_id"] = claims["id"]
		inputMap["production_id"] = id

		inputJson , _ := json.Marshal(productionUpdateInput)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateProductionErr.Error(),
		})

		productionId := uint(idInt)
		createdErrLog := logs.Logs{
			ProductionId: &productionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if updateProductionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": updateProductionErr.Error(),
		})
	}

	return c.Status(200).JSON(updateProduction)
}

func (h *productionHandler) DeleteProduction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteTransactionErr := h.historyService.DeleteProduction(idInt, uint(claims["id"].(float64)))

	if deleteTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["production_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteTransactionErr.Error(),
		})

		productionId := uint(idInt)
		createdErrLog := logs.Logs{
			ProductionId: &productionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if  deleteTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete production",
			"error": deleteTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete production",
	})
}

func (h *productionHandler) DetailProduction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detailProduction, detailProductionErr := h.productionService.DetailProduction(idInt)

	if detailProductionErr != nil {
		status := 400

		if detailProductionErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailProductionErr.Error(),
		})
	}

	return c.Status(200).JSON(detailProduction)
}
