package handler

import (
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerbaln"
	"ajebackend/model/notificationuser"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type minerbaLnHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	logService              logs.Service
	minerbaLnService        minerbaln.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}

func NewMinerbaLnHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, logService logs.Service, minerbaLnService minerbaln.Service, notificationUserService notificationuser.Service, v *validator.Validate) *minerbaLnHandler {
	return &minerbaLnHandler{
		transactionService,
		userService,
		historyService,
		logService,
		minerbaLnService,
		notificationUserService,
		v,
	}
}

func (h *minerbaLnHandler) CreateMinerbaLn(c *fiber.Ctx) error {
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

	inputCreateMinerbaLn := new(minerbaln.InputCreateMinerbaLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateMinerbaLn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateMinerbaLn)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_ln_period"] = inputCreateMinerbaLn.Period
		inputMap["list_transactions"] = inputCreateMinerbaLn.ListDataLn
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	// _, findMinerbaErr := h.minerbaService.GetReportMinerbaWithPeriod(inputCreateMinerbaLn.Period)

	// if findMinerbaErr == nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"error": "report with same period already exist",
	// 	})
	// }

	splitPeriod := strings.Split(inputCreateMinerbaLn.Period, " ")

	baseIdNumber := fmt.Sprintf("LML-%s-%s", helper.MonthStringToNumberString(splitPeriod[0]), splitPeriod[1])
	createMinerbaLn, createMinerbaLnErr := h.historyService.CreateMinerba(inputCreateMinerbaLn.Period, baseIdNumber, inputCreateMinerbaLn.ListDataLn, uint(claims["id"].(float64)))

	if createMinerbaLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["minerba_period"] = inputCreateMinerbaLn.Period
		inputMap["list_dn"] = inputCreateMinerbaLn.ListDataLn

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createMinerbaLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createMinerbaLnErr.Error(),
		})
	}

	return c.Status(201).JSON(createMinerbaLn)
}
