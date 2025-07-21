package handler

import (
	"mrpbackend/model/fuelratio"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type fuelratiosHandler struct {
	userService       user.Service
	fuelratiosService fuelratio.Service
	v                 *validator.Validate
}

func NewFuelRatioHandler(userService user.Service, fuelratiosService fuelratio.Service, v *validator.Validate) *fuelratiosHandler {
	return &fuelratiosHandler{
		userService,
		fuelratiosService,
		v,
	}
}

// Master Data

func (h *fuelratiosHandler) CreateFuelRatio(c *fiber.Ctx) error {
	fuelratiosInput := new(fuelratio.RegisterFuelRatioInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(fuelratiosInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*fuelratiosInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdFuelRatio, createdFuelRatioErr := h.fuelratiosService.CreateFuelRatio(*fuelratiosInput)

	if createdFuelRatioErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdFuelRatioErr.Error(),
		})
	}

	return c.Status(201).JSON(createdFuelRatio)
}

func (h *fuelratiosHandler) GetFuelRatio(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	brands, brandsErr := h.fuelratiosService.FindFuelRatio()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *fuelratiosHandler) GetFuelRatioById(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	heavyEquipments, heavyEquipmentsErr := h.fuelratiosService.FindFuelRatioById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *fuelratiosHandler) GetListFuelRatio(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page == "" {
		pageNumber = 1
	}

	var filterFuelRatio fuelratio.SortFilterFuelRatio

	filterFuelRatio.UnitId = c.Query("unit_id")
	filterFuelRatio.EmployeeId = c.Query("employee_id")
	filterFuelRatio.Shift = c.Query("shift")
	filterFuelRatio.FirstHM = c.Query("first_hm")
	filterFuelRatio.Status = c.Query("status")
	filterFuelRatio.Field = c.Query("field")
	filterFuelRatio.Sort = c.Query("sort")

	listFuelRatio, listFuelRatioErr := h.fuelratiosService.GetListFuelRatio(pageNumber, filterFuelRatio)

	if listFuelRatioErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listFuelRatioErr.Error(),
		})
	}

	return c.Status(200).JSON(listFuelRatio)
}

func (h *fuelratiosHandler) GetFindFuelRatioExport(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	var filterFuelRatio fuelratio.SortFilterFuelRatioSummary

	filterFuelRatio.UnitName = c.Query("unit_name")
	filterFuelRatio.Shift = c.Query("shift")
	filterFuelRatio.TotalRefill = c.Query("total_refill")
	filterFuelRatio.Consumption = c.Query("consumption")
	filterFuelRatio.Tolerance = c.Query("tolerance")
	filterFuelRatio.FirstHM = c.Query("first_hm")
	filterFuelRatio.LastHM = c.Query("last_hm")
	filterFuelRatio.Field = c.Query("field")
	filterFuelRatio.Sort = c.Query("sort")

	listFuelRatio, listFuelRatioErr := h.fuelratiosService.FindFuelRatioExport(filterFuelRatio)

	if listFuelRatioErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listFuelRatioErr.Error(),
		})
	}

	return c.Status(200).JSON(listFuelRatio)
}

func (h *fuelratiosHandler) GetListRangkuman(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page == "" {
		pageNumber = 1
	}

	var filterFuelRatio fuelratio.SortFilterFuelRatioSummary

	filterFuelRatio.UnitName = c.Query("unit_name")
	filterFuelRatio.Shift = c.Query("shift")
	filterFuelRatio.TotalRefill = c.Query("total_refill")
	filterFuelRatio.Consumption = c.Query("consumption")
	filterFuelRatio.Tolerance = c.Query("tolerance")
	filterFuelRatio.FirstHM = c.Query("first_hm")
	filterFuelRatio.LastHM = c.Query("last_hm")
	filterFuelRatio.Field = c.Query("field")
	filterFuelRatio.Sort = c.Query("sort")

	listFuelRatio, listFuelRatioErr := h.fuelratiosService.ListRangkuman(pageNumber, filterFuelRatio)

	if listFuelRatioErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listFuelRatioErr.Error(),
		})
	}

	return c.Status(200).JSON(listFuelRatio)
}

func (h *fuelratiosHandler) UpdateFuelRatio(c *fiber.Ctx) error {
	//Get User
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}
	// Check User Login
	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	//Get Input
	inputUpdateFuelRatio := new(fuelratio.RegisterFuelRatioInput)
	if err := c.BodyParser(inputUpdateFuelRatio); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateFuelRatio)
	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	//Get ID
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	updateFuelRatio, err := h.fuelratiosService.UpdateFuelRatio(*inputUpdateFuelRatio, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateFuelRatio)
}

func (h *fuelratiosHandler) DeleteFuelRatio(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID parameter",
			"error":   err.Error(),
		})
	}

	fuelratios, err := h.fuelratiosService.FindFuelRatioById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find fuelratio",
			"error":   err.Error(),
		})
	}

	// Optional: Use fuelratios.ID for extra safety
	if _, err := h.fuelratiosService.DeleteFuelRatio(fuelratios.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete fuelratio",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete fuelratio",
	})
}
