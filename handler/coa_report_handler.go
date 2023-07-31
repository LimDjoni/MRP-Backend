package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/coareport"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/useriupopk"
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

type coaReportHandler struct {
	coaReportService        coareport.Service
	logService              logs.Service
	userIupopkService       useriupopk.Service
	historyService          history.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}

func NewCoaReportHandler(coaReportService coareport.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
) *coaReportHandler {
	return &coaReportHandler{
		coaReportService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
	}
}

// Preview
func (h *coaReportHandler) ListCoaReportTransaction(c *fiber.Ctx) error {
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

	coaReportInput := new(coareport.CoaReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	listTransaction, listTransactionErr := h.coaReportService.GetTransactionCoaReport(coaReportInput.DateFrom, coaReportInput.DateTo, iupopkIdInt)

	if listTransactionErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": listTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(listTransaction)
}

func (h *coaReportHandler) CreateCoaReport(c *fiber.Ctx) error {
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

	coaReportInput := new(coareport.CoaReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	coaReport, coaReportErr := h.historyService.CreateCoaReport(coaReportInput.DateFrom, coaReportInput.DateTo, iupopkIdInt, uint(claims["id"].(float64)))

	if coaReportErr != nil {
		inputJson, _ := json.Marshal(coaReportInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": coaReportErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": coaReportErr.Error(),
		})
	}

	return c.Status(201).JSON(coaReport)
}

func (h *coaReportHandler) DeleteCoaReport(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	iupopkId := c.Params("iupopk_id")
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailCoaReport, detailCoaReportErr := h.coaReportService.GetDetailTransactionCoaReport(idInt, iupopkIdInt)

	if detailCoaReportErr != nil {
		status := 400

		if detailCoaReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete coa report",
			"error":   detailCoaReportErr.Error(),
		})
	}

	_, isDeletedCoaReportErr := h.historyService.DeleteCoaReport(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedCoaReportErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["coa_report_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedCoaReportErr.Error(),
		})

		coaReportId := uint(idInt)
		createdErrLog := logs.Logs{
			CoaReportId: &coaReportId,
			Input:       inputJson,
			Message:     messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedCoaReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete coa report",
			"error":   isDeletedCoaReportErr.Error(),
		})
	}

	if detailCoaReport.Detail.CoaReportDocumentLink != "" {
		documentLink := detailCoaReport.Detail.CoaReportDocumentLink

		documentLinkSplit := strings.Split(documentLink, "/")

		fileName := ""

		for i, v := range documentLinkSplit {
			if i == 3 {
				fileName += v + "/"
			}

			if i == 4 {
				fileName += v + "/"
			}

			if i == 5 {
				fileName += v
			}
		}
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["coa_report_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": isDeletedCoaReportErr.Error(),
			})

			coaReportId := uint(idInt)
			createdErrLog := logs.Logs{
				CoaReportId: &coaReportId,
				Input:       inputJson,
				Message:     messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete coa report aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete coa report",
	})
}

func (h *coaReportHandler) DetailCoaReport(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	iupopkId := c.Params("iupopk_id")
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailCoaReport, detailCoaReportErr := h.coaReportService.GetDetailTransactionCoaReport(idInt, iupopkIdInt)

	if detailCoaReportErr != nil {
		status := 400

		if detailCoaReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportErr.Error(),
		})
	}

	return c.Status(200).JSON(detailCoaReport)
}

func (h *coaReportHandler) UpdateDocumentCoaReport(c *fiber.Ctx) error {
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

	inputUpdateCoaReport := new(coareport.InputUpdateDocumentCoaReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateCoaReport); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update coa report",
		})
	}

	errors := h.v.Struct(*inputUpdateCoaReport)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update coa report",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReport
		inputMap["coa_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		coaReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			CoaReportId: &coaReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailCoaReport, detailCoaReportErr := h.coaReportService.GetDetailTransactionCoaReport(idInt, iupopkIdInt)

	if detailCoaReportErr != nil {
		status := 400

		if detailCoaReportErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportErr.Error(),
		})
	}

	updateCoaReport, updateCoaReportErr := h.historyService.UpdateDocumentCoaReport(idInt, *inputUpdateCoaReport, uint(claims["id"].(float64)), iupopkIdInt)

	if updateCoaReportErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReport
		inputMap["coa_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateCoaReportErr.Error(),
		})

		coaReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			CoaReportId: &coaReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateCoaReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateCoaReportErr.Error(),
			"message": "failed to update coa report",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "coa report"
	inputNotification.Period = fmt.Sprintf("%v/%v", detailCoaReport.Detail.DateFrom, detailCoaReport.Detail.DateTo)
	inputNotification.Status = "membuat dokumen"
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReport
		inputMap["coa_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		coaReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:       inputJson,
			Message:     messageJson,
			CoaReportId: &coaReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update coa report",
		})
	}

	return c.Status(200).JSON(updateCoaReport)
}

func (h *coaReportHandler) RequestCreateExcelCoaReport(c *fiber.Ctx) error {
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

	header := c.GetReqHeaders()

	detailCoaReport, detailCoaReportErr := h.coaReportService.GetDetailTransactionCoaReport(idInt, iupopkIdInt)

	if detailCoaReportErr != nil {
		status := 400

		if detailCoaReportErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportErr.Error(),
		})
	}

	var inputRequestCreateExcel coareport.InputRequestCreateUploadCoaReport
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.CoaReport = detailCoaReport.Detail
	inputRequestCreateExcel.ListTransaction = detailCoaReport.ListTransaction
	inputRequestCreateExcel.Iupopk = detailCoaReport.Detail.Iupopk
	hitJob, hitJobErr := h.coaReportService.RequestCreateExcelCoaReport(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *coaReportHandler) ListCoaReport(c *fiber.Ctx) error {
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

	var filterCoaReport coareport.SortFilterCoaReport

	filterCoaReport.Field = c.Query("field")
	filterCoaReport.Sort = c.Query("sort")
	filterCoaReport.DateStart = c.Query("date_start")
	filterCoaReport.DateEnd = c.Query("date_end")

	listCoaReport, listCoaReportErr := h.coaReportService.ListCoaReport(pageNumber, filterCoaReport, iupopkIdInt)

	if listCoaReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listCoaReportErr.Error(),
		})
	}

	return c.Status(200).JSON(listCoaReport)
}
