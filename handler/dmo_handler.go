package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/dmo"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
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
	notificationUserService notificationuser.Service
	v *validator.Validate
}

func NewDmoHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, dmoService dmo.Service, traderService trader.Service, traderDmoService traderdmo.Service, notificationUserService notificationuser.Service, v *validator.Validate) *dmoHandler {
	return &dmoHandler{
		transactionService,
		userService,
		historyService,
		logService,
		dmoService,
		traderService,
		traderDmoService,
		notificationUserService,
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateDmo := new(dmo.CreateDmoInput)

	// Binds the request body to the Person struct
	err := c.BodyParser(inputCreateDmo)

	if err != nil {
		inputCreateDmo.Period = strings.Replace(inputCreateDmo.Period, "\"", "", -1)

		formPart, errFormPart := c.MultipartForm()
		if 	errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}

		for _, trader := range formPart.Value["trader"] {
			var traderArrayInt []int
			errUnmarshal := json.Unmarshal([]byte(trader), &traderArrayInt)
			fmt.Println(errUnmarshal)
			inputCreateDmo.Trader = traderArrayInt
		}

		for _, barge := range formPart.Value["transaction_barge"] {
			var bargeArrayInt []int
			errUnmarshal := json.Unmarshal([]byte(barge), &bargeArrayInt)
			fmt.Println(errUnmarshal)
			inputCreateDmo.TransactionBarge = bargeArrayInt
		}

		for _, vessel := range formPart.Value["transaction_vessel"] {
			var bargeVesselInt []int
			errUnmarshal := json.Unmarshal([]byte(vessel), &bargeVesselInt)
			fmt.Println(errUnmarshal)
			inputCreateDmo.TransactionVessel = bargeVesselInt
		}

		for _, vesselAdjustment := range formPart.Value["vessel_adjustment"] {
			var newVesselAdjustmentArray []dmo.VesselAdjustmentInput
			errUnmarshal := json.Unmarshal([]byte(vesselAdjustment), &newVesselAdjustmentArray)
			fmt.Println(errUnmarshal)

			inputCreateDmo.VesselAdjustment = newVesselAdjustmentArray
		}
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

	if !inputCreateDmo.IsDocumentCustom {
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
	}

	if inputCreateDmo.IsDocumentCustom {
		formPart, _ := c.MultipartForm()

		if len(formPart.File) < 3 {
			return c.Status(400).JSON(fiber.Map{
				"message": "dmo created but job failed", "error": "bast, reconciliation_letter, statement_letter document is required",
			})
		}
		bastFile := formPart.File["bast"][0]
		reconciliationLetterFile := formPart.File["reconciliation_letter"][0]
		statementLetterFile := formPart.File["statement_letter"][0]
		_, reqJobDocumentCustomErr := h.transactionService.RequestCreateCustomDmo(createDmo, bastFile, reconciliationLetterFile, statementLetterFile, header["Authorization"])

		if reqJobDocumentCustomErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "dmo created but job failed", "error": reqJobDocumentCustomErr.Error(),
			})
		}
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

	var filterDmo dmo.FilterDmo

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		filterDmo.Quantity = 0
	} else {
		filterDmo.Quantity = quantity
	}

	filterDmo.CreatedStart = c.Query("created_start")
	filterDmo.CreatedEnd = c.Query("created_end")

	listDmo, listDmoErr := h.dmoService.GetListReportDmoAll(pageNumber, filterDmo)

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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
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

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

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

	var createNotif notification.InputNotification

	createNotif.Type = "dmo"
	createNotif.Status = "success delete dmo"
	createNotif.Period = findDmo.Detail.Period
	createNotif.EndUser = endUserDmo.Company.CompanyName

	_, createNotificationDeleteDmoErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)))

	if createNotificationDeleteDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt
		inputMap["dmo_period"] = findDmo.Detail.Period
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteDmoErr.Error(),
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputUpdateDmo := new(dmo.InputUpdateDocumentDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
			"message": "failed to update dmo",
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

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"

	if detailDmo.Detail.IsDocumentCustom {
		inputNotification.Status = "success upload document custom"
	} else {
		inputNotification.Status = "success create document"
	}

	inputNotification.Period = detailDmo.Detail.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo

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
			"message": "failed to create notification update dmo",
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

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
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

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" {
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

	if typeDocument == "bast" {
		isBast = true
	}

	if typeDocument == "statement_letter" {
		isStatementLetter = true
	}

	if typeDocument == "reconciliation_letter" {
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

func (h *dmoHandler) UpdateTrueIsSignedDmoDocument(c *fiber.Ctx) error {
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

	dataDmo, dataDmoErr := h.dmoService.GetDataDmo(idInt)

	if dataDmoErr != nil {
		status := 400

		if dataDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": dataDmoErr.Error(),
		})
	}

	var fileName string

	fileName = *dataDmo.IdNumber

	file, errFormFile := c.FormFile("document")

	if errFormFile != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": errFormFile.Error(),
		})
	}

	typeDocument := c.Params("type")

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" {
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

	if typeDocument == "bast" {
		isBast = true
		fileName += "/signed_bast.pdf"
	}

	if typeDocument == "statement_letter" {
		isStatementLetter = true
		fileName += "/signed_surat_pernyataan.pdf"
	}

	if typeDocument == "reconciliation_letter" {
		isReconciliationLetter = true
		fileName += "/signed_berita_acara.pdf"
	}

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = typeDocument
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt

		inputJson , _ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId: &dmoId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": uploadErr.Error(),
		})
	}

	updateSignedDocumentDmo, updateSignedDocumentDmoErr := h.historyService.UpdateTrueIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, idInt, uint(claims["id"].(float64)), up.Location)

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

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"
	inputNotification.Status = "success signed document"
	inputNotification.Period = dataDmo.Period
	inputNotification.Document = typeDocument
	inputNotification.EndUser = endUserDmo.Company.CompanyName

	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to create notification update signed dmo",
			"error": createdNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(updateSignedDocumentDmo)
}

func (h *dmoHandler) UpdateFalseIsSignedDmoDocument(c *fiber.Ctx) error {
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

	dataDmo, dataDmoErr := h.dmoService.GetDataDmo(idInt)

	if dataDmoErr != nil {
		status := 400

		if dataDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": dataDmoErr.Error(),
		})
	}

	var fileName string

	fileName = *dataDmo.IdNumber

	typeDocument := c.Params("type")

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" {
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

	if typeDocument == "bast" {
		isBast = true
		fileName += "/signed_bast.pdf"
	}

	if typeDocument == "statement_letter" {
		isStatementLetter = true
		fileName += "/signed_surat_pernyataan.pdf"
	}

	if typeDocument == "reconciliation_letter" {
		isReconciliationLetter = true
		fileName += "/signed_berita_acara.pdf"
	}

	deleteUpload, deleteUploadErr := awshelper.DeleteDocument(fileName)

	if deleteUploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["document_type"] = typeDocument
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt
		inputMap["response"] = deleteUpload
		inputJson , _ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteUploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId: &dmoId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error": deleteUploadErr.Error(),
		})
	}

	updateSignedDocumentDmo, updateSignedDocumentDmoErr := h.historyService.UpdateFalseIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, idInt, uint(claims["id"].(float64)))

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

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"
	inputNotification.Status = "success delete signed document"
	inputNotification.Period = dataDmo.Period
	inputNotification.Document = typeDocument
	inputNotification.EndUser = endUserDmo.Company.CompanyName

	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)))

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
			DmoId: &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to create notification delete signed dmo",
			"error": createdNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(updateSignedDocumentDmo)
}
