package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/coareportln"
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

type coaReportLnHandler struct {
	coaReportLnService      coareportln.Service
	logService              logs.Service
	userIupopkService       useriupopk.Service
	historyService          history.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}

func NewCoaReportLnHandler(coaReportLnService coareportln.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
) *coaReportLnHandler {
	return &coaReportLnHandler{
		coaReportLnService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
	}
}

func (h *coaReportLnHandler) ListCoaReportLnTransaction(c *fiber.Ctx) error {
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

	coaReportLnInput := new(coareportln.CoaReportLnInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportLnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportLnInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	listTransaction, listTransactionErr := h.coaReportLnService.GetTransactionCoaReportLn(coaReportLnInput.DateFrom, coaReportLnInput.DateTo, iupopkIdInt)

	if listTransactionErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": listTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(listTransaction)
}

func (h *coaReportLnHandler) CreateCoaReportLn(c *fiber.Ctx) error {
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

	coaReportLnInput := new(coareportln.CoaReportLnInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportLnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportLnInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	coaReportLn, coaReportLnErr := h.historyService.CreateCoaReportLn(coaReportLnInput.DateFrom, coaReportLnInput.DateTo, iupopkIdInt, uint(claims["id"].(float64)))

	if coaReportLnErr != nil {
		inputJson, _ := json.Marshal(coaReportLnInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": coaReportLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": coaReportLnErr.Error(),
		})
	}

	return c.Status(201).JSON(coaReportLn)
}

func (h *coaReportLnHandler) DeleteCoaReportLn(c *fiber.Ctx) error {
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

	detailCoaReportLn, detailCoaReportLnErr := h.coaReportLnService.GetDetailTransactionCoaReportLn(idInt, iupopkIdInt)

	if detailCoaReportLnErr != nil {
		status := 400

		if detailCoaReportLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete coa report ln",
			"error":   detailCoaReportLnErr.Error(),
		})
	}

	_, isDeletedCoaReportLnErr := h.historyService.DeleteCoaReportLn(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedCoaReportLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["coa_report_ln_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedCoaReportLnErr.Error(),
		})

		coaReportLnId := uint(idInt)
		createdErrLog := logs.Logs{
			CoaReportLnId: &coaReportLnId,
			Input:         inputJson,
			Message:       messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedCoaReportLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete coa report ln",
			"error":   isDeletedCoaReportLnErr.Error(),
		})
	}

	if detailCoaReportLn.Detail.CoaReportLnDocumentLink != "" {
		documentLink := detailCoaReportLn.Detail.CoaReportLnDocumentLink

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
			inputMap["coa_report_ln_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": isDeletedCoaReportLnErr.Error(),
			})

			coaReportLnId := uint(idInt)
			createdErrLog := logs.Logs{
				CoaReportLnId: &coaReportLnId,
				Input:         inputJson,
				Message:       messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete coa report ln aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete coa report ln",
	})
}

func (h *coaReportLnHandler) DetailCoaReportLn(c *fiber.Ctx) error {
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

	detailCoaReportLn, detailCoaReportLnErr := h.coaReportLnService.GetDetailTransactionCoaReportLn(idInt, iupopkIdInt)

	if detailCoaReportLnErr != nil {
		status := 400

		if detailCoaReportLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportLnErr.Error(),
		})
	}

	return c.Status(200).JSON(detailCoaReportLn)
}

func (h *coaReportLnHandler) UpdateDocumentCoaReportLn(c *fiber.Ctx) error {
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

	inputUpdateCoaReportLn := new(coareportln.InputUpdateDocumentCoaReportLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateCoaReportLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update coa report ln",
		})
	}

	errors := h.v.Struct(*inputUpdateCoaReportLn)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update coa report ln",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReportLn
		inputMap["coa_report_ln_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		coaReportLnId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:         inputJson,
			Message:       messageJson,
			CoaReportLnId: &coaReportLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailCoaReportLn, detailCoaReportLnErr := h.coaReportLnService.GetDetailTransactionCoaReportLn(idInt, iupopkIdInt)

	if detailCoaReportLnErr != nil {
		status := 400

		if detailCoaReportLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportLnErr.Error(),
		})
	}

	updateCoaReportLn, updateCoaReportLnErr := h.historyService.UpdateDocumentCoaReportLn(idInt, *inputUpdateCoaReportLn, uint(claims["id"].(float64)), iupopkIdInt)

	if updateCoaReportLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReportLn
		inputMap["coa_report_ln_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateCoaReportLnErr.Error(),
		})

		coaReportLnId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:         inputJson,
			Message:       messageJson,
			CoaReportLnId: &coaReportLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateCoaReportLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateCoaReportLnErr.Error(),
			"message": "failed to update coa report ln",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "coa report ln"
	inputNotification.Period = fmt.Sprintf("%v/%v", detailCoaReportLn.Detail.DateFrom, detailCoaReportLn.Detail.DateTo)
	inputNotification.Status = "membuat dokumen"
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateCoaReportLn
		inputMap["coa_report_ln_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		coaReportLnId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:         inputJson,
			Message:       messageJson,
			CoaReportLnId: &coaReportLnId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update coa report ln",
		})
	}

	return c.Status(200).JSON(updateCoaReportLn)
}

func (h *coaReportLnHandler) RequestCreateExcelCoaReportLn(c *fiber.Ctx) error {
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

	detailCoaReportLn, detailCoaReportLnErr := h.coaReportLnService.GetDetailTransactionCoaReportLn(idInt, iupopkIdInt)

	if detailCoaReportLnErr != nil {
		status := 400

		if detailCoaReportLnErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailCoaReportLnErr.Error(),
		})
	}

	var inputRequestCreateExcel coareportln.InputRequestCreateUploadCoaReportLn
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.CoaReport = detailCoaReportLn.Detail
	inputRequestCreateExcel.ListTransaction = detailCoaReportLn.ListTransaction
	inputRequestCreateExcel.Iupopk = detailCoaReportLn.Detail.Iupopk
	hitJob, hitJobErr := h.coaReportLnService.RequestCreateExcelCoaReportLn(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *coaReportLnHandler) ListCoaReportLn(c *fiber.Ctx) error {
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

	var filterCoaReportLn coareportln.SortFilterCoaReportLn

	filterCoaReportLn.Field = c.Query("field")
	filterCoaReportLn.Sort = c.Query("sort")
	filterCoaReportLn.DateStart = c.Query("date_start")
	filterCoaReportLn.DateEnd = c.Query("date_end")

	listCoaReportLn, listCoaReportLnErr := h.coaReportLnService.ListCoaReportLn(pageNumber, filterCoaReportLn, iupopkIdInt)

	if listCoaReportLnErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listCoaReportLnErr.Error(),
		})
	}

	return c.Status(200).JSON(listCoaReportLn)
}
