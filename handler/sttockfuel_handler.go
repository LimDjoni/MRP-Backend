package handler

import (
	"mrpbackend/model/stockfuel"
	"mrpbackend/model/user"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type stockfuelsHandler struct {
	userService       user.Service
	stockfuelsService stockfuel.Service
	v                 *validator.Validate
}

func NewStockFuelHandler(userService user.Service, stockfuelsService stockfuel.Service, v *validator.Validate) *stockfuelsHandler {
	return &stockfuelsHandler{
		userService,
		stockfuelsService,
		v,
	}
}

// Master Data
func (h *stockfuelsHandler) ListStockFuel(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	var filterStockFuel stockfuel.StockFuelSummary

	filterStockFuel.Month = c.Query("month")

	listStockFuel, listStockFuelErr := h.stockfuelsService.ListStockFuel(filterStockFuel)

	if listStockFuelErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listStockFuelErr.Error(),
		})
	}

	return c.Status(200).JSON(listStockFuel)
}
