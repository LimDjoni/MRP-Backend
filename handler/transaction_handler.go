package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
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

type transactionHandler struct {
	transactionService transaction.Service
	userService user.Service
	historyService history.Service
	v               *validator.Validate
	logService logs.Service
}

func NewTransactionHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, v *validator.Validate, logService logs.Service) *transactionHandler {
	return &transactionHandler{
		transactionService,
		userService,
		historyService,
		v,
		logService,
	}
}

func (h *transactionHandler) CreateTransactionDN(c *fiber.Ctx) error {
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

	transactionInput := new(transaction.DataTransactionInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(transactionInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*transactionInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdTransaction, createdTransactionErr := h.historyService.CreateTransactionDN(*transactionInput, uint(claims["id"].(float64)))

	if createdTransactionErr != nil {
		inputJson , _ := json.Marshal(transactionInput)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": createdTransactionErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

 		return c.Status(400).JSON(fiber.Map{
			"error": createdTransactionErr.Error(),
		})
	}

	return c.Status(201).JSON(createdTransaction)
}

func (h *transactionHandler) ListDataDN(c *fiber.Ctx) error {
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

	var sortAndFilter transaction.SortAndFilter
	page := c.Query("page")
	sortAndFilter.Field = c.Query("field")
	sortAndFilter.Sort = c.Query("sort")

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		sortAndFilter.Quantity = 0
	} else {
		sortAndFilter.Quantity = quantity
	}

	sortAndFilter.ShipName = c.Query("ship_name")
	sortAndFilter.BargeName = c.Query("barge_name")
	sortAndFilter.ShippingFrom = c.Query("shipping_from")
	sortAndFilter.ShippingTo = c.Query("shipping_to")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if page == "" {
		pageNumber = 1
	}

	listDN, listDNErr := h.transactionService.ListDataDN(pageNumber, sortAndFilter)

	if listDNErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDNErr.Error(),
		})
	}

	return c.Status(200).JSON(listDN)
}

func (h *transactionHandler) DetailTransactionDN(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	detailTransactionDN, detailTransactionDNErr := h.transactionService.DetailTransactionDN(idInt)

	if detailTransactionDNErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailTransactionDNErr.Error(),
		})
	}

	return c.Status(200).JSON(detailTransactionDN)
}

func (h *transactionHandler) DeleteTransactionDN(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	deleteTransaction, deleteTransactionErr := h.historyService.DeleteTransactionDN(idInt, uint(claims["id"].(float64)))

	if deleteTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson ,_ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": deleteTransactionErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if  deleteTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete transaction",
			"error": deleteTransactionErr.Error(),
		})
	}

	if deleteTransaction == false && deleteTransactionErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to delete transaction",
			"error": deleteTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete transaction",
	})
}

func (h *transactionHandler) UpdateTransactionDN(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	transactionInput := new(transaction.DataTransactionInput)

	// Binds the request body to the Person struct
	if errParsing := c.BodyParser(transactionInput); errParsing != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errParsing.Error(),
		})
	}

	errors := h.v.Struct(*transactionInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateTransaction, updateTransactionErr := h.historyService.UpdateTransactionDN(idInt, *transactionInput ,uint(claims["id"].(float64)))

	if updateTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = transactionInput
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = id

		inputJson , _ := json.Marshal(transactionInput)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": updateTransactionErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if  updateTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update transaction",
			"error": updateTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTransaction)
}

func (h *transactionHandler) UpdateDocumentTransactionDN (c *fiber.Ctx) error {
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

	documentType := c.Params("type")
	responseErr := fiber.Map{
		"message": "failed to upload document",
	}

	file, errFormFile := c.FormFile("document")

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf" ) {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	switch documentType {
		case "skb","skab","bl","royalti_provision","royalti_final","cow","coa","invoice","lhv":
		default:
			return c.Status(400).JSON(fiber.Map{
				"error": "document type not found",
				"message": "failed to upload document",
			})
	}

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detailTransaction, detailTransactionErr := h.transactionService.DetailTransactionDN(idInt)

	if detailTransactionErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to upload document",
			"error": detailTransactionErr.Error(),
		})
	}

	fileName := fmt.Sprintf("%s/%s.pdf", detailTransaction.IdNumber, documentType)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson , _ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input: inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	editDocument, editDocumentErr := h.historyService.UploadDocumentTransactionDN(detailTransaction.ID, up.Location, uint(claims["id"].(float64)), documentType)

	if editDocumentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson , _ := json.Marshal(inputMap)
		messageJson ,_ := json.Marshal(map[string]interface{}{
			"error": editDocumentErr.Error(),
			"upload_response": up,
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input: inputJson,
			Message: messageJson,
		}
		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = editDocumentErr.Error()

		status := 400

		if editDocumentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(responseErr)
	}

	return c.Status(200).JSON(editDocument)
}
