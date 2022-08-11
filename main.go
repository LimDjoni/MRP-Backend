package main

import (
	"ajebackend/handler"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	transactionHandler := handler.NewTransactionHandler()

	app := fiber.New()

	apiV1 := app.Group("/api/v1") // /api

	apiV1.Get("/list", transactionHandler.HelloWorld)

	log.Fatal(app.Listen(":3000"))
}