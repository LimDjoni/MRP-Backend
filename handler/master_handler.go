package handler

import (
	"ajebackend/model/trader"
	"ajebackend/model/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
)

type masterHandler struct {
	userService user.Service
	traderService trader.Service
}

func NewMasterHandler(userService user.Service, traderService trader.Service) *masterHandler {
	return &masterHandler{
		userService,
		traderService,
	}
}

func (h *masterHandler) ListTrader(c *fiber.Ctx) error {
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
	return c.Status(200).JSON(listTrader)
}
