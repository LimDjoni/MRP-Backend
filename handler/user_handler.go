package handler

import (
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type userHandler struct {
	userService user.Service
	v               *validator.Validate
}

func NewUserHandler(userService user.Service, v *validator.Validate) *userHandler {
	return &userHandler{
		userService,
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

	loginUser, loginUserErr := h.userService.LoginUser(*loginInput)

	if loginUserErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "wrong email / username / password",
		})
	}

	return c.Status(200).JSON(loginUser)
}
