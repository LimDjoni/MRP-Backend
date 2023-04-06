package handler

import (
	"ajebackend/model/dmo"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/rkab"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type rkabHandler struct {
	rkabService             rkab.Service
	logService              logs.Service
	userIupopkService       useriupopk.Service
	historyService          history.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
}

func NewRkabHandler(
	rkabService rkab.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
) *rkabHandler {
	return &rkabHandler{
		rkabService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
	}
}

func (h *rkabHandler) CreateRkab(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	iupopkId := c.Params("iupopk_id")

	iupopkIdInt, err := strconv.Atoi(iupopkId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	inputCreateRkab := new(rkab.RkabInput)

	inputCreateDmo := new(dmo.CreateDmoInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputCreateRkab)

	if errParsing != nil {
		inputCreateDmo.Period = strings.Replace(inputCreateDmo.Period, "\"", "", -1)

		fmt.Println(errParsing)
	}

	errors := h.v.Struct(*inputCreateDmo)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateRkab
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

	file, errFormFile := c.FormFile("document")
	responseErr := fiber.Map{
		"message": "failed to upload document",
	}

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	// fileName := fmt.Sprintf()

	createRkab, createRkabErr := h.historyService.CreateRkab(*&inputCreateRkab, iupopkIdInt, uint(claims["id"].(float64)))

	if createRkabErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateDmo
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createRkabErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createRkabErr.Error(),
		})
	}

	return c.Status(201).JSON(createRkab)
}
