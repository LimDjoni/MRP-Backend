package handler

import (
	"ajebackend/helper"
	"ajebackend/model/awshelper"
	"ajebackend/model/dmo"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
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
	v *validator.Validate
}

func NewDmoHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, dmoService dmo.Service, traderService trader.Service, v *validator.Validate) *dmoHandler {
	return &dmoHandler{
		transactionService,
		userService,
		historyService,
		logService,
		dmoService,
		traderService,
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

	var listTrader []uint

	for _, value := range inputCreateDmo.Trader {
		listTrader = append(listTrader, value.ID)
	}

	_, checkListTraderErr := h.traderService.CheckListTrader(listTrader)

	if checkListTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":  "trader " + checkListTraderErr.Error(),
		})
	}

	_, checkEndUserErr := h.traderService.CheckEndUser(inputCreateDmo.EndUser.ID)

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

	if len(inputCreateDmo.TransactionBarge) > 0 {
		_, checkDmoBargeErr := h.transactionService.CheckDataDnAndDmo(inputCreateDmo.TransactionBarge)

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
	}

	if len(inputCreateDmo.TransactionVessel) > 0 {
		_, checkDmoVesselErr := h.transactionService.CheckDataDnAndDmo(inputCreateDmo.TransactionVessel)

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

	if findDmo.Detail.ReconciliationLetterDocumentLink != nil || findDmo.Detail.BASTDocumentLink != nil || findDmo.Detail.StatementLetterDocumentLink != nil {
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
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete dmo",
	})
}
