package handler

import (
	"mrpbackend/model/alatberat"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type alatBeratHandler struct {
	userService      user.Service
	alatBeratService alatberat.Service
	v                *validator.Validate
}

func NewAlatBeratHandler(userService user.Service, alatBeratService alatberat.Service, v *validator.Validate) *alatBeratHandler {
	return &alatBeratHandler{
		userService,
		alatBeratService,
		v,
	}
}

// Master Data

func (h *alatBeratHandler) CreateAlatBerat(c *fiber.Ctx) error {
	alatBeratInput := new(alatberat.RegisterAlatBeratInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(alatBeratInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*alatBeratInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdAlatBerat, createdAlatBeratErr := h.alatBeratService.CreateAlatBerat(*alatBeratInput)

	if createdAlatBeratErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdAlatBeratErr.Error(),
		})
	}

	return c.Status(201).JSON(createdAlatBerat)
}

func (h *alatBeratHandler) GetAlatBerat(c *fiber.Ctx) error {
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

	brands, brandsErr := h.alatBeratService.FindAlatBerat()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *alatBeratHandler) GetAlatBeratById(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.alatBeratService.FindAlatBeratById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *alatBeratHandler) GetListAlatBerat(c *fiber.Ctx) error {
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

	var filterAlatBerat alatberat.SortFilterAlatBerat

	filterAlatBerat.BrandId = c.Query("brand_id")
	filterAlatBerat.HeavyEquipmentId = c.Query("heavy_equipment_id")
	filterAlatBerat.SeriesId = c.Query("series_id")
	filterAlatBerat.Consumption = c.Query("consumption")
	filterAlatBerat.Tolerance = c.Query("tolerance")
	filterAlatBerat.Field = c.Query("field")
	filterAlatBerat.Sort = c.Query("sort")

	listAlatBerat, listAlatBeratErr := h.alatBeratService.GetListAlatBerat(pageNumber, filterAlatBerat)

	if listAlatBeratErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listAlatBeratErr.Error(),
		})
	}

	return c.Status(200).JSON(listAlatBerat)
}

func (h *alatBeratHandler) GetConsumption(c *fiber.Ctx) error {
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

	brandId := c.Params("brandId")

	brandIdInt, err := strconv.Atoi(brandId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	heavyEquipmentId := c.Params("heavyEquipmentId")

	heavyEquipmentIdInt, err := strconv.Atoi(heavyEquipmentId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	seriesId := c.Params("seriesId")

	seriesIdInt, err := strconv.Atoi(seriesId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	heavyEquipments, heavyEquipmentsErr := h.alatBeratService.FindConsumption(uint(brandIdInt), uint(heavyEquipmentIdInt), uint(seriesIdInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *alatBeratHandler) UpdateAlatBerat(c *fiber.Ctx) error {
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
	inputUpdateAlatBerat := new(alatberat.RegisterAlatBeratInput)
	if err := c.BodyParser(inputUpdateAlatBerat); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateAlatBerat)
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

	updateAlatBerat, err := h.alatBeratService.UpdateAlatBerat(*inputUpdateAlatBerat, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateAlatBerat)
}

func (h *alatBeratHandler) DeleteAlatBerat(c *fiber.Ctx) error {
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

	alatBerat, err := h.alatBeratService.FindAlatBeratById(uint(id))
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

	// Optional: Use alatBerat.ID for extra safety
	if _, err := h.alatBeratService.DeleteAlatBerat(alatBerat.ID); err != nil {
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
