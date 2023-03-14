package handler

import (
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type notificationHandler struct {
	userService             user.Service
	notificationService     notification.Service
	notificationUserService notificationuser.Service
	logsService             logs.Service
	v                       *validator.Validate
}

func NewNotificationHandler(userService user.Service, notificationService notification.Service, notificationUserService notificationuser.Service, logsService logs.Service, v *validator.Validate) *notificationHandler {
	return &notificationHandler{
		userService,
		notificationService,
		notificationUserService,
		logsService,
		v,
	}
}

func (h *notificationHandler) CreateNotification(c *fiber.Ctx) error {
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

	inputCreateNotification := new(notification.InputNotification)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateNotification); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateNotification)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateNotification
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors":  dataErrors,
			"message": "failed to create notification",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdNotification, createdNotificationErr := h.notificationUserService.CreateNotification(*inputCreateNotification, uint(claims["id"].(float64)), 11)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateNotification
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors":  createdNotificationErr.Error(),
			"message": "failed to create notification",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logsService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createdNotificationErr.Error(),
		})
	}

	return c.Status(201).JSON(createdNotification)
}

func (h *notificationHandler) GetNotification(c *fiber.Ctx) error {
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

	listNotification, listNotificationErr := h.notificationUserService.GetNotification(uint(claims["id"].(float64)))

	if listNotificationErr != nil && listNotificationErr.Error() != "record not found" {
		return c.Status(400).JSON(fiber.Map{
			"error": listNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": listNotification,
	})
}

func (h *notificationHandler) UpdateNotification(c *fiber.Ctx) error {
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

	updatedNotification, updatedNotificationErr := h.notificationUserService.UpdateReadNotification(uint(claims["id"].(float64)))

	if updatedNotificationErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   updatedNotificationErr.Error(),
			"message": "failed to update notification",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"list": updatedNotification,
	})
}

func (h *notificationHandler) DeleteNotification(c *fiber.Ctx) error {
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

	_, deletedNotificationErr := h.notificationService.DeleteNotification(uint(claims["id"].(float64)))

	if deletedNotificationErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deletedNotificationErr.Error(),
			"message": "failed to delete notification",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete notification",
	})
}
