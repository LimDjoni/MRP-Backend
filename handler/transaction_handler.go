package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/destination"
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

type transactionHandler struct {
	transactionService transaction.Service
	historyService     history.Service
	v                  *validator.Validate
	logService         logs.Service
	destinationService destination.Service
	userIupopkService  useriupopk.Service
}

func NewTransactionHandler(transactionService transaction.Service, historyService history.Service, v *validator.Validate, logService logs.Service, destinationService destination.Service, userIupopkService useriupopk.Service) *transactionHandler {
	return &transactionHandler{
		transactionService,
		historyService,
		v,
		logService,
		destinationService,
		userIupopkService,
	}
}

func (h *transactionHandler) CreateTransactionDN(c *fiber.Ctx) error {
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

	var uniqueErrors []transaction.ErrorResponseUnique

	isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique := h.transactionService.CheckDataUnique(*transactionInput)

	if isDpRoyaltyNtpnUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "dp_royalty_ntpn"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.DpRoyaltyNtpn
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isDpRoyaltyBillingCodeUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "dp_royalty_billing_code"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.DpRoyaltyBillingCode
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isPaymentDpRoyaltyNtpnUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "payment_dp_royalty_ntpn"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.PaymentDpRoyaltyNtpn
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isPaymentDpRoyaltyBillingCodeUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "payment_dp_royalty_billing_code"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.PaymentDpRoyaltyBillingCode
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if len(uniqueErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": uniqueErrors,
		})
	}

	createdTransaction, createdTransactionErr := h.historyService.CreateTransactionDN(*transactionInput, uint(claims["id"].(float64)), iupopkIdInt)

	if createdTransactionErr != nil {
		inputJson, _ := json.Marshal(transactionInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdTransactionErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createdTransactionErr.Error(),
		})
	}

	return c.Status(201).JSON(createdTransaction)
}

func (h *transactionHandler) ListData(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if iupopkErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkErr.Error(),
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	var sortAndFilter transaction.SortAndFilter
	page := c.Query("page")
	sortAndFilter.Field = c.Query("field")
	sortAndFilter.Sort = c.Query("sort")
	sortAndFilter.Quantity = c.Query("quantity")
	sortAndFilter.TugboatId = c.Query("tugboat_id")
	sortAndFilter.BargeId = c.Query("barge_id")
	sortAndFilter.VesselId = c.Query("vessel_id")
	sortAndFilter.ShippingStart = c.Query("shipping_start")
	sortAndFilter.ShippingEnd = c.Query("shipping_end")
	sortAndFilter.VerificationFilter = c.Query("verification_filter")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if page == "" {
		pageNumber = 1
	}

	typeTransaction := strings.ToUpper(c.Params("transaction_type"))
	listDN, listDNErr := h.transactionService.ListData(pageNumber, sortAndFilter, typeTransaction, iupopkIdInt)

	if listDNErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDNErr.Error(),
		})
	}

	return c.Status(200).JSON(listDN)
}

func (h *transactionHandler) DetailTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")
	typeTransaction := strings.ToUpper(c.Params("transaction_type"))
	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailTransactionDN, detailTransactionDNErr := h.transactionService.DetailTransaction(idInt, typeTransaction, iupopkIdInt)

	if detailTransactionDNErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailTransactionDNErr.Error(),
		})
	}

	return c.Status(200).JSON(detailTransactionDN)
}

func (h *transactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")
	typeTransaction := strings.ToUpper(c.Params("transaction_type"))
	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}
	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	findTransaction, findTransactionErr := h.transactionService.DetailTransaction(idInt, typeTransaction, iupopkIdInt)

	if findTransactionErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": findTransactionErr.Error(),
		})
	}

	if findTransaction.MinerbaId != nil || findTransaction.DmoId != nil || findTransaction.MinerbaLnId != nil || findTransaction.GroupingVesselLnId != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "can't delete transaction because it's bound to minerba or dmo or group vessel",
		})
	}

	_, deleteTransactionErr := h.historyService.DeleteTransaction(idInt, uint(claims["id"].(float64)), typeTransaction, iupopkIdInt)

	if deleteTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteTransactionErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input:         inputJson,
			Message:       messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete transaction",
			"error":   deleteTransactionErr.Error(),
		})
	}

	documentLink := findTransaction.SkbDocumentLink

	if findTransaction.SkbDocumentLink != "" || findTransaction.SkabDocumentLink != "" || findTransaction.BLDocumentLink != "" || findTransaction.RoyaltiProvisionDocumentLink != "" || findTransaction.RoyaltiFinalDocumentLink != "" || findTransaction.COWDocumentLink != "" || findTransaction.COADocumentLink != "" || findTransaction.InvoiceAndContractDocumentLink != "" || findTransaction.LHVDocumentLink != "" {
		if findTransaction.SkabDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.SkabDocumentLink
		}
		if findTransaction.BLDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.BLDocumentLink
		}
		if findTransaction.RoyaltiProvisionDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.RoyaltiProvisionDocumentLink
		}
		if findTransaction.RoyaltiFinalDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.RoyaltiFinalDocumentLink
		}
		if findTransaction.COWDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.COWDocumentLink
		}
		if findTransaction.COADocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.COADocumentLink
		}
		if findTransaction.InvoiceAndContractDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.InvoiceAndContractDocumentLink
		}
		if findTransaction.LHVDocumentLink != "" && documentLink == "" {
			documentLink = findTransaction.LHVDocumentLink
		}

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
			inputMap["transaction_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error":     deleteAwsErr.Error(),
				"id_number": findTransaction.IdNumber,
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete transaction aws",
				"error":   deleteAwsErr.Error(),
			})
		}
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

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")

	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
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

	updateTransaction, updateTransactionErr := h.historyService.UpdateTransactionDN(idInt, *transactionInput, uint(claims["id"].(float64)), iupopkIdInt)

	if updateTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = transactionInput
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = id

		inputJson, _ := json.Marshal(transactionInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateTransactionErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input:         inputJson,
			Message:       messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update transaction",
			"error":   updateTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTransaction)
}

func (h *transactionHandler) UpdateDocumentTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")
	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	documentType := c.Params("type")
	responseErr := fiber.Map{
		"message": "failed to upload document",
	}

	file, errFormFile := c.FormFile("document")

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	switch documentType {
	case "skb", "skab", "bl", "royalti_provision", "royalti_final", "cow", "coa", "invoice", "lhv":
	default:
		return c.Status(400).JSON(fiber.Map{
			"error":   "document type not found",
			"message": "failed to upload document",
		})
	}

	typeTransaction := strings.ToUpper(c.Params("transaction_type"))
	detailTransaction, detailTransactionErr := h.transactionService.DetailTransaction(idInt, typeTransaction, iupopkIdInt)

	if detailTransactionErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to upload document",
			"error":   detailTransactionErr.Error(),
		})
	}

	fileName := fmt.Sprintf("AJE/T%s/%v/%v_%s.pdf", typeTransaction, *detailTransaction.IdNumber, *detailTransaction.IdNumber, documentType)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input:         inputJson,
			Message:       messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	editDocument, editDocumentErr := h.historyService.UploadDocumentTransaction(detailTransaction.ID, up.Location, uint(claims["id"].(float64)), documentType, typeTransaction, iupopkIdInt)

	if editDocumentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           editDocumentErr.Error(),
			"upload_response": up,
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input:         inputJson,
			Message:       messageJson,
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

func (h *transactionHandler) CreateTransactionLN(c *fiber.Ctx) error {
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

	var uniqueErrors []transaction.ErrorResponseUnique

	isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique := h.transactionService.CheckDataUnique(*transactionInput)
	if isDpRoyaltyNtpnUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "dp_royalty_ntpn"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.DpRoyaltyNtpn
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isDpRoyaltyBillingCodeUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "dp_royalty_billing_code"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.DpRoyaltyBillingCode
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isPaymentDpRoyaltyNtpnUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "payment_dp_royalty_ntpn"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.PaymentDpRoyaltyNtpn
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if isPaymentDpRoyaltyBillingCodeUnique == true {
		var uniqueErr transaction.ErrorResponseUnique
		uniqueErr.FailedField = "payment_dp_royalty_billing_code"
		uniqueErr.Tag = "unique"
		uniqueErr.Value = *transactionInput.PaymentDpRoyaltyBillingCode
		uniqueErrors = append(uniqueErrors, uniqueErr)
	}

	if len(uniqueErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": uniqueErrors,
		})
	}

	createdTransactionLN, createdTransactionLNErr := h.historyService.CreateTransactionLN(*transactionInput, uint(claims["id"].(float64)), iupopkIdInt)

	if createdTransactionLNErr != nil {
		inputJson, _ := json.Marshal(transactionInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdTransactionLNErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createdTransactionLNErr.Error(),
		})
	}

	return c.Status(201).JSON(createdTransactionLN)
}

func (h *transactionHandler) UpdateTransactionLN(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")
	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
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

	updateTransaction, updateTransactionErr := h.historyService.UpdateTransactionLN(idInt, *transactionInput, uint(claims["id"].(float64)), iupopkIdInt)

	if updateTransactionErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = transactionInput
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = id

		inputJson, _ := json.Marshal(transactionInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateTransactionErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionId: &transactionId,
			Input:         inputJson,
			Message:       messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateTransactionErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update transaction",
			"error":   updateTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTransaction)
}
