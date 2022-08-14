package handler

import (
	"github.com/gofiber/fiber/v2"
)

type transactionHandler struct {
}

func NewTransactionHandler() *transactionHandler {
	return &transactionHandler{}
}

func (h *transactionHandler) HelloWorld(c *fiber.Ctx) error {

	a := map[string]interface{}{
		"msg": "Hello World",
	}
	return c.Status(200).JSON(a)
}
