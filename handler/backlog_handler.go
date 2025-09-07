package handler

import (
	"mrpbackend/model/backlog"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type backlogsHandler struct {
	userService     user.Service
	backlogsService backlog.Service
	v               *validator.Validate
}

func NewBackLogHandler(userService user.Service, backlogsService backlog.Service, v *validator.Validate) *backlogsHandler {
	return &backlogsHandler{
		userService,
		backlogsService,
		v,
	}
}

// Master Data

func (h *backlogsHandler) CreateBackLog(c *fiber.Ctx) error {
	backlogsInput := new(backlog.RegisterBackLogInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(backlogsInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*backlogsInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdBackLog, createdBackLogErr := h.backlogsService.CreateBackLog(*backlogsInput)

	if createdBackLogErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdBackLogErr.Error(),
		})
	}

	return c.Status(201).JSON(createdBackLog)
}

func (h *backlogsHandler) GetBackLog(c *fiber.Ctx) error {
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

	brands, brandsErr := h.backlogsService.FindBackLog()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *backlogsHandler) GetBackLogById(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.backlogsService.FindBackLogById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *backlogsHandler) GetListBackLog(c *fiber.Ctx) error {
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

	var filterBackLog backlog.SortFilterBackLog

	filterBackLog.BrandName = c.Query("brand_name")
	filterBackLog.UnitId = c.Query("unit_id")
	filterBackLog.HMBreakdown = c.Query("hm_breakdown")
	filterBackLog.Problem = c.Query("problem")
	filterBackLog.Component = c.Query("component")
	filterBackLog.DateOfInspection = c.Query("date_of_inspection")
	filterBackLog.PlanReplaceRepair = c.Query("plan_replace_repair")
	filterBackLog.HMReady = c.Query("hm_ready")
	filterBackLog.PPNumber = c.Query("pp_number")
	filterBackLog.PONumber = c.Query("po_number")
	filterBackLog.Status = c.Query("status")
	filterBackLog.AgingBacklogByDate = c.Query("aging_backlog_by_date")
	filterBackLog.Field = c.Query("field")
	filterBackLog.Sort = c.Query("sort")

	listBackLog, listBackLogErr := h.backlogsService.GetListBackLog(pageNumber, filterBackLog)

	if listBackLogErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listBackLogErr.Error(),
		})
	}

	return c.Status(200).JSON(listBackLog)
}

func (h *backlogsHandler) UpdateBackLog(c *fiber.Ctx) error {
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
	inputUpdateBackLog := new(backlog.RegisterBackLogInput)
	if err := c.BodyParser(inputUpdateBackLog); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateBackLog)
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

	updateBackLog, err := h.backlogsService.UpdateBackLog(*inputUpdateBackLog, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateBackLog)
}

func (h *backlogsHandler) DeleteBackLog(c *fiber.Ctx) error {
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

	backlogs, err := h.backlogsService.FindBackLogById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find backlog",
			"error":   err.Error(),
		})
	}

	// Optional: Use backlogs.ID for extra safety
	if _, err := h.backlogsService.DeleteBackLog(backlogs.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete backlog",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete backlog",
	})
}

func (h *backlogsHandler) GetListDashboardBackLog(c *fiber.Ctx) error {
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
	var filterBackLog backlog.SortFilterDashboardBacklog

	filterBackLog.Year = c.Query("year")

	brands, brandsErr := h.backlogsService.ListDashboardBackLog(filterBackLog)

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}
