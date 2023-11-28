package handler

import (
	"ajebackend/model/history"
	"ajebackend/model/jettybalance"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/pitloss"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type jettyBalanceHandler struct {
	pitLossService          pitloss.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	notificationUserService notificationuser.Service
	userIupopkService       useriupopk.Service
	jettyBalanceService     jettybalance.Service
}

func NewJettyBalanceHandler(pitLossService pitloss.Service, historyService history.Service, v *validator.Validate, logService logs.Service, notificationUserService notificationuser.Service, userIupopkService useriupopk.Service, jettyBalanceService jettybalance.Service) *jettyBalanceHandler {
	return &jettyBalanceHandler{
		pitLossService,
		historyService,
		v,
		logService,
		notificationUserService,
		userIupopkService,
		jettyBalanceService,
	}
}

func (h *jettyBalanceHandler) CreateJettyBalance(c *fiber.Ctx) error {
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

	inputJettyPitLoss := new(pitloss.InputJettyPitLoss)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputJettyPitLoss); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputJettyPitLoss)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["year"] = inputJettyPitLoss.Year
		inputMap["jetty_id"] = inputJettyPitLoss.JettyId
		inputMap["input"] = inputJettyPitLoss
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createJettyPitLoss, createJettyPitLossErr := h.historyService.CreateJettyBalance(*inputJettyPitLoss, iupopkIdInt, uint(claims["id"].(float64)))

	if createJettyPitLossErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["year"] = inputJettyPitLoss.Year
		inputMap["jetty_id"] = inputJettyPitLoss.JettyId
		inputMap["input"] = inputJettyPitLoss
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": createJettyPitLossErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": createJettyPitLossErr.Error(),
		})
	}

	return c.Status(201).JSON(createJettyPitLoss)
}

func (h *jettyBalanceHandler) DetailJettyBalance(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detail, detailErr := h.pitLossService.DetailJettyBalance(idInt, iupopkIdInt)

	if detailErr != nil {

		status := 400

		if detailErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	return c.Status(200).JSON(detail)
}

func (h *jettyBalanceHandler) UpdateJettyBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
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

	inputJettyPitLoss := new(pitloss.InputUpdateJettyPitLoss)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputJettyPitLoss); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputJettyPitLoss)
	idUint := uint(idInt)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})

		inputMap["input"] = inputJettyPitLoss
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			JettyBalanceId: &idUint,
			Input:          inputJson,
			Message:        messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateJettyBalance, updateJettyBalanceErr := h.historyService.UpdateJettyBalance(idInt, *inputJettyPitLoss, iupopkIdInt, uint(claims["id"].(float64)))

	if updateJettyBalanceErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = inputJettyPitLoss
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": updateJettyBalanceErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:          inputJson,
			Message:        messageJson,
			JettyBalanceId: &idUint,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": updateJettyBalanceErr.Error(),
		})
	}

	return c.Status(201).JSON(updateJettyBalance)
}

func (h *jettyBalanceHandler) ListJettyBalance(c *fiber.Ctx) error {
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

	var filterJettyBalance jettybalance.SortFilterJettyBalance

	filterJettyBalance.Field = c.Query("field")
	filterJettyBalance.Sort = c.Query("sort")
	filterJettyBalance.JettyId = c.Query("jetty_id")
	filterJettyBalance.Year = c.Query("year")

	listJettyBalance, listJettyBalanceErr := h.jettyBalanceService.ListJettyBalance(pageNumber, filterJettyBalance, iupopkIdInt)

	if listJettyBalanceErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listJettyBalanceErr.Error(),
		})
	}

	return c.Status(200).JSON(listJettyBalance)
}

func (h *jettyBalanceHandler) DeleteJettyBalance(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteJettyBalanceErr := h.historyService.DeleteJettyBalance(idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if deleteJettyBalanceErr != nil {

		status := 400

		if deleteJettyBalanceErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": deleteJettyBalanceErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete jetty balance",
	})
}
