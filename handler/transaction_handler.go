package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"strconv"
)

type transactionHandler struct {
	transactionService transaction.Service
	userService user.Service
	historyService history.Service
}

func NewTransactionHandler(transactionService transaction.Service, userService user.Service, historyService history.Service) *transactionHandler {
	return &transactionHandler{
		transactionService,
		userService,
		historyService,
	}
}

func (h *transactionHandler) CreateTransactionDN(c *fiber.Ctx) error {
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

	transactionInput := new(transaction.DataTransactionInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(transactionInput); err != nil {
		return c.Status(400).JSON(err)
	}

	createdTransaction, createdTransactionErr := h.historyService.CreateTransactionDN(*transactionInput, uint(claims["id"].(float64)))
	//response := map[string]interface{}{}
	//
	if createdTransactionErr != nil {
 		return c.Status(400).JSON(createdTransactionErr.Error())
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

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		response := map[string]interface{}{
			"error": err.Error(),
		}
		return c.Status(400).JSON(response)
	}
	if page == "" {
		pageNumber = 1
	}

	listDN, listDNErr := h.transactionService.ListDataDN(pageNumber)

	if listDNErr != nil {
		return c.Status(400).JSON(listDNErr.Error())
	}

	return c.Status(200).JSON(listDN)
}

func (h *transactionHandler) DetailTransactionDN(c *fiber.Ctx) error {
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
		response := map[string]interface{}{
			"error": "data not found",
		}
		return c.Status(404).JSON(response)
	}

	detailTransactionDN, detailTransactionDNErr := h.transactionService.DetailTransactionDN(idInt)

	if detailTransactionDNErr != nil {
		response := map[string]interface{}{
			"error": detailTransactionDNErr.Error(),
		}
		return c.Status(404).JSON(response)
	}

	return c.Status(200).JSON(detailTransactionDN)
}

func (h *transactionHandler) DeleteTransactionDN(c *fiber.Ctx) error {
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
		response := map[string]interface{}{
			"error": "data not found",
		}
		return c.Status(404).JSON(response)
	}

	deleteTransaction, deleteTransactionErr := h.historyService.DeleteTransactionDN(idInt, uint(claims["id"].(float64)))

	if deleteTransactionErr != nil {
		response := map[string]interface{}{
			"message": "failed to delete transaction",
			"error": deleteTransactionErr.Error(),
		}
		return c.Status(400).JSON(response)
	}

	if deleteTransaction == false && deleteTransactionErr != nil {
		response := map[string]interface{}{
			"message": "failed to delete transaction",
		}
		return c.Status(400).JSON(response)
	}

	response := map[string]interface{}{
		"message": "success delete transaction",
	}
	return c.Status(200).JSON(response)
}

func (h *transactionHandler) UpdateTransactionDN(c *fiber.Ctx) error {
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
		response := map[string]interface{}{
			"error": "data not found",
		}
		return c.Status(404).JSON(response)
	}

	transactionInput := new(transaction.DataTransactionInput)

	// Binds the request body to the Person struct
	if errParsing := c.BodyParser(transactionInput); errParsing != nil {
		return c.Status(400).JSON(errParsing)
	}


	updateTransaction, updateTransactionErr := h.historyService.UpdateTransactionDN(idInt, *transactionInput ,uint(claims["id"].(float64)))

	if updateTransactionErr != nil {
		response := map[string]interface{}{
			"message": "failed to update transaction",
			"error": updateTransactionErr.Error(),
		}
		return c.Status(400).JSON(response)
	}

	return c.Status(200).JSON(updateTransaction)
}

func (h *transactionHandler) UpdateDocumentTransactionDN (c *fiber.Ctx) error {
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

	documentType := c.Params("type")

	switch documentType {
		case "skb","skab","bl","royalti_provision","royalti_final","cow","coa","invoice","lhv":
		default:
			response := map[string]interface{}{
				"error": "document type not found",
			}
			return c.Status(400).JSON(response)
	}


	idInt, err := strconv.Atoi(id)

	if err != nil {
		response := map[string]interface{}{
			"error": "data not found",
		}
		return c.Status(400).JSON(response)
	}

	detailTransaction, detailTransactionErr := h.transactionService.DetailTransactionDN(idInt)

	if detailTransactionErr != nil {
		response := map[string]interface{}{
			"message": "failed to upload document",
			"error": detailTransactionErr.Error(),
		}
		return c.Status(400).JSON(response)
	}

	fileName := fmt.Sprintf("%s/%s.pdf", detailTransaction.IdNumber, documentType)

	file, errFormFile := c.FormFile("document")
	responseErr := map[string]interface{}{
		"message": "failed to upload document",
	}
	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	editDocument, editDocumentErr := h.historyService.UploadDocument(detailTransaction.ID, up.Location, uint(claims["id"].(float64)), documentType)

	if editDocumentErr != nil {
		responseErr["error"] = editDocumentErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	return c.Status(200).JSON(editDocument)
}
