package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/insw"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type inswHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	groupingVesselLnService groupingvesselln.Service
	notificationUserService notificationuser.Service
}

func NewInswHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, v *validator.Validate, logService logs.Service, groupingVesselLnService groupingvesselln.Service, notificationUserService notificationuser.Service) *inswHandler {
	return &inswHandler{
		transactionService,
		userService,
		historyService,
		v,
		logService,
		groupingVesselLnService,
		notificationUserService,
	}
}

func (h *inswHandler) CreateInsw(c *fiber.Ctx) error {
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

	inputCreateInsw := new(insw.InputCreateInsw)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateInsw); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateInsw)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_month"] = inputCreateInsw.Month
		inputMap["insw_year"] = inputCreateInsw.Year
		inputMap["list_grouping_vessel_ln"] = inputCreateInsw.ListGroupingVesselId
		inputMap["input"] = inputCreateInsw
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

	baseIdNumber := fmt.Sprintf("INSW-%s-%v", helper.MonthLongToNumberString(inputCreateInsw.Month), inputCreateInsw.Year)
	createInsw, createInswErr := h.historyService.CreateInsw(inputCreateInsw.ListGroupingVesselId, inputCreateInsw.Month, inputCreateInsw.Year, baseIdNumber, uint(claims["id"].(float64)))

	if createInswErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_month"] = inputCreateInsw.Month
		inputMap["insw_year"] = inputCreateInsw.Year
		inputMap["list_grouping_vessel_ln"] = inputCreateInsw.ListGroupingVesselId
		inputMap["input"] = inputCreateInsw
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createInswErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createInswErr.Error(),
		})
	}

	return c.Status(201).JSON(createInsw)
}

func (h *inswHandler) ListGroupingVesselWithoutInsw(c *fiber.Ctx) error {
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

	listGroupingVesselWithoutInsw, listGroupingVesselWithoutInswErr := h.groupingVesselLnService.ListGroupingVesselLnWithoutInsw()

	if listGroupingVesselWithoutInswErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listGroupingVesselWithoutInswErr.Error(),
		})
	}

	return c.Status(200).JSON(listGroupingVesselWithoutInsw)
}

func (h *inswHandler) DetailInsw(c *fiber.Ctx) error {
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

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt)

	if detailInswErr != nil {

		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	return c.Status(200).JSON(detailInsw)
}

func (h *inswHandler) DeleteInsw(c *fiber.Ctx) error {
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

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt)

	if detailInswErr != nil {

		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	_, deleteInswErr := h.historyService.DeleteInsw(idInt, uint(claims["id"].(float64)))

	if deleteInswErr != nil {

		status := 400

		if deleteInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": deleteInswErr.Error(),
		})
	}

	if detailInsw.Detail.InswDocumentLink != "" {

	}

	if detailInsw.Detail.InswDocumentLink != "" {
		fileName := fmt.Sprintf("INSW/%s/", *detailInsw.Detail.IdNumber)
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["insw_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			inswId := uint(idInt)
			createdErrLog := logs.Logs{
				InswId:  &inswId,
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if deleteAwsErr.Error() == "record not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"message": "failed to delete minerba ln aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	var createNotif notification.InputNotification

	createNotif.Type = "insw"
	createNotif.Status = "menghapus"
	createNotif.Period = fmt.Sprintf("%s %v", detailInsw.Detail.Month, detailInsw.Detail.Year)

	_, createNotificationDeleteInsw := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationDeleteInsw != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_id"] = idInt
		inputMap["month"] = detailInsw.Detail.Month
		inputMap["year"] = detailInsw.Detail.Year
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteInsw.Error(),
		})
		inswId := uint(idInt)
		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			InswId:  &inswId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteInsw.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete insw",
	})
}
