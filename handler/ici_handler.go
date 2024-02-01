package handler

import (
	"ajebackend/model/ici"
	"ajebackend/model/logs"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type iciHandler struct {
	iciService        ici.Service
	logService        logs.Service
	v                 *validator.Validate
	userIupopkService useriupopk.Service
}

func NewIciHandler(iciService ici.Service, logService logs.Service, v *validator.Validate, userIupopkService useriupopk.Service) *iciHandler {
	return &iciHandler{
		iciService,
		logService,
		v,
		userIupopkService,
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

	// harus ada iupopk
	iupopkId := c.Params("iupopk_id")
	iupopkIdInt, err := strconv.Atoi(iupopkId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || !checkUser.IsActive {
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

	createIci, createIciErr := h.iciService.CreateIci(*inputIci, iupopkIdInt)

	if createIciErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createIciErr.Error(),
		})
	}

	return c.Status(201).JSON(createIci)
}
