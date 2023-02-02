package handler

import (
	"ajebackend/model/master/destination"
	"ajebackend/model/user"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type masterHandler struct {
	destinationService destination.Service
	userService        user.Service
}

func NewMasterHandler(destinationService destination.Service, userService user.Service) *masterHandler {
	return &masterHandler{
		destinationService,
		userService,
	}
}

func (h *masterHandler) GetDestination(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	destinations, destinationsErr := h.destinationService.GetDestination()

	if destinationsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": destinationsErr.Error(),
		})
	}

	return c.Status(200).JSON(destinations)
}
