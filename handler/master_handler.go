package handler

import (
	"mrpbackend/model/allmaster"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type masterHandler struct {
	userService      user.Service
	allMasterService allmaster.Service
	v                *validator.Validate
}

func NewMasterHandler(userService user.Service, allMasterService allmaster.Service, v *validator.Validate) *masterHandler {
	return &masterHandler{
		userService,
		allMasterService,
		v,
	}
}

// Master Data
func (h *masterHandler) CreateUserRole(c *fiber.Ctx) error {
	userRoleInput := new(allmaster.RegisterUserRoleInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(userRoleInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*userRoleInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdUserRole, createdUserRoleErr := h.allMasterService.CreateUserRole(*userRoleInput)

	if createdUserRoleErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdUserRoleErr.Error(),
		})
	}

	return c.Status(201).JSON(createdUserRole)
}

func (h *masterHandler) CreateBrand(c *fiber.Ctx) error {
	brandInput := new(allmaster.RegisterBrandInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(brandInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*brandInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdBrand, createdBrandErr := h.allMasterService.CreateBrand(*brandInput)

	if createdBrandErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdBrandErr.Error(),
		})
	}

	return c.Status(201).JSON(createdBrand)
}

func (h *masterHandler) CreateHeavyEquipment(c *fiber.Ctx) error {
	HeavyEquipmentInput := new(allmaster.RegisterHeavyEquipmentInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(HeavyEquipmentInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*HeavyEquipmentInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createCreateHeavyEquipmentInput, createCreateHeavyEquipmentInputErr := h.allMasterService.CreateHeavyEquipment(*HeavyEquipmentInput)

	if createCreateHeavyEquipmentInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createCreateHeavyEquipmentInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createCreateHeavyEquipmentInput)
}

func (h *masterHandler) CreateSeries(c *fiber.Ctx) error {
	SeriesInput := new(allmaster.RegisterSeriesInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(SeriesInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*SeriesInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createCreateSeriesInput, createCreateSeriesInputErr := h.allMasterService.CreateSeries(*SeriesInput)

	if createCreateSeriesInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createCreateSeriesInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createCreateSeriesInput)
}

func (h *masterHandler) CreateKartuKeluarga(c *fiber.Ctx) error {
	KartuKeluargaInput := new(allmaster.RegisterKartuKeluargaInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(KartuKeluargaInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*KartuKeluargaInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createKartuKeluargasInput, createKartuKeluargasInputErr := h.allMasterService.CreateKartuKeluarga(*KartuKeluargaInput)

	if createKartuKeluargasInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createKartuKeluargasInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createKartuKeluargasInput)
}

func (h *masterHandler) CreateKTP(c *fiber.Ctx) error {
	KTPInput := new(allmaster.RegisterKTPInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(KTPInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*KTPInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createKTPsInput, createKTPsInputErr := h.allMasterService.CreateKTP(*KTPInput)

	if createKTPsInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createKTPsInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createKTPsInput)
}

func (h *masterHandler) CreatePendidikan(c *fiber.Ctx) error {
	PendidikanInput := new(allmaster.RegisterPendidikanInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(PendidikanInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*PendidikanInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createPendidikansInput, createPendidikansInputErr := h.allMasterService.CreatePendidikan(*PendidikanInput)

	if createPendidikansInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createPendidikansInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createPendidikansInput)
}

func (h *masterHandler) CreateDOH(c *fiber.Ctx) error {
	DOHInput := new(allmaster.RegisterDOHInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(DOHInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*DOHInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createDOHsInput, createDOHsInputErr := h.allMasterService.CreateDOH(*DOHInput)

	if createDOHsInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createDOHsInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createDOHsInput)
}

func (h *masterHandler) CreateJabatan(c *fiber.Ctx) error {
	JabatanInput := new(allmaster.RegisterJabatanInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(JabatanInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*JabatanInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createJabatansInput, createJabatansInputErr := h.allMasterService.CreateJabatan(*JabatanInput)

	if createJabatansInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createJabatansInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createJabatansInput)
}

func (h *masterHandler) CreateSertifikat(c *fiber.Ctx) error {
	SertifikatInput := new(allmaster.RegisterSertifikatInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(SertifikatInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*SertifikatInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createSertifikatsInput, createSertifikatsInputErr := h.allMasterService.CreateSertifikat(*SertifikatInput)

	if createSertifikatsInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createSertifikatsInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createSertifikatsInput)
}

func (h *masterHandler) CreateMCU(c *fiber.Ctx) error {
	MCUInput := new(allmaster.RegisterMCUInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(MCUInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*MCUInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createMCUsInput, createMCUsInputErr := h.allMasterService.CreateMCU(*MCUInput)

	if createMCUsInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createMCUsInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createMCUsInput)
}

func (h *masterHandler) CreateHistory(c *fiber.Ctx) error {
	HistoryInput := new(allmaster.RegisterHistoryInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(HistoryInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*HistoryInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createHistorysInput, createHistorysInputErr := h.allMasterService.CreateHistory(*HistoryInput)

	if createHistorysInputErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createHistorysInputErr.Error(),
		})
	}

	return c.Status(201).JSON(createHistorysInput)
}

func (h *masterHandler) GetUserRole(c *fiber.Ctx) error {
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

	userRoles, userRolesErr := h.allMasterService.FindUserRole()

	if userRolesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": userRolesErr.Error(),
		})
	}

	return c.Status(200).JSON(userRoles)
}

func (h *masterHandler) GetUserRoleById(c *fiber.Ctx) error {
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

	userRoles, userRolesErr := h.allMasterService.FindUserRoleById(uint(idInt))

	if userRolesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": userRolesErr.Error(),
		})
	}

	return c.Status(200).JSON(userRoles)
}

func (h *masterHandler) GetBrand(c *fiber.Ctx) error {
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

	brands, brandsErr := h.allMasterService.FindBrand()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *masterHandler) GetBrandById(c *fiber.Ctx) error {
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

	brands, brandsErr := h.allMasterService.FindBrandById(uint(idInt))

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *masterHandler) GetHeavyEquipment(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.allMasterService.FindHeavyEquipment()

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *masterHandler) GetHeavyEquipmentById(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.allMasterService.FindHeavyEquipmentById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *masterHandler) GetHeavyEquipmentByBrandId(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.allMasterService.FindHeavyEquipmentByBrandID(uint(brandIdInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *masterHandler) GetSeries(c *fiber.Ctx) error {
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

	series, seriesErr := h.allMasterService.FindSeries()

	if seriesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": seriesErr.Error(),
		})
	}

	return c.Status(200).JSON(series)
}

func (h *masterHandler) GetSeriesById(c *fiber.Ctx) error {
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

	series, seriesErr := h.allMasterService.FindSeriesById(uint(idInt))

	if seriesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": seriesErr.Error(),
		})
	}

	return c.Status(200).JSON(series)
}

func (h *masterHandler) GetSeriesByBrandAndEquipmentdID(c *fiber.Ctx) error {
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

	heavyEquipmentId := c.Params("heavyequipmentId")

	heavyEquipmentIdInt, err := strconv.Atoi(heavyEquipmentId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	series, seriesErr := h.allMasterService.FindSeriesByBrandAndEquipmentdID(uint(brandIdInt), uint(heavyEquipmentIdInt))

	if seriesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": seriesErr.Error(),
		})
	}

	return c.Status(200).JSON(series)
}

func (h *masterHandler) GetDepartment(c *fiber.Ctx) error {
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

	brands, brandsErr := h.allMasterService.FindDepartment()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *masterHandler) GetRole(c *fiber.Ctx) error {
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

	brands, brandsErr := h.allMasterService.FindRole()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *masterHandler) GetPosition(c *fiber.Ctx) error {
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

	brands, brandsErr := h.allMasterService.FindPosition()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *masterHandler) UpdateDOH(c *fiber.Ctx) error {
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
	inputUpdateDOH := new(allmaster.RegisterDOHInput)
	if err := c.BodyParser(inputUpdateDOH); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateDOH)
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

	updateDOH, err := h.allMasterService.UpdateDOH(*inputUpdateDOH, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateDOH)
}

func (h *masterHandler) DeleteDOH(c *fiber.Ctx) error {
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

	employee, err := h.allMasterService.FindDOHById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.allMasterService.DeleteDOH(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete doh",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete doh",
	})
}

func (h *masterHandler) GetDohKontrak(c *fiber.Ctx) error {
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

	var filterDoh allmaster.SortFilterDohKontrak

	filterDoh.Field = c.Query("field")
	filterDoh.PT = c.Query("pt")
	filterDoh.Year = c.Query("year")
	filterDoh.Sort = c.Query("sort")
	filterDoh.CodeEmp = c.Query("code_emp")

	listDoh, listDohErr := h.allMasterService.FindDohKontrak(pageNumber, filterDoh)

	if listDohErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDohErr.Error(),
		})
	}

	return c.Status(200).JSON(listDoh)
}

func (h *masterHandler) UpdateJabatan(c *fiber.Ctx) error {
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
	inputUpdateJabatan := new(allmaster.RegisterJabatanInput)
	if err := c.BodyParser(inputUpdateJabatan); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateJabatan)
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

	updateJabatan, err := h.allMasterService.UpdateJabatan(*inputUpdateJabatan, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateJabatan)
}

func (h *masterHandler) DeleteJabatan(c *fiber.Ctx) error {
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

	employee, err := h.allMasterService.FindJabatanById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.allMasterService.DeleteJabatan(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete jabatan",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete jabatan",
	})
}

func (h *masterHandler) UpdateSertifikat(c *fiber.Ctx) error {
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
	inputUpdateSertifikat := new(allmaster.RegisterSertifikatInput)
	if err := c.BodyParser(inputUpdateSertifikat); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateSertifikat)
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

	updateSertifikat, err := h.allMasterService.UpdateSertifikat(*inputUpdateSertifikat, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateSertifikat)
}

func (h *masterHandler) DeleteSertifikat(c *fiber.Ctx) error {
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

	employee, err := h.allMasterService.FindSertifikatById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.allMasterService.DeleteSertifikat(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete sertifikat",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete sertifikat",
	})
}

func (h *masterHandler) UpdateMCU(c *fiber.Ctx) error {
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
	inputUpdateMCU := new(allmaster.RegisterMCUInput)
	if err := c.BodyParser(inputUpdateMCU); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateMCU)
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

	updateMCU, err := h.allMasterService.UpdateMCU(*inputUpdateMCU, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateMCU)
}

func (h *masterHandler) DeleteMCU(c *fiber.Ctx) error {
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

	employee, err := h.allMasterService.FindMCUById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.allMasterService.DeleteMCU(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete mcu",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete mcu",
	})
}

func (h *masterHandler) GetMCUBerkala(c *fiber.Ctx) error {
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

	var filterMCU allmaster.SortFilterDohKontrak

	filterMCU.PT = c.Query("pt")
	filterMCU.Year = c.Query("year")
	filterMCU.Field = c.Query("field")
	filterMCU.Sort = c.Query("sort")
	filterMCU.CodeEmp = c.Query("code_emp")

	listDoh, listDohErr := h.allMasterService.FindMCUBerkala(pageNumber, filterMCU)

	if listDohErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDohErr.Error(),
		})
	}

	return c.Status(200).JSON(listDoh)
}

func (h *masterHandler) UpdateHistory(c *fiber.Ctx) error {
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
	inputUpdateHistory := new(allmaster.RegisterHistoryInput)
	if err := c.BodyParser(inputUpdateHistory); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate Input
	errors := h.v.Struct(*inputUpdateHistory)
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

	updateHistory, err := h.allMasterService.UpdateHistory(*inputUpdateHistory, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateHistory)
}

func (h *masterHandler) DeleteHistory(c *fiber.Ctx) error {
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

	employee, err := h.allMasterService.FindHistoryById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.allMasterService.DeleteHistory(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete history",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete history",
	})
}

func (h *masterHandler) GenerateSideBar(c *fiber.Ctx) error {
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

	userId := c.Params("userId")

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	employees, employeesErr := h.allMasterService.GenerateSideBar(uint(userIdInt))

	if employeesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": employeesErr.Error(),
		})
	}

	return c.Status(200).JSON(employees)
}
