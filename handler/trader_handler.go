package handler

import (
	"ajebackend/model/company"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
)

type traderHandler struct {
	userService user.Service
	traderService trader.Service
	companyService company.Service
	logsService logs.Service
	v *validator.Validate
}

func NewTraderHandler(userService user.Service, traderService trader.Service, companyService company.Service, logsService logs.Service, v *validator.Validate) *traderHandler {
	return &traderHandler{
		userService,
		traderService,
		companyService,
		logsService,
		v,
	}
}

func (h *traderHandler) ListTrader(c *fiber.Ctx) error {
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

	listTrader, listTraderErr := h.traderService.ListTrader()

	if listTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listTraderErr.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"trader": listTrader,
	})
}
