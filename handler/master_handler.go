package handler

import (
	"ajebackend/model/counter"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type masterHandler struct {
	destinationService destination.Service
	userService        user.Service
	allMasterService   allmaster.Service
	counterService     counter.Service
	v                  *validator.Validate
}

func NewMasterHandler(destinationService destination.Service, userService user.Service, allMasterService allmaster.Service, counterService counter.Service, v *validator.Validate) *masterHandler {
	return &masterHandler{
		destinationService,
		userService,
		allMasterService,
		counterService,
		v,
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

func (h *masterHandler) GetListMaster(c *fiber.Ctx) error {
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

	listMaster, listMasterErr := h.allMasterService.ListMasterData()

	if listMasterErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listMasterErr.Error(),
		})
	}

	return c.Status(200).JSON(listMaster)
}

func (h *masterHandler) UpdateCounter(c *fiber.Ctx) error {
	updateCounterErr := h.counterService.UpdateCounter()

	if updateCounterErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": updateCounterErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "counter has been updated",
	})
}

func (h *masterHandler) CreateIupopk(c *fiber.Ctx) error {
	iupopkInput := new(iupopk.InputIupopk)

	// Binds the request body to the Person struct
	if err := c.BodyParser(iupopkInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*iupopkInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdIupopk, createdIupopkErr := h.counterService.CreateIupopk(*iupopkInput)

	if createdIupopkErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdIupopkErr.Error(),
		})
	}

	return c.Status(200).JSON(createdIupopk)
}
