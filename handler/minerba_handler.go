package handler

import (
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"strconv"
)

type minerbaHandler struct {
	transactionService transaction.Service
	userService user.Service
	historyService history.Service
	logService logs.Service
	minerbaService minerba.Service
}

func NewMinerbaHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, minerbaService minerba.Service) *minerbaHandler {
	return &minerbaHandler{
		transactionService,
		userService,
		historyService,
		logService,
		minerbaService,
	}
}

func (h *minerbaHandler) ListDataDNWithoutMinerba(c *fiber.Ctx) error {
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

	listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr := h.transactionService.ListDataDNWithoutMinerba()

	if listDataDNWithoutMinerbaErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDataDNWithoutMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": listDataDNWithoutMinerba,
	})
}

func (h *minerbaHandler) CreateMinerba(c *fiber.Ctx) error {
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

	inputCreateMinerba := new(minerba.InputCreateMinerba)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateMinerba); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, findMinerbaErr := h.minerbaService.GetReportMinerbaWithPeriod(inputCreateMinerba.Period)

	if findMinerbaErr == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "report with same period already exist",
		})
	}

	_, checkMinerbaTransactionErr := h.transactionService.CheckDataDNAndMinerba(inputCreateMinerba.ListDataDn)

	if checkMinerbaTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_period"] = inputCreateMinerba.Period
		inputMap["list_dn"] = inputCreateMinerba.ListDataDn

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": checkMinerbaTransactionErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": checkMinerbaTransactionErr.Error(),
		})
	}

	createMinerba, createMinerbaErr := h.historyService.CreateMinerba(inputCreateMinerba.Period, inputCreateMinerba.Period, inputCreateMinerba.ListDataDn, uint(claims["id"].(float64)))

	if createMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_period"] = inputCreateMinerba.Period
		inputMap["list_dn"] = inputCreateMinerba.ListDataDn

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createMinerbaErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createMinerbaErr.Error(),
		})
	}

	return c.Status(201).JSON(createMinerba)
}

func (h *minerbaHandler) DeleteMinerba(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete minerba",
			"error": "record not found",
		})
	}

	_, getDataMinerbaErr := h.minerbaService.GetDataMinerba(idInt)

	if getDataMinerbaErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete minerba",
			"error": "record not found",
		})
	}

	_, deleteMinerbaErr := h.historyService.DeleteMinerba(idInt, uint(claims["id"].(float64)))

	if deleteMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteMinerbaErr.Error(),
		})

		minerbaId := uint(idInt)
		createdErrLog := logs.Logs{
			MinerbaId: &minerbaId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if  deleteMinerbaErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete minerba",
			"error": deleteMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete minerba",
	})
}
