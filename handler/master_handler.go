package handler

import (
	"ajebackend/model/counter"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type masterHandler struct {
	destinationService destination.Service
	userService        user.Service
	allMasterService   allmaster.Service
	counterService     counter.Service
	v                  *validator.Validate
}

func NewMasterHandler(destinationService destination.Service, userService user.Service, allMasterService allmaster.Service, counterService counter.Service, v *validator.Validate) *masterHandler {
	return &masterHandler{
		destinationService,
		userService,
		allMasterService,
		counterService,
		v,
	}
}

func (h *masterHandler) GetDestination(c *fiber.Ctx) error {
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

	destinations, destinationsErr := h.destinationService.GetDestination()

	if destinationsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": destinationsErr.Error(),
		})
	}

	return c.Status(200).JSON(destinations)
}

func (h *masterHandler) GetListMaster(c *fiber.Ctx) error {

	listMaster, listMasterErr := h.allMasterService.ListMasterData()

	if listMasterErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listMasterErr.Error(),
		})
	}

	return c.Status(200).JSON(listMaster)
}

func (h *masterHandler) UpdateCounter(c *fiber.Ctx) error {
	updateCounterErr := h.counterService.UpdateCounter()

	if updateCounterErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": updateCounterErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "counter has been updated",
	})
}

// Master Data

func (h *masterHandler) CreateIupopk(c *fiber.Ctx) error {
	iupopkInput := new(iupopk.InputIupopk)

	// Binds the request body to the Person struct
	if err := c.BodyParser(iupopkInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*iupopkInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdIupopk, createdIupopkErr := h.counterService.CreateIupopk(*iupopkInput)

	if createdIupopkErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdIupopkErr.Error(),
		})
	}

	return c.Status(201).JSON(createdIupopk)
}

func (h *masterHandler) CreateBarge(c *fiber.Ctx) error {
	bargeInput := new(allmaster.InputBarge)

	// Binds the request body to the Person struct
	if err := c.BodyParser(bargeInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*bargeInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createBarge, createBargeErr := h.allMasterService.CreateBarge(*bargeInput)

	if createBargeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createBargeErr.Error(),
		})
	}

	return c.Status(201).JSON(createBarge)
}

func (h *masterHandler) CreateTugboat(c *fiber.Ctx) error {
	tugboatInput := new(allmaster.InputTugboat)

	// Binds the request body to the Person struct
	if err := c.BodyParser(tugboatInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*tugboatInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createTugboat, createTugboatErr := h.allMasterService.CreateTugboat(*tugboatInput)

	if createTugboatErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createTugboatErr.Error(),
		})
	}

	return c.Status(201).JSON(createTugboat)
}

func (h *masterHandler) CreateVessel(c *fiber.Ctx) error {
	vesselInput := new(allmaster.InputVessel)

	// Binds the request body to the Person struct
	if err := c.BodyParser(vesselInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*vesselInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createVessel, createVesselErr := h.allMasterService.CreateVessel(*vesselInput)

	if createVesselErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createVesselErr.Error(),
		})
	}

	return c.Status(201).JSON(createVessel)
}

func (h *masterHandler) CreatePortLocation(c *fiber.Ctx) error {
	portLocationInput := new(allmaster.InputPortLocation)

	// Binds the request body to the Person struct
	if err := c.BodyParser(portLocationInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*portLocationInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createPortLocation, createPortLocationErr := h.allMasterService.CreatePortLocation(*portLocationInput)

	if createPortLocationErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createPortLocationErr.Error(),
		})
	}

	return c.Status(201).JSON(createPortLocation)
}

func (h *masterHandler) CreatePort(c *fiber.Ctx) error {
	portInput := new(allmaster.InputPort)

	// Binds the request body to the Person struct
	if err := c.BodyParser(portInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*portInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createPort, createPortErr := h.allMasterService.CreatePort(*portInput)

	if createPortErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createPortErr.Error(),
		})
	}

	return c.Status(201).JSON(createPort)
}

func (h *masterHandler) CreateCompany(c *fiber.Ctx) error {
	companyInput := new(allmaster.InputCompany)

	// Binds the request body to the Person struct
	if err := c.BodyParser(companyInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*companyInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createCompany, createCompanyErr := h.allMasterService.CreateCompany(*companyInput)

	if createCompanyErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createCompanyErr.Error(),
		})
	}

	return c.Status(201).JSON(createCompany)
}

func (h *masterHandler) CreateTrader(c *fiber.Ctx) error {
	traderInput := new(allmaster.InputTrader)

	// Binds the request body to the Person struct
	if err := c.BodyParser(traderInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*traderInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createTrader, createTraderErr := h.allMasterService.CreateTrader(*traderInput)

	if createTraderErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createTraderErr.Error(),
		})
	}

	return c.Status(201).JSON(createTrader)
}

func (h *masterHandler) CreateIndustryType(c *fiber.Ctx) error {
	industryTypeInput := new(allmaster.InputIndustryType)

	// Binds the request body to the Person struct
	if err := c.BodyParser(industryTypeInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*industryTypeInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createIndustryType, createIndustryTypeErr := h.allMasterService.CreateIndustryType(*industryTypeInput)

	if createIndustryTypeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createIndustryTypeErr.Error(),
		})
	}

	return c.Status(201).JSON(createIndustryType)
}

func (h *masterHandler) UpdateBarge(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	bargeInput := new(allmaster.InputBarge)

	// Binds the request body to the Person struct
	if err := c.BodyParser(bargeInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*bargeInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateBarge, updateBargeErr := h.allMasterService.UpdateBarge(idInt, *bargeInput)

	if updateBargeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateBargeErr.Error(),
		})
	}

	return c.Status(200).JSON(updateBarge)
}

func (h *masterHandler) UpdateTugboat(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	tugboatInput := new(allmaster.InputTugboat)

	// Binds the request body to the Person struct
	if err := c.BodyParser(tugboatInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*tugboatInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateTugboat, updateTugboatErr := h.allMasterService.UpdateTugboat(idInt, *tugboatInput)

	if updateTugboatErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateTugboatErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTugboat)
}

func (h *masterHandler) UpdateVessel(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	vesselInput := new(allmaster.InputVessel)

	// Binds the request body to the Person struct
	if err := c.BodyParser(vesselInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*vesselInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateVessel, updateVesselErr := h.allMasterService.UpdateVessel(idInt, *vesselInput)

	if updateVesselErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateVesselErr.Error(),
		})
	}

	return c.Status(200).JSON(updateVessel)
}

func (h *masterHandler) UpdatePortLocation(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	portLocationInput := new(allmaster.InputPortLocation)

	// Binds the request body to the Person struct
	if err := c.BodyParser(portLocationInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*portLocationInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updatePortLocation, updatePortLocationErr := h.allMasterService.UpdatePortLocation(idInt, *portLocationInput)

	if updatePortLocationErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updatePortLocationErr.Error(),
		})
	}

	return c.Status(200).JSON(updatePortLocation)
}

func (h *masterHandler) UpdatePort(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	portInput := new(allmaster.InputPort)

	// Binds the request body to the Person struct
	if err := c.BodyParser(portInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*portInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updatePort, updatePortErr := h.allMasterService.UpdatePort(idInt, *portInput)

	if updatePortErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updatePortErr.Error(),
		})
	}

	return c.Status(200).JSON(updatePort)
}

func (h *masterHandler) UpdateCompany(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	companyInput := new(allmaster.InputCompany)

	// Binds the request body to the Person struct
	if err := c.BodyParser(companyInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*companyInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateCompany, updateCompanyErr := h.allMasterService.UpdateCompany(idInt, *companyInput)

	if updateCompanyErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateCompanyErr.Error(),
		})
	}

	return c.Status(200).JSON(updateCompany)
}

func (h *masterHandler) UpdateTrader(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	traderInput := new(allmaster.InputTrader)

	// Binds the request body to the Person struct
	if err := c.BodyParser(traderInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*traderInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateTrader, updateTraderErr := h.allMasterService.UpdateTrader(idInt, *traderInput)

	if updateTraderErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateTraderErr.Error(),
		})
	}

	return c.Status(200).JSON(updateTrader)
}

func (h *masterHandler) UpdateIndustryType(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	industryTypeInput := new(allmaster.InputIndustryType)

	// Binds the request body to the Person struct
	if err := c.BodyParser(industryTypeInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*industryTypeInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateIndustryType, updateIndustryTypeErr := h.allMasterService.UpdateIndustryType(idInt, *industryTypeInput)

	if updateIndustryTypeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": updateIndustryTypeErr.Error(),
		})
	}

	return c.Status(200).JSON(updateIndustryType)
}

func (h *masterHandler) DeleteBarge(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteBargeErr := h.allMasterService.DeleteBarge(idInt)

	if deleteBargeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteBargeErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeleteTugboat(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteTugboatErr := h.allMasterService.DeleteTugboat(idInt)

	if deleteTugboatErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteTugboatErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeleteVessel(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteVesselErr := h.allMasterService.DeleteVessel(idInt)

	if deleteVesselErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteVesselErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeletePortLocation(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deletePortLocationErr := h.allMasterService.DeletePortLocation(idInt)

	if deletePortLocationErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deletePortLocationErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeletePort(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deletePortErr := h.allMasterService.DeletePort(idInt)

	if deletePortErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deletePortErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeleteCompany(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteCompanyErr := h.allMasterService.DeleteCompany(idInt)

	if deleteCompanyErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteCompanyErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeleteTrader(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteTraderErr := h.allMasterService.DeleteTrader(idInt)

	if deleteTraderErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteTraderErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) DeleteIndustryType(c *fiber.Ctx) error {

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	_, deleteIndustryTypeErr := h.allMasterService.DeleteIndustryType(idInt)

	if deleteIndustryTypeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   deleteIndustryTypeErr.Error(),
			"message": "failed to delete",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete",
	})
}

func (h *masterHandler) ListTrader(c *fiber.Ctx) error {

	listTrader, listTraderErr := h.allMasterService.ListTrader()

	if listTraderErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listTraderErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"traders": listTrader,
	})
}

func (h *masterHandler) ListCompany(c *fiber.Ctx) error {

	listCompany, listCompanyErr := h.allMasterService.ListCompany()

	if listCompanyErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listCompanyErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"companies": listCompany,
	})
}
