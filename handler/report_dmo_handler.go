package handler

import (
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/reportdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type reportDmoHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	logService              logs.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
	reportDmoService        reportdmo.ReportDmo
}

func NewReportDmoHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, notificationUserService notificationuser.Service, v *validator.Validate, reportDmoService reportdmo.ReportDmo) *reportDmoHandler {
	return &reportDmoHandler{
		transactionService,
		userService,
		historyService,
		logService,
		notificationUserService,
		v,
		reportDmoService,
	}
}

func (h *reportDmoHandler) CreateReportDmo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateReportDmo := new(reportdmo.InputCreateReportDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateReportDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateReportDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_period"] = inputCreateReportDmo.Period
		inputMap["transactions"] = inputCreateReportDmo.Transactions
		inputMap["grouping_vessel"] = inputCreateReportDmo.GroupingVessels
		inputMap["input"] = inputCreateReportDmo
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

	splitPeriod := strings.Split(inputCreateReportDmo.Period, " ")

	baseIdNumber := fmt.Sprintf("LDO-%s-%s", helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1])

	createReportDmo, createReportDmoErr := h.historyService.CreateReportDmo(*&inputCreateReportDmo, baseIdNumber, uint(claims["id"].(float64)))

	if createReportDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_period"] = inputCreateReportDmo.Period
		inputMap["transactions"] = inputCreateReportDmo.Transactions
		inputMap["grouping_vessel"] = inputCreateReportDmo.GroupingVessels
		inputMap["input"] = inputCreateReportDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createReportDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createReportDmoErr.Error(),
		})
	}

	return c.Status(201).JSON(createReportDmo)
}

func (h *reportDmoHandler) UpdateDocumentReportDmo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputUpdateDocumentReportDmo := new(reportdmo.InputUpdateDocumentReportDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateDocumentReportDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update minerba",
		})
	}

	errors := h.v.Struct(*inputUpdateDocumentReportDmo)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update report dmo",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDocumentReportDmo
		inputMap["report_dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		reportDmoId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			ReportDmoId: &reportDmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateReportDmo, updateReportDmoErr := h.historyService.UpdateDocumentReportDmo(idInt, *inputUpdateDocumentReportDmo, uint(claims["id"].(float64)))

	if updateReportDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDocumentReportDmo

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateReportDmoErr.Error(),
		})
		reportDmoId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			ReportDmoId: &reportDmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateReportDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateReportDmoErr.Error(),
			"message": "failed to update minerba",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "report dmo"
	inputNotification.Status = "membuat"
	inputNotification.Period = updateReportDmo.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDocumentReportDmo

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		reportDmoId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			ReportDmoId: &reportDmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update report dmo",
		})
	}

	return c.Status(200).JSON(updateReportDmo)
}

func (h *reportDmoHandler) RequestCreateExcelReportDmo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
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

	header := c.GetReqHeaders()

	detailReportDmo, detailReportDmoErr := h.transactionService.GetDetailReportDmo(idInt)

	if detailReportDmoErr != nil {
		status := 400

		if detailReportDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailReportDmoErr.Error(),
		})
	}

	var inputRequestCreateReportDmo transaction.InputRequestCreateReportDmo
	inputRequestCreateReportDmo.Authorization = header["Authorization"]
	inputRequestCreateReportDmo.ReportDmo = detailReportDmo.Detail
	inputRequestCreateReportDmo.GroupingVessels = detailReportDmo.Transactions
	inputRequestCreateReportDmo.Transactions = detailReportDmo.GroupingVessels

	hitJob, hitJobErr := h.transactionService.RequestCreateReportDmo(inputRequestCreateReportDmo)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}
