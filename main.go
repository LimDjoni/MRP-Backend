package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	apiV1 := app.Group("/api/v1") // /api

	apiV1.Get("/list", handler)

	log.Fatal(app.Listen(":3000"))
}