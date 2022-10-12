package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/dmo"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
	"ajebackend/model/traderdmo"
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

type dmoHandler struct {
	transactionService transaction.Service
	userService user.Service
	historyService history.Service
	logService logs.Service
	dmoService dmo.Service
	traderService trader.Service
	traderDmoService traderdmo.Service
	v *validator.Validate
}

func NewDmoHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, dmoService dmo.Service, traderService trader.Service, traderDmoService traderdmo.Service, v *validator.Validate) *dmoHandler {
	return &dmoHandler{
		transactionService,
		userService,
		historyService,
		logService,
		dmoService,
		traderService,
		traderDmoService,
		v,
	}
}

func (h *dmoHandler) CreateDmo(c *fiber.Ctx) error {
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

	inputCreateDmo := new(dmo.CreateDmoInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
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

	header := c.GetReqHeaders()

	for _, valueVessel := range inputCreateDmo.TransactionVessel {
		for _, valueBarge := range inputCreateDmo.TransactionBarge {
			if valueVessel == valueBarge {
				return c.Status(400).JSON(fiber.Map{
					"error": "please check transaction is in vessel & barge",
				})
			}
		}
	}

	if len(inputCreateDmo.TransactionVessel) > 0 && len(inputCreateDmo.VesselAdjustment) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "please check vessel adjustment",
		})
	}

	if len(inputCreateDmo.TransactionVessel) == 0 && len(inputCreateDmo.TransactionBarge) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "please check there is no transaction vessel and transaction barge",
		})
	}


	_, checkListTraderErr := h.traderService.CheckListTrader(inputCreateDmo.Trader)

	if checkListTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":  "trader " + checkListTraderErr.Error(),
		})
	}

	_, checkEndUserErr := h.traderService.CheckEndUser(inputCreateDmo.EndUser)

	if checkEndUserErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "trader end user " + checkEndUserErr.Error(),
		})
	}

	_, findDmoErr := h.dmoService.GetReportDmoWithPeriod(inputCreateDmo.Period)

	if findDmoErr == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "report with same period already exist",
		})
	}

	var dataTransactions []transaction.Transaction

	if len(inputCreateDmo.TransactionBarge) > 0 {
		checkDmoBarge, checkDmoBargeErr := h.transactionService.CheckDataDnAndDmo(inputCreateDmo.TransactionBarge)

		if checkDmoBargeErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
			inputMap["input"] = inputCreateDmo

			inputJson ,_ := json.Marshal(inputMap)
			messageJson ,_ := json.Marshal(map[string]interface{}{
				"error": checkDmoBargeErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input: inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if checkDmoBargeErr.Error() == "please check there is transaction not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"error": "transaction barge " + checkDmoBargeErr.Error(),
			})
		}

		for _, v := range checkDmoBarge {
			dataTransactions = append(dataTransactions, v)
		}
	}

	if len(inputCreateDmo.TransactionVessel) > 0 {
		checkDmoVessel, checkDmoVesselErr := h.transactionService.CheckDataDnAndDmo(inputCreateDmo.TransactionVessel)

		if checkDmoVesselErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
			inputMap["input"] = inputCreateDmo

			inputJson ,_ := json.Marshal(inputMap)
			messageJson ,_ := json.Marshal(map[string]interface{}{
				"error": checkDmoVesselErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input: inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if checkDmoVesselErr.Error() == "please check there is transaction not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"error": "transaction vessel " + checkDmoVesselErr.Error(),
			})
		}

		for _, v := range checkDmoVessel {
			dataTransactions = append(dataTransactions, v)
		}
	}

	splitPeriod := strings.Split(inputCreateDmo.Period, " ")

	baseIdNumber := fmt.Sprintf("DD-%s-%s",  helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1])
	createDmo, createDmoErr := h.historyService.CreateDmo(*inputCreateDmo, baseIdNumber, uint(claims["id"].(float64)))

	if createDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
		inputMap["input"] = inputCreateDmo
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createDmoErr.Error(),
		})
	}

	list, endUser, listTraderDmoErr := h.traderDmoService.TraderListWithDmoId(int(createDmo.ID))

	if listTraderDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
		inputMap["input"] = inputCreateDmo
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": listTraderDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but get list trader failed",
			"error": listTraderDmoErr.Error(),
		})
	}

	var reqInputCreateUploadDmo transaction.InputRequestCreateUploadDmo

	reqInputCreateUploadDmo.Authorization = header["Authorization"]
	reqInputCreateUploadDmo.BastNumber = fmt.Sprintf("%s/BAST/%s", splitPeriod[1], helper.CreateIdNumber(int(createDmo.ID)))
	reqInputCreateUploadDmo.DataDmo = createDmo
	reqInputCreateUploadDmo.DataTransactions = dataTransactions
	reqInputCreateUploadDmo.Trader = list
	reqInputCreateUploadDmo.TraderEndUser = endUser

	_, requestJobDmoErr := h.transactionService.RequestCreateDmo(reqInputCreateUploadDmo)

	if requestJobDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_dn_vessel"] = inputCreateDmo.TransactionVessel
		inputMap["input"] = inputCreateDmo
		inputMap["input_job"] = reqInputCreateUploadDmo
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": requestJobDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but job failed",
			"error": requestJobDmoErr.Error(),
			"dmo": createDmo,
		})
	}

	return c.Status(201).JSON(createDmo)
}

func (h *dmoHandler) ListDmo(c *fiber.Ctx) error {
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

	listDmo, listDmoErr := h.dmoService.GetListReportDmoAll(pageNumber)

	if listDmoErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(listDmo)
}

func (h *dmoHandler) DetailDmo(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	detailDmo, detailDmoErr := h.transactionService.GetDetailDmo(idInt)

	if detailDmoErr != nil {
		status := 400

		if  detailDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(detailDmo)
}

func (h *dmoHandler) ListDataDNWithoutDmo(c *fiber.Ctx) error {
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

	listDataDNWithoutDmo, listDataDNWithoutDmoErr := h.transactionService.ListDataDNWithoutDmo()

	if listDataDNWithoutDmoErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDataDNWithoutDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": listDataDNWithoutDmo,
	})
}

func (h *dmoHandler) DeleteDmo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
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
			"message": "failed to delete dmo",
			"error": "record not found",
		})
	}

	findDmo, findDmoErr := h.transactionService.GetDetailDmo(idInt)

	if findDmoErr != nil {
		status := 400

		if findDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error": findDmoErr.Error(),
		})
	}

	if findDmo.Detail.IsBastDocumentSigned != false || findDmo.Detail.IsReconciliationLetterSigned != false || findDmo.Detail.IsStatementLetterSigned != false {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error": "document dmo is already signed",
		})
	}

	fileName := fmt.Sprintf("%s/", *findDmo.Detail.IdNumber)
	_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

	if deleteAwsErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteAwsErr.Error(),
			"id_number": findDmo.Detail.IdNumber,
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to delete dmo aws",
			"error": deleteAwsErr.Error(),
		})
	}


	_, deleteDmoErr := h.historyService.DeleteDmo(idInt, uint(claims["id"].(float64)))

	if deleteDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId: &findDmo.Detail.ID,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if  deleteDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error": deleteDmoErr.Error(),
		})
	}



	return c.Status(200).JSON(fiber.Map{
		"message": "success delete dmo",
	})
}

func (h *dmoHandler) UpdateDocumentDmo(c *fiber.Ctx) error {
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

	inputUpdateDmo := new(dmo.InputUpdateDocumentDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
			"message": "failed to update minerba",
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update dmo",
			"error": "record not found",
		})
	}

	errors := h.v.Struct(*inputUpdateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})
		
		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailDmo, detailDmoErr := h.transactionService.GetDetailDmo(idInt)

	if detailDmoErr != nil {
		status := 400

		if  detailDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailDmoErr.Error(),
		})
	}

	if detailDmo.Detail.BASTDocumentLink != nil || detailDmo.Detail.ReconciliationLetterDocumentLink != nil || detailDmo.Detail.StatementLetterDocumentLink != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "document already has been created",
		})
	}

	updateDocumentDmo, updateDocumentDmoErr := h.historyService.UpdateDocumentDmo(idInt, *inputUpdateDmo, uint(claims["id"].(float64)))

	if updateDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateDocumentDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": updateDocumentDmoErr.Error(),
			"message": "failed to update document dmo",
		})
	}

	return c.Status(200).JSON(updateDocumentDmo)
}

func (h *dmoHandler) UpdateIsDownloadedDocumentDmo(c *fiber.Ctx) error {
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

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error": "record not found",
		})
	}

	var isBast bool
	var isStatementLetter bool
	var isReconciliationLetter bool

	typeDocument := c.Params("type")

	if typeDocument != "bast_document_link" && typeDocument != "statement_letter_document_link" && typeDocument != "reconciliation_letter_document_link" {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"message": "failed to update downloaded dmo",
			"error": "type not found",
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error": "type not found",
		})
	}

	if typeDocument == "bast_document_link" {
		isBast = true
	}

	if typeDocument == "statement_letter_document_link" {
		isStatementLetter = true
	}

	if typeDocument == "reconciliation_letter_document_link" {
		isReconciliationLetter = true
	}

	updateDownloadedDocumentDmo, updateDownloadedDocumentDmoErr := h.historyService.UpdateIsDownloadedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, idInt, uint(claims["id"].(float64)))

	if updateDownloadedDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateDownloadedDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateDownloadedDocumentDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error": updateDownloadedDocumentDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(updateDownloadedDocumentDmo)
}

func (h *dmoHandler) UpdateIsSignedDocumentDmo(c *fiber.Ctx) error {
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

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": "record not found",
		})
	}

	var isBast bool
	var isStatementLetter bool
	var isReconciliationLetter bool

	typeDocument := c.Params("type")

	if typeDocument != "bast_document_link" && typeDocument != "statement_letter_document_link" && typeDocument != "reconciliation_letter_document_link" {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"message": "failed to update signed dmo",
			"error": "type not found",
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": "type not found",
		})
	}

	if typeDocument == "bast_document_link" {
		isBast = true
	}

	if typeDocument == "statement_letter_document_link" {
		isStatementLetter = true
	}

	if typeDocument == "reconciliation_letter_document_link" {
		isReconciliationLetter = true
	}

	updateSignedDocumentDmo, updateSignedDocumentDmoErr := h.historyService.UpdateIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, idInt, uint(claims["id"].(float64)))

	if updateSignedDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateSignedDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateSignedDocumentDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": updateSignedDocumentDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(updateSignedDocumentDmo)
}
