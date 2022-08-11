package handler

import (
	"ajebackend/model/transaction"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type transactionHandler struct {

}

func NewTransactionHandler() *transactionHandler {
	return &transactionHandler {
	}
}

func (h *transactionHandler) HelloWorld (c *fiber.Ctx) error {
	dbUrlStg := "host=localhost user=postgres password=postgres dbname=deli_aje_development port=5432 sslmode=disable TimeZone=Asia/Shanghai"


	db, err := gorm.Open(postgres.Open(dbUrlStg), &gorm.Config{})

	if err != nil {
		return c.Status(400).JSON(err)
	}

	var testTransaction transaction.Test
	testTransaction.Date = "2022-10-10"
	createErr := db.Create(&testTransaction).Error

	fmt.Println(createErr)

	a := map[string]interface{}{
		"msg": "Hello World",
	}
	return c.Status(200).JSON(a)
}
