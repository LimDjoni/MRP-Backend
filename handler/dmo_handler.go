package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/counter"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
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

type dmoHandler struct {
	transactionService      transaction.Service
	historyService          history.Service
	logService              logs.Service
	dmoService              dmo.Service
	traderService           trader.Service
	traderDmoService        traderdmo.Service
	notificationUserService notificationuser.Service
	companyService          company.Service
	v                       *validator.Validate
	dmoVesselService        dmovessel.Service
	userIupopkService       useriupopk.Service
	counterService          counter.Service
}

func NewDmoHandler(
	transactionService transaction.Service,
	historyService history.Service,
	logService logs.Service,
	dmoService dmo.Service,
	traderService trader.Service,
	traderDmoService traderdmo.Service,
	notificationUserService notificationuser.Service,
	companyService company.Service,
	v *validator.Validate,
	dmoVesselService dmovessel.Service,
	userIupopkService useriupopk.Service,
	counterService counter.Service) *dmoHandler {
	return &dmoHandler{
		transactionService,
		historyService,
		logService,
		dmoService,
		traderService,
		traderDmoService,
		notificationUserService,
		companyService,
		v,
		dmoVesselService,
		userIupopkService,
		counterService,
	}
}

func (h *dmoHandler) CreateDmo(c *fiber.Ctx) error {
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

	inputCreateDmo := new(dmo.CreateDmoInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputCreateDmo)

	if errParsing != nil {
		inputCreateDmo.Period = strings.Replace(inputCreateDmo.Period, "\"", "", -1)

		formPart, errFormPart := c.MultipartForm()
		if errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}

		if len(inputCreateDmo.Trader) == 0 {
			for _, trader := range formPart.Value["trader"] {
				var traderArrayInt []int
				errUnmarshal := json.Unmarshal([]byte(trader), &traderArrayInt)
				fmt.Println(errUnmarshal)
				inputCreateDmo.Trader = traderArrayInt
			}
		}

		if len(inputCreateDmo.TransactionBarge) == 0 {
			for _, barge := range formPart.Value["transaction_barge"] {
				var bargeArrayInt []int
				errUnmarshal := json.Unmarshal([]byte(barge), &bargeArrayInt)
				fmt.Println(errUnmarshal)
				inputCreateDmo.TransactionBarge = bargeArrayInt
			}
		}

		if len(inputCreateDmo.GroupingVessel) == 0 {
			for _, vessel := range formPart.Value["grouping_vessel"] {
				var groupingVesselInt []int
				errUnmarshal := json.Unmarshal([]byte(vessel), &groupingVesselInt)
				fmt.Println(errUnmarshal)
				inputCreateDmo.GroupingVessel = groupingVesselInt
			}
		}
	}

	errors := h.v.Struct(*inputCreateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["transaction_barge"] = inputCreateDmo.TransactionBarge
		inputMap["grouping_vessel"] = inputCreateDmo.GroupingVessel
		inputMap["input"] = inputCreateDmo
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
	header := c.GetReqHeaders()

	if len(inputCreateDmo.TransactionBarge) == 1 {
		if inputCreateDmo.TransactionBarge[0] == 0 {
			var tempTransactionBarge []int
			inputCreateDmo.TransactionBarge = tempTransactionBarge
		}
	}

	if len(inputCreateDmo.GroupingVessel) == 1 {
		if inputCreateDmo.GroupingVessel[0] == 0 {
			var tempGroupingVessel []int
			inputCreateDmo.GroupingVessel = tempGroupingVessel
		}
	}

	if len(inputCreateDmo.Trader) == 1 {
		if inputCreateDmo.Trader[0] == 0 {
			var tempTrader []int
			inputCreateDmo.Trader = tempTrader
		}
	}

	if len(inputCreateDmo.GroupingVessel) == 0 && len(inputCreateDmo.TransactionBarge) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "please check there is no grouping vessel and transaction barge",
		})
	}

	if len(inputCreateDmo.Trader) > 0 {
		_, checkListTraderErr := h.traderService.CheckListTrader(inputCreateDmo.Trader)

		if checkListTraderErr != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "trader " + checkListTraderErr.Error(),
			})
		}
	}

	_, checkEndUserErr := h.traderService.CheckEndUser(inputCreateDmo.EndUser)

	if checkEndUserErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "trader end user " + checkEndUserErr.Error(),
		})
	}

	var dataTransactions []transaction.Transaction
	if len(inputCreateDmo.TransactionBarge) > 0 {
		checkDmoBarge, checkDmoBargeErr := h.transactionService.CheckDataDnAndDmo(inputCreateDmo.TransactionBarge, iupopkIdInt)

		if checkDmoBargeErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_group_vessel"] = inputCreateDmo.GroupingVessel
			inputMap["input"] = inputCreateDmo

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": checkDmoBargeErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
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

	if len(inputCreateDmo.GroupingVessel) > 0 {
		_, checkDmoVesselErr := h.transactionService.CheckGroupingVesselAndDmo(inputCreateDmo.GroupingVessel, iupopkIdInt)

		if checkDmoVesselErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_group_vessel"] = inputCreateDmo.GroupingVessel
			inputMap["input"] = inputCreateDmo
			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": checkDmoVesselErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if checkDmoVesselErr.Error() == "please check there is transaction not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"error": "grouping vessel " + checkDmoVesselErr.Error(),
			})
		}

	}

	splitDocumentDate := strings.Split(inputCreateDmo.DocumentDate, "-")

	createDmo, createDmoErr := h.historyService.CreateDmo(*inputCreateDmo, uint(claims["id"].(float64)), iupopkIdInt)

	if createDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_group_vessel"] = inputCreateDmo.GroupingVessel
		inputMap["input"] = inputCreateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
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
		inputMap["list_group_vessel"] = inputCreateDmo.GroupingVessel
		inputMap["input"] = inputCreateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": listTraderDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but get list trader failed",
			"error":   listTraderDmoErr.Error(),
		})
	}

	listTransactionDmo, listTransactionDmoErr := h.transactionService.GetDataDmo(createDmo.ID, iupopkIdInt)

	if listTransactionDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_dn_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_group_vessel"] = inputCreateDmo.GroupingVessel
		inputMap["input"] = inputCreateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": listTransactionDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but get list transaction failed",
			"error":   listTransactionDmoErr.Error(),
		})
	}

	var reqInputCreateUploadDmo transaction.InputRequestCreateUploadDmo

	reqInputCreateUploadDmo.Authorization = header["Authorization"]
	idNumberSplit := strings.Split(*createDmo.IdNumber, "-")

	counterDetail, counterDetailErr := h.counterService.GetCounter(iupopkIdInt)

	if counterDetailErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = inputCreateDmo.Period
		inputMap["list_transaction_barge"] = inputCreateDmo.TransactionBarge
		inputMap["list_grouping_vessel"] = inputCreateDmo.GroupingVessel
		inputMap["input"] = inputCreateDmo
		inputMap["input_job"] = reqInputCreateUploadDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": counterDetailErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "dmo created but job failed", "error": counterDetailErr.Error(),
		})
	}

	var bastFormat string
	bastFormatSplit := strings.Split(counterDetail.BastFormat, "/")

	for i, v := range bastFormatSplit {
		switch v {
		case "BAST":
			bastFormat += "BAST"
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "CODE":
			bastFormat += idNumberSplit[1]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "COUNTER":
			if len(idNumberSplit) == 1 {
				bastFormat += "0" + idNumberSplit[4]
			} else {
				bastFormat += idNumberSplit[4]
			}
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "MM":
			bastFormat += splitDocumentDate[1]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "YYYY":
			bastFormat += splitDocumentDate[0]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		}
	}

	reqInputCreateUploadDmo.BastNumber = bastFormat
	reqInputCreateUploadDmo.DataDmo = createDmo
	reqInputCreateUploadDmo.Trader = list
	reqInputCreateUploadDmo.TraderEndUser = endUser
	reqInputCreateUploadDmo.ListTransactionBarge = listTransactionDmo.ListTransactionBarge
	reqInputCreateUploadDmo.ListTransactionGroupingVessel = listTransactionDmo.ListTransactionGroupingVessel
	reqInputCreateUploadDmo.ListGroupingVessel = listTransactionDmo.ListGroupingVessel
	reqInputCreateUploadDmo.Iupopk = createDmo.Iupopk

	if !inputCreateDmo.IsDocumentCustom {
		_, requestJobDmoErr := h.transactionService.RequestCreateDmo(reqInputCreateUploadDmo)

		if requestJobDmoErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_transaction_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_grouping_vessel"] = inputCreateDmo.GroupingVessel
			inputMap["input"] = inputCreateDmo
			inputMap["input_job"] = reqInputCreateUploadDmo
			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": requestJobDmoErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "Dmo already has been created, but job failed",
				"error":   requestJobDmoErr.Error(),
				"dmo":     createDmo,
			})
		}
	}

	if inputCreateDmo.IsDocumentCustom {
		formPart, _ := c.MultipartForm()

		if len(formPart.File) < 1 {
			return c.Status(400).JSON(fiber.Map{
				"message": "dmo created but job failed", "error": "reconciliation_letter document is required",
			})
		}
		reconciliationLetterFile := formPart.File["reconciliation_letter"][0]

		_, reqJobDocumentCustomErr := h.transactionService.RequestCreateCustomDmo(createDmo, endUser, reconciliationLetterFile, header["Authorization"], reqInputCreateUploadDmo)

		if reqJobDocumentCustomErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = inputCreateDmo.Period
			inputMap["list_transaction_barge"] = inputCreateDmo.TransactionBarge
			inputMap["list_grouping_vessel"] = inputCreateDmo.GroupingVessel
			inputMap["input"] = inputCreateDmo
			inputMap["input_job"] = reqInputCreateUploadDmo
			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": reqJobDocumentCustomErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "dmo created but job failed", "error": reqJobDocumentCustomErr.Error(),
			})
		}
	}

	return c.Status(201).JSON(createDmo)
}

func (h *dmoHandler) UpdateDmo(c *fiber.Ctx) error {
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

	idDmo := uint(idInt)

	inputUpdateDmo := new(dmo.UpdateDmoInput)

	if err := c.BodyParser(inputUpdateDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	detailDmo, detailDmoErr := h.transactionService.GetDetailDmo(idInt, iupopkIdInt)

	if detailDmoErr != nil {
		status := 400

		if detailDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailDmoErr.Error(),
		})
	}

	errors := h.v.Struct(*inputUpdateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt
		inputMap["transaction_barge"] = inputUpdateDmo.TransactionBarge
		inputMap["grouping_vessel"] = inputUpdateDmo.GroupingVessel
		inputMap["input"] = inputUpdateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &idDmo,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}
	header := c.GetReqHeaders()

	if len(inputUpdateDmo.GroupingVessel) == 0 && len(inputUpdateDmo.TransactionBarge) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "please check there is no grouping vessel and transaction barge",
		})
	}

	updateDmo, updateDmoErr := h.historyService.UpdateDmo(*inputUpdateDmo, idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if updateDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = id
		inputMap["list_dn_barge"] = inputUpdateDmo.TransactionBarge
		inputMap["list_group_vessel"] = inputUpdateDmo.GroupingVessel
		inputMap["input"] = inputUpdateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &idDmo,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateDmoErr.Error(),
		})
	}

	list, endUser, listTraderDmoErr := h.traderDmoService.TraderListWithDmoId(int(updateDmo.ID))

	if listTraderDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = id
		inputMap["list_dn_barge"] = inputUpdateDmo.TransactionBarge
		inputMap["list_group_vessel"] = inputUpdateDmo.GroupingVessel
		inputMap["input"] = inputUpdateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": listTraderDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but get list trader failed",
			"error":   listTraderDmoErr.Error(),
		})
	}

	listTransactionDmo, listTransactionDmoErr := h.transactionService.GetDataDmo(updateDmo.ID, iupopkIdInt)

	if listTransactionDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = id
		inputMap["list_dn_barge"] = inputUpdateDmo.TransactionBarge
		inputMap["list_group_vessel"] = inputUpdateDmo.GroupingVessel
		inputMap["input"] = inputUpdateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": listTransactionDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been created, but get list transaction failed",
			"error":   listTransactionDmoErr.Error(),
		})
	}

	splitDocumentDate := strings.Split(detailDmo.Detail.DocumentDate, "-")
	idNumberSplit := strings.Split(*detailDmo.Detail.IdNumber, "-")

	counterDetail, counterDetailErr := h.counterService.GetCounter(iupopkIdInt)

	if counterDetailErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_period"] = detailDmo.Detail.Period
		inputMap["list_transaction_barge"] = inputUpdateDmo.TransactionBarge
		inputMap["list_grouping_vessel"] = inputUpdateDmo.GroupingVessel
		inputMap["input"] = inputUpdateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": counterDetailErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &idDmo,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "Dmo already has been updated, but job failed",
			"error":   counterDetailErr.Error(),
			"dmo":     updateDmo,
		})
	}

	var bastFormat string
	bastFormatSplit := strings.Split(counterDetail.BastFormat, "/")

	for i, v := range bastFormatSplit {
		switch v {
		case "BAST":
			bastFormat += "BAST"
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "CODE":
			bastFormat += idNumberSplit[1]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "COUNTER":
			if len(idNumberSplit) == 1 {
				bastFormat += "0" + idNumberSplit[4]
			} else {
				bastFormat += idNumberSplit[4]
			}
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "MM":
			bastFormat += splitDocumentDate[1]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		case "YYYY":
			bastFormat += splitDocumentDate[0]
			if i < len(bastFormatSplit)-1 {
				bastFormat += "/"
			}
		}
	}

	var reqInputCreateUploadDmo transaction.InputRequestCreateUploadDmo

	reqInputCreateUploadDmo.Authorization = header["Authorization"]
	reqInputCreateUploadDmo.BastNumber = bastFormat
	reqInputCreateUploadDmo.DataDmo = updateDmo
	reqInputCreateUploadDmo.Trader = list
	reqInputCreateUploadDmo.TraderEndUser = endUser
	reqInputCreateUploadDmo.ListTransactionBarge = listTransactionDmo.ListTransactionBarge
	reqInputCreateUploadDmo.ListGroupingVessel = listTransactionDmo.ListGroupingVessel
	reqInputCreateUploadDmo.ListTransactionGroupingVessel = listTransactionDmo.ListTransactionGroupingVessel
	reqInputCreateUploadDmo.Iupopk = updateDmo.Iupopk

	if !detailDmo.Detail.IsDocumentCustom {
		_, requestJobDmoErr := h.transactionService.RequestCreateDmo(reqInputCreateUploadDmo)

		if requestJobDmoErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["dmo_period"] = detailDmo.Detail.Period
			inputMap["list_transaction_barge"] = inputUpdateDmo.TransactionBarge
			inputMap["list_grouping_vessel"] = inputUpdateDmo.GroupingVessel
			inputMap["input"] = inputUpdateDmo
			inputMap["input_job"] = reqInputCreateUploadDmo
			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": requestJobDmoErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
				DmoId:   &idDmo,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "Dmo already has been updated, but job failed",
				"error":   requestJobDmoErr.Error(),
				"dmo":     updateDmo,
			})
		}
	}
	return c.Status(200).JSON(updateDmo)
}

func (h *dmoHandler) ListDmo(c *fiber.Ctx) error {
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

	var filterDmo dmo.FilterAndSortDmo

	filterDmo.Quantity = c.Query("quantity")
	filterDmo.Month = c.Query("month")
	filterDmo.Year = c.Query("year")
	filterDmo.BuyerId = c.Query("buyer_id")
	filterDmo.Field = c.Query("field")
	filterDmo.Sort = c.Query("sort")

	listDmo, listDmoErr := h.dmoService.GetListReportDmoAll(pageNumber, filterDmo, iupopkIdInt)

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

	detailDmo, detailDmoErr := h.transactionService.GetDetailDmo(idInt, iupopkIdInt)

	if detailDmoErr != nil {
		status := 400

		if detailDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailDmoErr.Error(),
		})
	}

	return c.Status(200).JSON(detailDmo)
}

func (h *dmoHandler) ListTransactionForDmo(c *fiber.Ctx) error {
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

	var listTransactionForDmo transaction.ChooseTransactionDmo

	listBarge, listBargeErr := h.transactionService.ListDataDNBargeWithoutVessel(iupopkIdInt)

	if listBargeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listBargeErr.Error(),
		})
	}

	listTransactionForDmo.BargeTransaction = listBarge

	listGroupingVessel, listGroupingVesselErr := h.dmoVesselService.ListGroupingVesselWithoutDmo(iupopkIdInt)

	if listGroupingVesselErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listGroupingVesselErr.Error(),
		})
	}

	listTransactionForDmo.GroupingVesselTransaction = listGroupingVessel
	return c.Status(200).JSON(listTransactionForDmo)
}

func (h *dmoHandler) DeleteDmo(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error":   "record not found",
		})
	}

	findDmo, findDmoErr := h.transactionService.GetDetailDmo(idInt, iupopkIdInt)

	if findDmoErr != nil {
		status := 400

		if findDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error":   findDmoErr.Error(),
		})
	}

	if findDmo.Detail.IsBastDocumentSigned != false || findDmo.Detail.IsReconciliationLetterSigned != false || findDmo.Detail.IsStatementLetterSigned != false || findDmo.Detail.IsReconciliationLetterEndUserSigned != false {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error":   "document dmo is already signed",
		})
	}

	documentLink := findDmo.Detail.ReconciliationLetterDocumentLink

	if documentLink != nil {
		documentLinkSplit := strings.Split(*documentLink, "/")

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
			inputMap["dmo_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error":     deleteAwsErr.Error(),
				"id_number": findDmo.Detail.IdNumber,
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete dmo aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

	_, deleteDmoErr := h.historyService.DeleteDmo(idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if deleteDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId:   &findDmo.Detail.ID,
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete dmo",
			"error":   deleteDmoErr.Error(),
		})
	}

	var createNotif notification.InputNotification

	createNotif.Type = "dmo"
	createNotif.Status = "menghapus"
	createNotif.Period = findDmo.Detail.Period
	createNotif.EndUser = endUserDmo.Company.CompanyName
	_, createNotificationDeleteDmoErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationDeleteDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt
		inputMap["dmo_period"] = findDmo.Detail.Period
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
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

	inputUpdateDmo := new(dmo.InputUpdateDocumentDmo)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateDmo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update dmo",
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update dmo",
			"error":   "record not found",
		})
	}

	errors := h.v.Struct(*inputUpdateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailDmo, detailDmoErr := h.transactionService.GetDetailDmo(idInt, iupopkIdInt)

	if detailDmoErr != nil {
		status := 400

		if detailDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailDmoErr.Error(),
		})
	}

	updateDocumentDmo, updateDocumentDmoErr := h.historyService.UpdateDocumentDmo(idInt, *inputUpdateDmo, uint(claims["id"].(float64)), iupopkIdInt)

	if updateDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateDocumentDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateDocumentDmoErr.Error(),
			"message": "failed to update document dmo",
		})
	}

	getEndUserDmo, getEndUserDmoErr := h.traderDmoService.GetTraderEndUserDmo(idInt)

	if getEndUserDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": getEndUserDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if getEndUserDmoErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   getEndUserDmoErr.Error(),
			"message": "failed to get end user dmo",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"

	inputNotification.Status = "membuat"
	inputNotification.EndUser = getEndUserDmo.Company.CompanyName
	inputNotification.Period = detailDmo.Detail.Period
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateDmo

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
			"message": "failed to create notification update document dmo",
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

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error":   "record not found",
		})
	}

	var isBast bool
	var isStatementLetter bool
	var isReconciliationLetter bool
	var isReconciliationLetterEndUser bool

	typeDocument := c.Params("type")

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" && typeDocument != "reconciliation_letter_end_user" {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"message": "failed to update downloaded dmo",
			"error":   "type not found",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error":   "type not found",
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

	if typeDocument == "reconciliation_letter_end_user" {
		isReconciliationLetterEndUser = true
	}

	updateDownloadedDocumentDmo, updateDownloadedDocumentDmoErr := h.historyService.UpdateIsDownloadedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if updateDownloadedDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateDownloadedDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateDownloadedDocumentDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update downloaded dmo",
			"error":   updateDownloadedDocumentDmoErr.Error(),
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

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   "record not found",
		})
	}

	var isBast bool
	var isStatementLetter bool
	var isReconciliationLetter bool
	var isReconciliationLetterEndUser bool

	dataDmo, dataDmoErr := h.dmoService.GetDataDmo(idInt, iupopkIdInt)

	if dataDmoErr != nil {
		status := 400

		if dataDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   dataDmoErr.Error(),
		})
	}

	var fileName string

	idNumber := *dataDmo.IdNumber
	fileName = dataDmo.Iupopk.Code
	fileName += "/LBU/"
	fileName += idNumber

	file, errFormFile := c.FormFile("document")

	if errFormFile != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   errFormFile.Error(),
		})
	}

	typeDocument := c.Params("type")

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" && typeDocument != "reconciliation_letter_end_user" {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"message": "failed to update signed dmo",
			"error":   "type not found",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   "type not found",
		})
	}

	if typeDocument == "bast" {
		isBast = true
		fileName += fmt.Sprintf("/%v_signed_bast.pdf", idNumber)
	}

	if typeDocument == "statement_letter" {
		isStatementLetter = true
		fileName += fmt.Sprintf("/%v_signed_surat_pernyataan.pdf", idNumber)
	}

	if typeDocument == "reconciliation_letter" {
		isReconciliationLetter = true
		fileName += fmt.Sprintf("/%v_signed_berita_acara.pdf", idNumber)
	}

	if typeDocument == "reconciliation_letter_end_user" {
		isReconciliationLetterEndUser = true
		fileName += fmt.Sprintf("/%v_signed_berita_acara_pengguna_akhir.pdf", idNumber)
	}

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = typeDocument
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId:   &dmoId,
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   uploadErr.Error(),
		})
	}

	updateSignedDocumentDmo, updateSignedDocumentDmoErr := h.historyService.UpdateTrueIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, idInt, uint(claims["id"].(float64)), up.Location, iupopkIdInt)

	if updateSignedDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateSignedDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateSignedDocumentDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   updateSignedDocumentDmoErr.Error(),
		})
	}

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"
	inputNotification.Status = "mengupload"
	inputNotification.Period = dataDmo.Period
	inputNotification.Document = typeDocument
	inputNotification.EndUser = endUserDmo.Company.CompanyName

	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to create notification update signed dmo",
			"error":   createdNotificationErr.Error(),
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

	dmoId := uint(idInt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   "record not found",
		})
	}

	var isBast bool
	var isStatementLetter bool
	var isReconciliationLetter bool
	var isReconciliationLetterEndUser bool

	dataDmo, dataDmoErr := h.dmoService.GetDataDmo(idInt, iupopkIdInt)

	if dataDmoErr != nil {
		status := 400

		if dataDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   dataDmoErr.Error(),
		})
	}

	var fileName string

	fileName = dataDmo.Iupopk.Code
	fileName += "/LBU/"
	fileName += *dataDmo.IdNumber
	fileName += "/"
	fileName += *dataDmo.IdNumber
	fileName += "_"
	typeDocument := c.Params("type")

	if typeDocument != "bast" && typeDocument != "statement_letter" && typeDocument != "reconciliation_letter" && typeDocument != "reconciliation_letter_end_user" {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"message": "failed to update signed dmo",
			"error":   "type not found",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   "type not found",
		})
	}

	if typeDocument == "bast" {
		isBast = true
		fileName += "signed_bast.pdf"
	}

	if typeDocument == "statement_letter" {
		isStatementLetter = true
		fileName += "signed_surat_pernyataan.pdf"
	}

	if typeDocument == "reconciliation_letter" {
		isReconciliationLetter = true
		fileName += "signed_berita_acara.pdf"
	}

	if typeDocument == "reconciliation_letter_end_user" {
		isReconciliationLetterEndUser = true
		fileName += "signed_berita_acara_pengguna_akhir.pdf"
	}

	deleteUpload, deleteUploadErr := awshelper.DeleteDocument(fileName)

	if deleteUploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["document_type"] = typeDocument
		inputMap["user_id"] = claims["id"]
		inputMap["dmo_id"] = idInt
		inputMap["response"] = deleteUpload
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteUploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			DmoId:   &dmoId,
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   deleteUploadErr.Error(),
		})
	}

	updateSignedDocumentDmo, updateSignedDocumentDmoErr := h.historyService.UpdateFalseIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if updateSignedDocumentDmoErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateSignedDocumentDmoErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateSignedDocumentDmoErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update signed dmo",
			"error":   updateSignedDocumentDmoErr.Error(),
		})
	}

	endUserDmo, _ := h.traderDmoService.GetTraderEndUserDmo(idInt)

	var inputNotification notification.InputNotification
	inputNotification.Type = "dmo"
	inputNotification.Status = "menghapus"
	inputNotification.Period = dataDmo.Period
	inputNotification.Document = typeDocument
	inputNotification.EndUser = endUserDmo.Company.CompanyName

	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["type"] = typeDocument
		inputMap["dmo_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			DmoId:   &dmoId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"message": "failed to create notification delete signed dmo",
			"error":   createdNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(updateSignedDocumentDmo)
}
