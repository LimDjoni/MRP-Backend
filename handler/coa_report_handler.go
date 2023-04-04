package handler

import (
	"ajebackend/model/coareport"
	"ajebackend/model/logs"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type coaReportHandler struct {
	coaReportService  coareport.Service
	logService        logs.Service
	userIupopkService useriupopk.Service
	v                 *validator.Validate
}

func NewCoaReportHandler(coaReportService coareport.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	v *validator.Validate,
) *coaReportHandler {
	return &coaReportHandler{
		coaReportService,
		logService,
		userIupopkService,
		v,
	}
}

func (h *coaReportHandler) ListCoaReportTransaction(c *fiber.Ctx) error {
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

	coaReportInput := new(coareport.CoaReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	listTransaction, listTransactionErr := h.coaReportService.GetTransactionCoaReport(coaReportInput.DateFrom, coaReportInput.DateTo, iupopkIdInt)

	if listTransactionErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": listTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(listTransaction)
}

func (h *coaReportHandler) CreateCoaReport(c *fiber.Ctx) error {
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

	coaReportInput := new(coareport.CoaReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(coaReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*coaReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	listTransaction, listTransactionErr := h.coaReportService.GetTransactionCoaReport(coaReportInput.DateFrom, coaReportInput.DateTo, iupopkIdInt)

	if listTransactionErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": listTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(listTransaction)
}
