package handler

import (
	"ajebackend/model/transaction"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{
		transactionService,
	}
}

func (h *transactionHandler) CreateTransactionDN(c *fiber.Ctx) error {
	transactionInput := new(transaction.DataTransactionInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(transactionInput); err != nil {
		return c.Status(400).JSON(err)
	}

	createdTransaction, createdTransactionErr := h.transactionService.CreateTransactionDN(*transactionInput)
	//response := map[string]interface{}{}
	//
	if createdTransactionErr != nil {
 		return c.Status(400).JSON(createdTransactionErr.Error())
	}

	return c.Status(201).JSON(createdTransaction)
}

func (h *transactionHandler) ListDataDN(c *fiber.Ctx) error {
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
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		response := map[string]interface{}{
			"error": "data tidak ditemukan",
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

func (h *transactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		response := map[string]interface{}{
			"error": "data tidak ditemukan",
		}
		return c.Status(404).JSON(response)
	}

	deleteTransaction, deleteTransactionErr := h.transactionService.DeleteTransaction(idInt)

	if deleteTransactionErr != nil {
		response := map[string]interface{}{
			"error": deleteTransactionErr.Error(),
		}
		return c.Status(400).JSON(response)
	}

	if deleteTransaction == false {
		response := map[string]interface{}{
			"message": "data tidak terhapus",
		}
		return c.Status(400).JSON(response)
	}

	response := map[string]interface{}{
		"message": "data berhasil dihapus",
	}
	return c.Status(200).JSON(response)
}
