package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
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
	reportDmoService        reportdmo.Service
}

func NewReportDmoHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, notificationUserService notificationuser.Service, v *validator.Validate, reportDmoService reportdmo.Service) *reportDmoHandler {
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

	baseIdNumber := fmt.Sprintf("LDO-AJE-%s-%s", helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1][len(splitPeriod[1])-2:])

	createReportDmo, createReportDmoErr := h.historyService.CreateReportDmo(*inputCreateReportDmo, baseIdNumber, uint(claims["id"].(float64)))

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
			"message": "failed to update report dmo",
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
			"message": "failed to update report dmo",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "report dmo"
	inputNotification.Status = "membuat"
	inputNotification.Period = updateReportDmo.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), 11)

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
	inputRequestCreateReportDmo.GroupingVessels = detailReportDmo.GroupingVessels
	inputRequestCreateReportDmo.Transactions = detailReportDmo.Transactions

	hitJob, hitJobErr := h.transactionService.RequestCreateReportDmo(inputRequestCreateReportDmo)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *reportDmoHandler) UpdateReportDmo(c *fiber.Ctx) error {
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

	inputUpdateReportDmo := new(reportdmo.InputUpdateReportDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateReportDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	idInt, err := strconv.Atoi(id)

	reportDmoId := uint(idInt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update report dmo",
			"error":   "record not found",
		})
	}

	errors := h.v.Struct(*inputUpdateReportDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_id"] = idInt
		inputMap["list"] = inputUpdateReportDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			ReportDmoId: &reportDmoId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	_, findDetailReportDmoErr := h.transactionService.GetDetailReportDmo(idInt)

	if findDetailReportDmoErr != nil {
		status := 400

		if findDetailReportDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": findDetailReportDmoErr.Error(),
		})
	}

	updateReportDmo, updateReportDmoErr := h.historyService.UpdateTransactionReportDmo(idInt, *inputUpdateReportDmo, uint(claims["id"].(float64)))

	if updateReportDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_id"] = idInt
		inputMap["list"] = inputUpdateReportDmo

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateReportDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			ReportDmoId: &reportDmoId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateReportDmoErr.Error(),
		})
	}

	var createNotif notification.InputNotification

	createNotif.Type = "report dmo"
	createNotif.Status = "mengedit"
	createNotif.Period = updateReportDmo.Period

	_, createNotifUpdateReportDmo := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), 11)

	if createNotifUpdateReportDmo != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_id"] = idInt
		inputMap["list"] = inputUpdateReportDmo

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotifUpdateReportDmo.Error(),
		})

		createdErrLog := logs.Logs{
			ReportDmoId: &reportDmoId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotifUpdateReportDmo.Error(),
		})
	}

	return c.Status(200).JSON(updateReportDmo)
}

func (h *reportDmoHandler) DeleteReportDmo(c *fiber.Ctx) error {
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
			"message": "failed to delete report dmo",
			"error":   "record not found",
		})
	}

	detailReportDmo, getDetailReportDmoErr := h.transactionService.GetDetailReportDmo(idInt)

	if getDetailReportDmoErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete report dmo",
			"error":   "record not found",
		})
	}

	_, deleteReportDmoErr := h.historyService.DeleteReportDmo(idInt, uint(claims["id"].(float64)))

	if deleteReportDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteReportDmoErr.Error(),
		})

		reportDmoId := uint(idInt)
		createdErrLog := logs.Logs{
			ReportDmoId: &reportDmoId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteReportDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete report",
			"error":   deleteReportDmoErr.Error(),
		})
	}

	if detailReportDmo.Detail.RecapDmoDocumentLink != nil || detailReportDmo.Detail.DetailDmoDocumentLink != nil {

		fileName := fmt.Sprintf("AJE/LDO/%s/", *detailReportDmo.Detail.IdNumber)

		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["report_dmo_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			reportDmoId := uint(idInt)
			createdErrLog := logs.Logs{
				ReportDmoId: &reportDmoId,
				Input:       inputJson,
				Message:     messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete report dmo aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	var createNotif notification.InputNotification

	createNotif.Type = "report dmo"
	createNotif.Status = "menghapus"
	createNotif.Period = detailReportDmo.Detail.Period

	_, createNotificationDeleteReportDmoErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), 11)

	if createNotificationDeleteReportDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["report_dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteReportDmoErr.Error(),
		})

		reportDmoId := uint(idInt)
		createdErrLog := logs.Logs{
			ReportDmoId: &reportDmoId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteReportDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete report dmo",
	})
}

func (h *reportDmoHandler) GetListForReport(c *fiber.Ctx) error {
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

	getListForReport, getListForReportErr := h.transactionService.GetListForReport()

	if getListForReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": getListForReportErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": getListForReport,
	})
}

func (h *reportDmoHandler) CheckValidPeriodReportDmo(c *fiber.Ctx) error {
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

	inputCheckReportDmo := new(reportdmo.CheckReportDmoPeriod)

	if err := c.BodyParser(inputCheckReportDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCheckReportDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid period",
			"errors":  dataErrors,
		})
	}

	_, detailReportDmoErr := h.reportDmoService.GetReportDmoWithPeriod(inputCheckReportDmo.Period)

	if detailReportDmoErr != nil && detailReportDmoErr.Error() == "record not found" {
		return c.Status(200).JSON(fiber.Map{
			"message": "valid period",
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"message": "invalid period",
	})
}

func (h *reportDmoHandler) DetailReportDmo(c *fiber.Ctx) error {
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

	return c.Status(200).JSON(detailReportDmo)
}

func (h *reportDmoHandler) ListReportDmo(c *fiber.Ctx) error {
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

	var filterReportDmo reportdmo.FilterAndSortReportDmo

	filterReportDmo.Quantity = c.Query("quantity")
	filterReportDmo.UpdatedStart = c.Query("updated_start")
	filterReportDmo.UpdatedEnd = c.Query("updated_end")
	filterReportDmo.Month = c.Query("month")
	filterReportDmo.Year = c.Query("year")
	filterReportDmo.Field = c.Query("field")
	filterReportDmo.Sort = c.Query("sort")

	listReportDmo, listReportDmoErr := h.reportDmoService.GetListReportDmoAll(pageNumber, filterReportDmo)

	if listReportDmoErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listReportDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(listReportDmo)
}
