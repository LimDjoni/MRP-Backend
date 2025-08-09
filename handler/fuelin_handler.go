package handler

import (
	"mrpbackend/model/fuelin"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type fuelInHandler struct {
	userService   user.Service
	fuelInService fuelin.Service
	v             *validator.Validate
}

func NewFuelInHandler(userService user.Service, fuelInService fuelin.Service, v *validator.Validate) *fuelInHandler {
	return &fuelInHandler{
		userService,
		fuelInService,
		v,
	}
}

// Master Data

func (h *fuelInHandler) CreateFuelIn(c *fiber.Ctx) error {
	fuelInInput := new(fuelin.RegisterFuelInInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(fuelInInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*fuelInInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdFuelIn, createdFuelInErr := h.fuelInService.CreateFuelIn(*fuelInInput)

	if createdFuelInErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdFuelInErr.Error(),
		})
	}

	return c.Status(201).JSON(createdFuelIn)
}

func (h *fuelInHandler) GetFuelIn(c *fiber.Ctx) error {
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

	brands, brandsErr := h.fuelInService.FindFuelIn()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *fuelInHandler) GetFuelInById(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.fuelInService.FindFuelInById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *fuelInHandler) GetListFuelIn(c *fiber.Ctx) error {
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

	var filterFuelIn fuelin.SortFilterFuelIn

	filterFuelIn.Vendor = c.Query("vendor")
	filterFuelIn.Code = c.Query("code")
	filterFuelIn.NomorSuratJalan = c.Query("nomor_surat_jalan")
	filterFuelIn.NomorPlatMobil = c.Query("nomor_plat_mobil")
	filterFuelIn.Qty = c.Query("qty")
	filterFuelIn.QtyNow = c.Query("qty_now")
	filterFuelIn.Driver = c.Query("driver")
	filterFuelIn.TujuanAwal = c.Query("tujuan_awal")
	filterFuelIn.Field = c.Query("field")
	filterFuelIn.Sort = c.Query("sort")

	listFuelIn, listFuelInErr := h.fuelInService.GetListFuelIn(pageNumber, filterFuelIn)

	if listFuelInErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listFuelInErr.Error(),
		})
	}

	return c.Status(200).JSON(listFuelIn)
}

func (h *fuelInHandler) UpdateFuelIn(c *fiber.Ctx) error {
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
	inputUpdateFuelIn := new(fuelin.RegisterFuelInInput)
	if err := c.BodyParser(inputUpdateFuelIn); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateFuelIn)
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

	updateFuelIn, err := h.fuelInService.UpdateFuelIn(*inputUpdateFuelIn, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateFuelIn)
}

func (h *fuelInHandler) DeleteFuelIn(c *fiber.Ctx) error {
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

	fuelIn, err := h.fuelInService.FindFuelInById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find alat berat",
			"error":   err.Error(),
		})
	}

	// Optional: Use fuelIn.ID for extra safety
	if _, err := h.fuelInService.DeleteFuelIn(fuelIn.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete alat berat",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete alat berat",
	})
}
