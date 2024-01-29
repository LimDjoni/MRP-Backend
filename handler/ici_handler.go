package handler

import (
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/ici"
	"ajebackend/validatorfunc"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type iciHandler struct {
	iciService     ici.Service
	historyService history.Service
	logService     logs.Service
	v              *validator.Validate
}

func NewIciHandler(iciService ici.Service, historyService history.Service, logService logs.Service, v *validator.Validate) *iciHandler {
	return &iciHandler{
		iciService:     iciService,
		historyService: historyService,
		logService:     logService,
		v:              v,
	}
}

func (h *iciHandler) CreateIci(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputIci := new(ici.InputCreateUpdateIci)

	if err := c.BodyParser(inputIci); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputIci)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)

	}

}
