package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/contract"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type contractHandler struct {
	contractService   contract.Service
	historyService    history.Service
	logService        logs.Service
	v                 *validator.Validate
	userIupopkService useriupopk.Service
	allMasterService  allmaster.Service
}

func NewContractHandler(contractService contract.Service, historyService history.Service, logService logs.Service, v *validator.Validate, userIupopkService useriupopk.Service, allMasterService allmaster.Service) *contractHandler {
	return &contractHandler{
		contractService,
		historyService,
		logService,
		v,
		userIupopkService,
		allMasterService,
	}
}

func (h *contractHandler) ListContract(c *fiber.Ctx) error {
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

	var sortAndFilter contract.FilterAndSortContract
	page := c.Query("page")
	sortAndFilter.Field = c.Query("field")
	sortAndFilter.Sort = c.Query("sort")
	sortAndFilter.Quantity = c.Query("quantity")
	sortAndFilter.ContractDateStart = c.Query("contract_date_start")
	sortAndFilter.ContractDateEnd = c.Query("contract_date_end")
	sortAndFilter.ContractNumber = c.Query("contract_number")
	sortAndFilter.CustomerId = c.Query("customer_id")
	sortAndFilter.ValidityStart = c.Query("validity_start")
	sortAndFilter.ValidityEnd = c.Query("validity_end")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if page == "" {
		pageNumber = 1
	}

	listContract, listContractErr := h.contractService.GetListReportContractAll(pageNumber, sortAndFilter, iupopkIdInt)

	if listContractErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": listContractErr.Error(),
		})
	}

	return c.Status(200).JSON(listContract)
}

func (h *contractHandler) CreateContract(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
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

	iupopkData, iupopkDataErr := h.allMasterService.FindIupopk(iupopkIdInt)

	if iupopkDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkDataErr.Error(),
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	contractInput := new(contract.InputCreateUpdateContract)

	// Binds the request body to the Person struct
	if err := c.BodyParser(contractInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*contractInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	var location string
	file, errFormFile := c.FormFile("file")

	if errFormFile != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to create contract file",
			"error":   errFormFile.Error(),
		})
	}

	if file != nil {
		fileName := fmt.Sprintf("%s/CONTRACT/%s/%s.pdf", iupopkData.Code, contractInput.ContractNumber, contractInput.ContractNumber)

		up, uploadErr := awshelper.UploadDocument(file, fileName)

		if uploadErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["file"] = file
			inputMap["user_id"] = claims["id"]
			inputMap["type"] = "CONTRACT"

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": uploadErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to create contract file",
				"error":   uploadErr.Error(),
			})
		}
		location = up.Location
	}

	createdContract, createdContractErr := h.historyService.CreateContract(*contractInput, location, iupopkIdInt, uint(claims["id"].(float64)))

	if createdContractErr != nil {
		inputJson, _ := json.Marshal(contractInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdContractErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createdContractErr.Error(),
		})
	}

	return c.Status(201).JSON(createdContract)
}

func (h *contractHandler) UpdateContract(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	id := c.Params("id")
	iupopkId := c.Params("iupopk_id")

	idInt, err := strconv.Atoi(id)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	iupopkData, iupopkDataErr := h.allMasterService.FindIupopk(iupopkIdInt)

	if iupopkDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkDataErr.Error(),
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	contractInput := new(contract.InputCreateUpdateContract)

	// Binds the request body to the Person struct
	if errParsing := c.BodyParser(contractInput); errParsing != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errParsing.Error(),
		})
	}

	errors := h.v.Struct(*contractInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	var location string
	file, _ := c.FormFile("file")

	// if errFormFile != nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"message": "failed to update contract file",
	// 		"error":   errFormFile.Error(),
	// 	})
	// }

	if file != nil {
		fileName := fmt.Sprintf("%s/CONTRACT/%s/%s.pdf", iupopkData.Code, contractInput.ContractNumber, contractInput.ContractNumber)

		up, uploadErr := awshelper.UploadDocument(file, fileName)

		if uploadErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["file"] = file
			inputMap["user_id"] = claims["id"]
			inputMap["contract_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": uploadErr.Error(),
			})

			contractId := uint(idInt)
			createdErrLog := logs.Logs{
				ContractId: &contractId,
				Input:      inputJson,
				Message:    messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to update contract file",
				"error":   uploadErr.Error(),
			})
		}
		location = up.Location
	}

	updateContract, updateContractErr := h.historyService.UpdateContract(idInt, *contractInput, location, iupopkIdInt, uint(claims["id"].(float64)))

	if updateContractErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["input"] = contractInput
		inputMap["user_id"] = claims["id"]
		inputMap["contract_id"] = id

		inputJson, _ := json.Marshal(contractInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateContractErr.Error(),
		})

		contractId := uint(idInt)
		createdErrLog := logs.Logs{
			ContractId: &contractId,
			Input:      inputJson,
			Message:    messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if updateContractErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to update contract",
			"error":   updateContractErr.Error(),
		})
	}

	return c.Status(200).JSON(updateContract)
}

func (h *contractHandler) DetailContract(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detailContract, detailContractErr := h.contractService.GetDataContract(idInt, iupopkIdInt)

	if detailContractErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailContractErr.Error(),
		})
	}

	return c.Status(200).JSON(detailContract)
}
