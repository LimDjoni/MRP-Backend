package handler

import (
	"ajebackend/model/user"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type userHandler struct {
	userService       user.Service
	userIupopkService useriupopk.Service
	v                 *validator.Validate
}

func NewUserHandler(userService user.Service, userIupopkService useriupopk.Service, v *validator.Validate) *userHandler {
	return &userHandler{
		userService,
		userIupopkService,
		v,
	}
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	registerInput := new(user.RegisterUserInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(registerInput); err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	errors := h.v.Struct(*registerInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})

	}

	newUser, newUserErr := h.userService.RegisterUser(*registerInput)

	if newUserErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": newUserErr.Error(),
		})
	}

	return c.Status(201).JSON(newUser)
}

func (h *userHandler) LoginUser(c *fiber.Ctx) error {
	loginInput := new(user.LoginUserInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(loginInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginUser, loginUserErr := h.userIupopkService.LoginUser(*loginInput)

	if loginUserErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": loginUserErr.Error(),
		})
	}

	return c.Status(200).JSON(loginUser)
}

func (h *userHandler) Validate(c *fiber.Ctx) error {
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

	return c.Status(200).JSON(fiber.Map{
		"message": "validate",
	})
}

func (h *userHandler) CreateUserIupopk(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("iupopk_id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	createUserIupopk, createUserIupopkErr := h.userIupopkService.CreateUserIupopk(int(claims["id"].(float64)), idInt)

	status := 400
	if createUserIupopkErr != nil {

		if createUserIupopkErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": createUserIupopkErr.Error(),
		})
	}

	return c.Status(201).JSON(createUserIupopk)
}

func (h *userHandler) DeleteUserIupopk(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("iupopk_id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	deleteUserIupopkErr := h.userIupopkService.DeleteUserIupopk(int(claims["id"].(float64)), idInt)

	status := 400
	if deleteUserIupopkErr != nil {

		if deleteUserIupopkErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": deleteUserIupopkErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "user iupopk has been deleted",
	})
}
