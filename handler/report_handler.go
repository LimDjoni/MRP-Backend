package handler

import (
"ajebackend/model/logs"
"ajebackend/model/transaction"
"ajebackend/model/user"
"ajebackend/validatorfunc"
"fmt"
"github.com/go-playground/validator/v10"
"github.com/gofiber/fiber/v2"
"github.com/golang-jwt/jwt/v4"
"reflect"
"time"
)

type reportHandler struct {
	transactionService transaction.Service
	userService user.Service
	v               *validator.Validate
	logService logs.Service
}

func NewReportHandler(transactionService transaction.Service, userService user.Service, v *validator.Validate, logService logs.Service) *reportHandler {
	return &reportHandler{
		transactionService,
		userService,
		v,
		logService,
	}
}

func (h *reportHandler) Report(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64  {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	reportInput := new(transaction.InputRequestGetReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(reportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*reportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	if reportInput.Year == 0 {
		year, _, _ := time.Now().Date()
		reportInput.Year = year
	}

	reportRecap, reportDetail, reportErr := h.transactionService.GetReport(reportInput.Year)

	if reportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error": reportErr.Error(),
		})
	}

	reportRecap.Year = reportInput.Year
	reportRecap.ProductionPlan = reportInput.ProductionPlan
	reportRecap.PercentageProductionObligation = reportInput.PercentageProductionObligation
	reportRecap.ProductionObligation = reportInput.ProductionPlan * reportInput.PercentageProductionObligation / 100

	percentage := reportRecap.Total / reportRecap.ProductionObligation * 100
	reportRecap.FulfillmentPercentageProductionObligation = fmt.Sprintf("%.2f%%", percentage)
	reportRecap.ProrateProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total / reportRecap.ProductionPlan * 100)
	reportRecap.FulfillmentOfProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total / reportRecap.ProductionPlan * 100)


	return c.Status(200).JSON(fiber.Map{
		"detail": reportDetail,
		"recap": reportRecap,
	})
}
