package handler

import (
	"ajebackend/model/logs"
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type traderHandler struct {
	traderService     trader.Service
	companyService    company.Service
	traderDmoService  traderdmo.Service
	logsService       logs.Service
	v                 *validator.Validate
	userIupopkService useriupopk.Service
}

func NewTraderHandler(traderService trader.Service, companyService company.Service, traderDmoService traderdmo.Service, logsService logs.Service, v *validator.Validate, userIupopkService useriupopk.Service) *traderHandler {
	return &traderHandler{
		traderService,
		companyService,
		traderDmoService,
		logsService,
		v,
		userIupopkService,
	}
}

func (h *traderHandler) ListTrader(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	listTrader, listTraderErr := h.traderService.ListTrader()

	if listTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listTraderErr.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"traders": listTrader,
	})
}

func (h *traderHandler) CreateTrader(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateTrader := new(trader.InputCreateUpdateTrader)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateTrader); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateTrader)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateTrader
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors":  dataErrors,
			"message": "failed to create trader",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createTrader, createTraderErr := h.traderService.CreateTrader(*inputCreateTrader)

	if createTraderErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateTrader
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":   createTraderErr.Error(),
			"message": "failed to create trader",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		messageErr := createTraderErr.Error()
		if createTraderErr.Error() == "record not found" {
			messageErr = "company " + createTraderErr.Error()
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": messageErr,
		})
	}

	return c.Status(201).JSON(createTrader)
}

func (h *traderHandler) UpdateTrader(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

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

	inputTraderUpdate := new(trader.InputCreateUpdateTrader)

	// Binds the request body to the Person struct
	if errParsing := c.BodyParser(inputTraderUpdate); errParsing != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errParsing.Error(),
		})
	}

	errors := h.v.Struct(*inputTraderUpdate)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	_, findTraderErr := h.traderService.DetailTrader(idInt)

	if findTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "trader record not found",
		})
	}

	_, findCompanyErr := h.companyService.DetailCompany(inputTraderUpdate.CompanyId)

	if findCompanyErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "company record not found",
		})
	}

	updateTrader, updateTraderErr := h.traderService.UpdateTrader(*inputTraderUpdate, idInt)

	if updateTraderErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = inputTraderUpdate
		inputMap["user_id"] = claims["id"]
		inputMap["trader_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateTraderErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if updateTraderErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update trader",
			"error":   updateTraderErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTrader)
}

func (h *traderHandler) DeleteTrader(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "record not found",
			"message": "failed to delete trader",
		})
	}

	_, findTraderErr := h.traderService.DetailTrader(idInt)

	if findTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   findTraderErr.Error(),
			"message": "failed to delete trader",
		})
	}

	listDmoWithTraderId, listDmoWithTraderIdErr := h.traderDmoService.DmoIdListWithTraderId(idInt)

	if findTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   listDmoWithTraderIdErr.Error(),
			"message": "failed to delete trader",
		})
	}

	if len(listDmoWithTraderId) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "trader is already used in dmo",
			"message": "failed to delete trader",
		})
	}

	_, deleteTraderErr := h.traderService.DeleteTrader(idInt)

	if deleteTraderErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["trader_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteTraderErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		status := 400

		if deleteTraderErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete trader",
			"error":   deleteTraderErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete trader",
	})
}

func (h *traderHandler) DetailTrader(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

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

	detailTrader, detailTraderErr := h.traderService.DetailTrader(idInt)

	if detailTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailTraderErr.Error(),
		})
	}

	return c.Status(200).JSON(detailTrader)
}
