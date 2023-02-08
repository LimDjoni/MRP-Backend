package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerbaln"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
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

type minerbaLnHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	logService              logs.Service
	minerbaLnService        minerbaln.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}

func NewMinerbaLnHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, minerbaLnService minerbaln.Service, notificationUserService notificationuser.Service, v *validator.Validate) *minerbaLnHandler {
	return &minerbaLnHandler{
		transactionService,
		userService,
		historyService,
		logService,
		minerbaLnService,
		notificationUserService,
		v,
	}
}

func (h *minerbaLnHandler) CreateMinerbaLn(c *fiber.Ctx) error {
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

	inputCreateMinerbaLn := new(minerbaln.InputCreateMinerbaLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateMinerbaLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateMinerbaLn)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_period"] = inputCreateMinerbaLn.Period
		inputMap["list_transactions"] = inputCreateMinerbaLn.ListDataLn
		inputMap["input"] = inputCreateMinerbaLn
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

	_, findMinerbaErr := h.minerbaLnService.GetReportMinerbaLnWithPeriod(inputCreateMinerbaLn.Period)

	if findMinerbaErr == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "report with same period already exist",
		})
	}

	splitPeriod := strings.Split(inputCreateMinerbaLn.Period, " ")

	baseIdNumber := fmt.Sprintf("LSL-%s-%s", helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1])
	createMinerbaLn, createMinerbaLnErr := h.historyService.CreateMinerbaLn(inputCreateMinerbaLn.Period, baseIdNumber, inputCreateMinerbaLn.ListDataLn, uint(claims["id"].(float64)))

	if createMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_period"] = inputCreateMinerbaLn.Period
		inputMap["list_dn"] = inputCreateMinerbaLn.ListDataLn

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createMinerbaLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createMinerbaLnErr.Error(),
		})
	}

	return c.Status(201).JSON(createMinerbaLn)
}

func (h *minerbaLnHandler) UpdateMinerbaLn(c *fiber.Ctx) error {
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

	inputUpdateMinerbaLn := new(minerbaln.InputUpdateMinerbaLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateMinerbaLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	idInt, err := strconv.Atoi(id)

	minerbaLnId := uint(idInt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update minerba ln",
			"error":   "record not found",
		})
	}

	errors := h.v.Struct(*inputUpdateMinerbaLn)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputMap["list_ln"] = inputUpdateMinerbaLn.ListDataLn
		inputMap["input"] = inputUpdateMinerbaLn
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			MinerbaLnId: &minerbaLnId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	_, findDetailMinerbaLnErr := h.transactionService.GetDetailMinerbaLn(idInt)

	if findDetailMinerbaLnErr != nil {
		status := 400

		if findDetailMinerbaLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": findDetailMinerbaLnErr.Error(),
		})
	}

	_, checkMinerbaLnUpdateTransactionErr := h.transactionService.CheckDataLnAndMinerbaLnUpdate(inputUpdateMinerbaLn.ListDataLn, idInt)

	if checkMinerbaLnUpdateTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputMap["list_ln"] = inputUpdateMinerbaLn.ListDataLn
		inputMap["input"] = inputUpdateMinerbaLn
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": checkMinerbaLnUpdateTransactionErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if checkMinerbaLnUpdateTransactionErr.Error() == "please check there is transaction not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": checkMinerbaLnUpdateTransactionErr.Error(),
		})
	}

	updateMinerbaLn, updateMinerbaLnErr := h.historyService.UpdateMinerbaLn(idInt, inputUpdateMinerbaLn.ListDataLn, uint(claims["id"].(float64)))

	if updateMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputMap["list_ln"] = inputUpdateMinerbaLn.ListDataLn
		inputMap["input"] = inputUpdateMinerbaLn
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateMinerbaLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateMinerbaLnErr.Error(),
		})
	}

	var createNotif notification.InputNotification

	createNotif.Type = "minerba ln"
	createNotif.Status = "mengedit"
	createNotif.Period = updateMinerbaLn.Period

	_, createNotificationUpdateMinerbaLnErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationUpdateMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputMap["list_ln"] = inputUpdateMinerbaLn.ListDataLn
		inputMap["input"] = inputUpdateMinerbaLn
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationUpdateMinerbaLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationUpdateMinerbaLnErr.Error(),
		})
	}

	return c.Status(200).JSON(updateMinerbaLn)
}

func (h *minerbaLnHandler) DeleteMinerbaLn(c *fiber.Ctx) error {
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
			"message": "failed to delete minerba",
			"error":   "record not found",
		})
	}

	dataMinerbaLn, getDataMinerbaLnErr := h.minerbaLnService.GetDataMinerbaLn(idInt)

	if getDataMinerbaLnErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete minerba ln",
			"error":   "record not found",
		})
	}

	_, deleteMinerbaLnErr := h.historyService.DeleteMinerbaLn(idInt, uint(claims["id"].(float64)))

	if deleteMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteMinerbaLnErr.Error(),
		})

		minerbaLnId := uint(idInt)
		createdErrLog := logs.Logs{
			MinerbaLnId: &minerbaLnId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteMinerbaLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete minerba",
			"error":   deleteMinerbaLnErr.Error(),
		})
	}

	if dataMinerbaLn.SP3MELNDocumentLink != nil {
		fileName := fmt.Sprintf("LSL/%s/", *dataMinerbaLn.IdNumber)
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["minerba_ln_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			minerbaLnId := uint(idInt)
			createdErrLog := logs.Logs{
				MinerbaLnId: &minerbaLnId,
				Input:       inputJson,
				Message:     messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if deleteMinerbaLnErr.Error() == "record not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"message": "failed to delete minerba ln aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	var createNotif notification.InputNotification

	createNotif.Type = "minerba ln"
	createNotif.Status = "menghapus"
	createNotif.Period = dataMinerbaLn.Period

	_, createNotificationDeleteMinerbaLnErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationDeleteMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_id"] = idInt
		inputMap["minerba_ln_period"] = dataMinerbaLn.Period
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteMinerbaLnErr.Error(),
		})
		minerbaLnId := uint(idInt)
		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteMinerbaLnErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete minerba ln",
	})
}

func (h *minerbaLnHandler) ListDataLNWithoutMinerba(c *fiber.Ctx) error {
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

	listDataLNWithoutMinerba, listDataLNWithoutMinerbaErr := h.transactionService.ListDataLNWithoutMinerba()

	if listDataLNWithoutMinerbaErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDataLNWithoutMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": listDataLNWithoutMinerba,
	})
}

func (h *minerbaLnHandler) ListMinerbaLn(c *fiber.Ctx) error {
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

	var filterMinerbaLn minerbaln.FilterAndSortMinerbaLn

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		filterMinerbaLn.Quantity = 0
	} else {
		filterMinerbaLn.Quantity = quantity
	}

	filterMinerbaLn.CreatedStart = c.Query("created_start")
	filterMinerbaLn.CreatedEnd = c.Query("created_end")
	filterMinerbaLn.Field = c.Query("field")
	filterMinerbaLn.Sort = c.Query("sort")

	listMinerbaLn, listMinerbaLnErr := h.minerbaLnService.GetListReportMinerbaLnAll(pageNumber, filterMinerbaLn)

	if listMinerbaLnErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listMinerbaLnErr.Error(),
		})
	}

	return c.Status(200).JSON(listMinerbaLn)
}

func (h *minerbaLnHandler) DetailMinerbaLn(c *fiber.Ctx) error {
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

	detailMinerbaLn, detailMinerbaLnErr := h.transactionService.GetDetailMinerbaLn(idInt)

	if detailMinerbaLnErr != nil {
		status := 400

		if detailMinerbaLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaLnErr.Error(),
		})
	}

	return c.Status(200).JSON(detailMinerbaLn)
}

func (h *minerbaLnHandler) CheckValidPeriodMinerbaLn(c *fiber.Ctx) error {
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

	inputCheckMinerbaLn := new(minerbaln.CheckMinerbaLnPeriod)

	if err := c.BodyParser(inputCheckMinerbaLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCheckMinerbaLn)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid period",
			"errors":  dataErrors,
		})
	}

	_, detailMinerbaLnErr := h.minerbaLnService.GetReportMinerbaLnWithPeriod(inputCheckMinerbaLn.Period)

	if detailMinerbaLnErr != nil && detailMinerbaLnErr.Error() == "record not found" {
		return c.Status(200).JSON(fiber.Map{
			"message": "valid period",
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"message": "invalid period",
	})
}

func (h *minerbaLnHandler) UpdateDocumentMinerbaLn(c *fiber.Ctx) error {
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

	inputUpdateMinerbaLn := new(minerbaln.InputUpdateDocumentMinerbaLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateMinerbaLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update minerba ln",
		})
	}

	errors := h.v.Struct(*inputUpdateMinerbaLn)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update minerba",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerbaLn
		inputMap["minerba_ln_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		minerbaLnId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailMinerbaLn, detailMinerbaLnErr := h.transactionService.GetDetailMinerbaLn(idInt)

	if detailMinerbaLnErr != nil {
		status := 400

		if detailMinerbaLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaLnErr.Error(),
		})
	}

	updateMinerbaLn, updateMinerbaLnErr := h.historyService.UpdateDocumentMinerbaLn(idInt, *inputUpdateMinerbaLn, uint(claims["id"].(float64)))

	if updateMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerbaLn

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateMinerbaLnErr.Error(),
		})

		minerbaLnId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			MinerbaLnId: &minerbaLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateMinerbaLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateMinerbaLnErr.Error(),
			"message": "failed to update minerba ln",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "minerba ln"
	inputNotification.Status = "membuat"
	inputNotification.Period = detailMinerbaLn.Detail.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerbaLn

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update minerba",
		})
	}

	return c.Status(200).JSON(updateMinerbaLn)
}

func (h *minerbaLnHandler) RequestCreateExcelMinerbaLn(c *fiber.Ctx) error {
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

	detailMinerbaLn, detailMinerbaLnErr := h.transactionService.GetDetailMinerbaLn(idInt)

	if detailMinerbaLnErr != nil {
		status := 400

		if detailMinerbaLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaLnErr.Error(),
		})
	}

	var inputRequestCreateExcel transaction.InputRequestCreateExcelMinerba
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.MinerbaId = idInt
	inputRequestCreateExcel.MinerbaNumber = *detailMinerbaLn.Detail.IdNumber
	inputRequestCreateExcel.MinerbaPeriod = detailMinerbaLn.Detail.Period
	inputRequestCreateExcel.Transactions = detailMinerbaLn.List

	hitJob, hitJobErr := h.transactionService.RequestCreateExcelLn(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}
