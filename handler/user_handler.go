package handler

import (
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       interface{}
}

var validate = validator.New()

func ValidateStruct(value interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(value)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = fieldErr.StructNamespace()
			element.Tag = fieldErr.Tag()
			element.Value = fieldErr.Value()
			errors = append(errors, &element)
		}
	}
	return errors
}

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{
		userService,
	}
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	registerInput := new(user.RegisterUserInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(registerInput); err != nil {
		return c.Status(400).JSON(err)
	}

	errors := ValidateStruct(*registerInput)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	newUser, newUserErr := h.userService.RegisterUser(*registerInput)

	if newUserErr != nil {
		return c.Status(400).JSON(newUserErr.Error())
	}

	return c.Status(201).JSON(newUser)
}

func (h *userHandler) LoginUser(c *fiber.Ctx) error {
	loginInput := new(user.LoginUserInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(loginInput); err != nil {
		return c.Status(400).JSON(err)
	}

	errors := ValidateStruct(*loginInput)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	loginUser, loginUserErr := h.userService.LoginUser(*loginInput)

	if loginUserErr != nil {
		return c.Status(400).JSON(loginUserErr.Error())
	}

	return c.Status(200).JSON(loginUser)
}
