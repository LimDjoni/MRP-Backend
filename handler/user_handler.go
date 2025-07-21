package handler

import (
	"mrpbackend/helper"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type userHandler struct {
	userService user.Service
	v           *validator.Validate
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

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	userchange := c.Locals("user").(*jwt.Token)

	claims := userchange.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	changePasswordInput := new(user.ChangePasswordInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(changePasswordInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	if changePasswordInput.NewPassword != changePasswordInput.ConfirmPassword || !helper.CheckPassword(changePasswordInput.OldPassword, checkUser.Password) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password tidak cocok",
		})
	}

	_, updUserErr := h.userService.ChangePassword(changePasswordInput.NewPassword, uint(claims["id"].(float64)))

	if updUserErr != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": updUserErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Sukses update password"})
}

func (h *userHandler) ResetPassword(c *fiber.Ctx) error {

	resetPasswordInput := new(user.ResetPasswordInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(resetPasswordInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newPassword := helper.CreateRandomPassword()

	_, resetUserErr := h.userService.ResetPassword(resetPasswordInput.Email, newPassword)

	if resetUserErr != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": resetUserErr.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": newPassword,
	})
}
