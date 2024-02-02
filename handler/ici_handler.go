package handler

import (
	"ajebackend/model/ici"
	"ajebackend/model/logs"
	"ajebackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type iciHandler struct {
	iciService ici.Service
	logService logs.Service
	v          *validator.Validate
}

func NewIciHandler(iciService ici.Service, logService logs.Service, v *validator.Validate) *iciHandler {
	return &iciHandler{
		iciService,
		logService,
		v,
	}
}

func (h *iciHandler) CreateIci(c *fiber.Ctx) error {

	//Get User
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	// Check User Login
	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	//Get input ICI
	inputIci := new(ici.InputCreateUpdateIci)

	//Check Input ICI
	if err := c.BodyParser(inputIci); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputIci)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createIci, createIciErr := h.iciService.CreateIci(*inputIci)

	if createIciErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createIciErr.Error(),
		})
	}

	return c.Status(201).JSON(createIci)
}

func (h *iciHandler) UpdateIci(c *fiber.Ctx) error {
	//Get User
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}
	// Check User Login
	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	//Get Input
	inputUpdateIci := new(ici.InputCreateUpdateIci)
	if err := c.BodyParser(inputUpdateIci); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateIci)
	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	//Get ID
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	updateIci, err := h.iciService.UpdateIci(*inputUpdateIci, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateIci)
}
