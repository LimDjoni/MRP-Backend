package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"strconv"
	"strings"
)

type minerbaHandler struct {
	transactionService transaction.Service
	userService user.Service
	historyService history.Service
	logService logs.Service
	minerbaService minerba.Service
	notificationUserService notificationuser.Service
	v *validator.Validate
}

func NewMinerbaHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, minerbaService minerba.Service, notificationUserService notificationuser.Service, v *validator.Validate) *minerbaHandler {
	return &minerbaHandler{
		transactionService,
		userService,
		historyService,
		logService,
		minerbaService,
		notificationUserService,
		v,
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateMinerba := new(minerba.InputCreateMinerba)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateMinerba); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateMinerba)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_period"] = inputCreateMinerba.Period
		inputMap["list_dn"] = inputCreateMinerba.ListDataDn
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	_, findMinerbaErr := h.minerbaService.GetReportMinerbaWithPeriod(inputCreateMinerba.Period)

	if findMinerbaErr == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "report with same period already exist",
		})
	}

	_, checkMinerbaTransactionErr := h.transactionService.CheckDataDnAndMinerba(inputCreateMinerba.ListDataDn)

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

		status := 400

		if checkMinerbaTransactionErr.Error() == "please check there is transaction not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": checkMinerbaTransactionErr.Error(),
		})
	}

	splitPeriod := strings.Split(inputCreateMinerba.Period, " ")

	baseIdNumber := fmt.Sprintf("LM-%s-%s",  helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1])
	createMinerba, createMinerbaErr := h.historyService.CreateMinerba(inputCreateMinerba.Period, baseIdNumber, inputCreateMinerba.ListDataDn, uint(claims["id"].(float64)))

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

func (h *minerbaHandler) UpdateMinerba(c *fiber.Ctx) error {
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

	id := c.Params("id")

	inputUpdateMinerba := new(minerba.InputUpdateMinerba)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateMinerba); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	idInt, err := strconv.Atoi(id)

	minerbaId := uint(idInt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update minerba",
			"error": "record not found",
		})
	}

	errors := h.v.Struct(*inputUpdateMinerba)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt
		inputMap["list_dn"] = inputUpdateMinerba.ListDataDn
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			MinerbaId: &minerbaId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	_, findDetailMinerbaErr := h.transactionService.GetDetailMinerba(idInt)

	if findDetailMinerbaErr != nil {
		status := 400

		if  findDetailMinerbaErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": findDetailMinerbaErr.Error(),
		})
	}

	_, checkMinerbaUpdateTransactionErr := h.transactionService.CheckDataDnAndMinerbaUpdate(inputUpdateMinerba.ListDataDn, idInt)

	if checkMinerbaUpdateTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt
		inputMap["list_dn"] = inputUpdateMinerba.ListDataDn

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": checkMinerbaUpdateTransactionErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if checkMinerbaUpdateTransactionErr.Error() == "please check there is transaction not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": checkMinerbaUpdateTransactionErr.Error(),
		})
	}

	updateMinerba, updateMinerbaErr := h.historyService.UpdateMinerba(idInt, inputUpdateMinerba.ListDataDn, uint(claims["id"].(float64)))

	if updateMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt
		inputMap["list_dn"] = inputUpdateMinerba.ListDataDn

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateMinerbaErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateMinerbaErr.Error(),
		})
	}

	var createNotif notification.InputNotification

	createNotif.Type = "minerba"
	createNotif.Status = "success update minerba"
	createNotif.Period = updateMinerba.Period

	_, createNotificationUpdateMinerbaErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationUpdateMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt
		inputMap["list_dn"] = inputUpdateMinerba.ListDataDn

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createNotificationUpdateMinerbaErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationUpdateMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(updateMinerba)
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
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

	dataMinerba, getDataMinerbaErr := h.minerbaService.GetDataMinerba(idInt)

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

	if dataMinerba.RecapDmoDocumentLink != nil || dataMinerba.DetailDmoDocumentLink != nil || dataMinerba.SP3MEDNDocumentLink != nil || dataMinerba.SP3MELNDocumentLink != nil || dataMinerba.INSWExportDocumentLink != nil {
		fileName := fmt.Sprintf("%s/", *dataMinerba.IdNumber)
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["minerba_id"] = idInt

			inputJson ,_ := json.Marshal(inputMap)
			messageJson ,_ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
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
				"message": "failed to delete minerba aws",
				"error": deleteAwsErr.Error(),
			})
		}

	}

	var createNotif notification.InputNotification

	createNotif.Type = "minerba"
	createNotif.Status = "success delete minerba"
	createNotif.Period = dataMinerba.Period

	_, createNotificationDeleteMinerbaErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationDeleteMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_id"] = idInt
		inputMap["minerba_period"] = dataMinerba.Period
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteMinerbaErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete minerba",
	})
}

func (h *minerbaHandler) UpdateDocumentMinerba(c *fiber.Ctx) error {
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

	inputUpdateMinerba := new(minerba.InputUpdateDocumentMinerba)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateMinerba); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
			"message": "failed to update minerba",
		})
	}

	errors := h.v.Struct(*inputUpdateMinerba)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update minerba",
			"error": "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerba
		inputMap["minerba_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		minerbaId := uint(idInt)

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			MinerbaId: &minerbaId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailMinerba, detailMinerbaErr := h.transactionService.GetDetailMinerba(idInt)

	if detailMinerbaErr != nil {
		status := 400

		if  detailMinerbaErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaErr.Error(),
		})
	}

	updateMinerba, updateMinerbaErr := h.historyService.UpdateDocumentMinerba(idInt, *inputUpdateMinerba, uint(claims["id"].(float64)))

	if updateMinerbaErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerba

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateMinerbaErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateMinerbaErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": updateMinerbaErr.Error(),
			"message": "failed to update minerba",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "minerba"
	inputNotification.Status = "success create document"
	inputNotification.Period = detailMinerba.Detail.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateMinerba

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createdNotificationErr.Error(),
			"message": "failed to create notification update minerba",
		})
	}

	return c.Status(200).JSON(updateMinerba)
}

func (h *minerbaHandler) ListMinerba(c *fiber.Ctx) error {
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

	var filterMinerba minerba.FilterMinerba

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		filterMinerba.Quantity = 0
	} else {
		filterMinerba.Quantity = quantity
	}

	filterMinerba.CreatedStart = c.Query("created_start")
	filterMinerba.CreatedEnd = c.Query("created_end")

	listMinerba, listMinerbaErr := h.minerbaService.GetListReportMinerbaAll(pageNumber, filterMinerba)

	if listMinerbaErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(listMinerba)
}

func (h *minerbaHandler) DetailMinerba(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detailMinerba, detailMinerbaErr := h.transactionService.GetDetailMinerba(idInt)

	if detailMinerbaErr != nil {
		status := 400

		if  detailMinerbaErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaErr.Error(),
		})
	}

	return c.Status(200).JSON(detailMinerba)
}

func (h *minerbaHandler) RequestCreateExcelMinerba(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	header := c.GetReqHeaders()

	detailMinerba, detailMinerbaErr := h.transactionService.GetDetailMinerba(idInt)

	if detailMinerbaErr != nil {
		status := 400

		if  detailMinerbaErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailMinerbaErr.Error(),
		})
	}

	var inputRequestCreateExcel transaction.InputRequestCreateExcelMinerba
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.MinerbaId = idInt
	inputRequestCreateExcel.MinerbaNumber = *detailMinerba.Detail.IdNumber
	inputRequestCreateExcel.MinerbaPeriod = detailMinerba.Detail.Period
	inputRequestCreateExcel.Transactions = detailMinerba.List

	hitJob, hitJobErr := h.transactionService.RequestCreateExcel(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *minerbaHandler) CheckValidPeriodMinerba(c *fiber.Ctx) error {
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

	inputCheckMinerba := new(minerba.CheckMinerbaPeriod)

	if err := c.BodyParser(inputCheckMinerba); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCheckMinerba)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid period",
			"errors": dataErrors,
		})
	}

	_, detailMinerbaErr := h.minerbaService.GetReportMinerbaWithPeriod(inputCheckMinerba.Period)

	if  detailMinerbaErr != nil && detailMinerbaErr.Error() == "record not found" {
		return c.Status(200).JSON(fiber.Map{
			"message": "valid period",
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"message": "invalid period",
	})

}
