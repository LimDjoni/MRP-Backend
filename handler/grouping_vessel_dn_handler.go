package handler

import (
	"ajebackend/model/destination"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type groupingVesselDnHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	groupingVesselDnService groupingvesseldn.Service
	destinationService      destination.Service
}

func NewGroupingVesselDnHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, v *validator.Validate, logService logs.Service, groupingVesselDnService groupingvesseldn.Service, destinationService destination.Service) *groupingVesselDnHandler {
	return &groupingVesselDnHandler{
		transactionService,
		userService,
		historyService,
		v,
		logService,
		groupingVesselDnService,
		destinationService,
	}
}

func (h *groupingVesselDnHandler) ListGroupingVesselDn(c *fiber.Ctx) error {
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

	var sortAndFilter groupingvesseldn.SortFilterGroupingVesselDn
	page := c.Query("page")
	sortAndFilter.Field = c.Query("field")
	sortAndFilter.Sort = c.Query("sort")

	quantity, errParsing := strconv.ParseFloat(c.Query("quantity"), 64)
	if errParsing != nil {
		sortAndFilter.Quantity = 0
	} else {
		sortAndFilter.Quantity = quantity
	}

	sortAndFilter.VesselName = c.Query("vessel_name")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if page == "" {
		pageNumber = 1
	}

	listGroupingVesselDn, listGroupingVesselDnErr := h.groupingVesselDnService.ListGroupingVesselDn(pageNumber, sortAndFilter)

	if listGroupingVesselDnErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listGroupingVesselDnErr.Error(),
		})
	}

	return c.Status(200).JSON(listGroupingVesselDn)
}

func (h *groupingVesselDnHandler) CreateGroupingVesselDn(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	checkUser, checkUserErr := h.userService.FindUser(uint(claims["id"].(float64)))

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	groupingVesselDnInput := new(groupingvesseldn.InputGroupingVesselDn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(groupingVesselDnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*groupingVesselDnInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	findDestination, findDestinationErr := h.destinationService.GetDestinationByName(*&groupingVesselDnInput.Destination)

	if findDestinationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "destination not found",
		})
	}

	groupingVesselDnInput.DestinationId = findDestination.ID

	createdGroupingVesselDn, createdGroupingVesselDnErr := h.historyService.CreateGroupingVesselDN(*groupingVesselDnInput, uint(claims["id"].(float64)))

	if createdGroupingVesselDnErr != nil {
		inputJson, _ := json.Marshal(groupingVesselDnInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdGroupingVesselDnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createdGroupingVesselDnErr.Error(),
		})
	}

	return c.Status(201).JSON(createdGroupingVesselDn)
}
