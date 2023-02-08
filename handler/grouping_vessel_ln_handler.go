package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
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

type groupingVesselLnHandler struct {
	transactionService      transaction.Service
	userService             user.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	groupingVesselLnService groupingvesselln.Service
}

func NewGroupingVesselLnHandler(transactionService transaction.Service, userService user.Service, historyService history.Service, v *validator.Validate, logService logs.Service, groupingVesselLnService groupingvesselln.Service) *groupingVesselLnHandler {
	return &groupingVesselLnHandler{
		transactionService,
		userService,
		historyService,
		v,
		logService,
		groupingVesselLnService,
	}
}

func (h *groupingVesselLnHandler) CreateGroupingVesselLn(c *fiber.Ctx) error {
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

	groupingVesselLnInput := new(groupingvesselln.InputGroupingVesselLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(groupingVesselLnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*groupingVesselLnInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdGroupingVesselLn, createdGroupingVesselLnErr := h.historyService.CreateGroupingVesselLN(*groupingVesselLnInput, uint(claims["id"].(float64)))

	if createdGroupingVesselLnErr != nil {
		inputJson, _ := json.Marshal(groupingVesselLnInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdGroupingVesselLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": createdGroupingVesselLnErr.Error(),
		})
	}

	return c.Status(201).JSON(createdGroupingVesselLn)
}

func (h *groupingVesselLnHandler) GetDetailGroupingVesselLn(c *fiber.Ctx) error {
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

	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}
	detailGroupingVesselLn, detailGroupingVesselLnErr := h.transactionService.GetDetailGroupingVesselLn(idInt)

	if detailGroupingVesselLnErr != nil {

		status := 400
		if detailGroupingVesselLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailGroupingVesselLnErr.Error(),
		})
	}

	return c.Status(200).JSON(detailGroupingVesselLn)
}

func (h *groupingVesselLnHandler) EditGroupingVesselLn(c *fiber.Ctx) error {
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

	editGroupingVesselLnInput := new(groupingvesselln.InputEditGroupingVesselLn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(editGroupingVesselLnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*editGroupingVesselLnInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	id := c.Params("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	editGroupingVesselLn, editGroupingVesselLnErr := h.historyService.EditGroupingVesselLn(idInt, *editGroupingVesselLnInput, uint(claims["id"].(float64)))

	if editGroupingVesselLnErr != nil {
		inputJson, _ := json.Marshal(editGroupingVesselLnInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": editGroupingVesselLnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": editGroupingVesselLnErr.Error(),
		})
	}

	return c.Status(200).JSON(editGroupingVesselLn)
}

func (h *groupingVesselLnHandler) UploadDocumentGroupingVesselLn(c *fiber.Ctx) error {
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

	id := c.Params("id")

	documentType := c.Params("type")
	responseErr := fiber.Map{
		"message": "failed to upload document",
	}

	file, errFormFile := c.FormFile("document")

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	switch documentType {
	case "peb",
		"insurance",
		"ls_export",
		"navy",
		"ska_coo",
		"coa_cow",
		"bl_mv":
	default:
		return c.Status(400).JSON(fiber.Map{
			"error":   "document type not found",
			"message": "failed to upload document",
		})
	}

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	detailGroupingVesselLn, detailGroupingVesselLnErr := h.transactionService.GetDetailGroupingVesselLn(idInt)

	if detailGroupingVesselLnErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to upload document",
			"error":   detailGroupingVesselLnErr.Error(),
		})
	}

	fileName := fmt.Sprintf("GML/%v/%v_%s.pdf", *detailGroupingVesselLn.Detail.IdNumber, *detailGroupingVesselLn.Detail.IdNumber, documentType)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["grouping_vessel_ln_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		groupingVesselLnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselLnId: &groupingVesselLnId,
			Input:              inputJson,
			Message:            messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	editDocument, editDocumentErr := h.historyService.UploadDocumentGroupingVesselLn(detailGroupingVesselLn.Detail.ID, up.Location, uint(claims["id"].(float64)), documentType)

	if editDocumentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           editDocumentErr.Error(),
			"upload_response": up,
		})

		groupingVesselLnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselLnId: &groupingVesselLnId,
			Input:              inputJson,
			Message:            messageJson,
		}
		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = editDocumentErr.Error()

		status := 400

		if editDocumentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(responseErr)
	}

	return c.Status(200).JSON(editDocument)
}

func (h *groupingVesselLnHandler) DeleteGroupingVesselLn(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete minerba",
			"error":   "record not found",
		})
	}

	detailGroupingVesselLn, detailGroupingVesselLnErr := h.transactionService.GetDetailGroupingVesselLn(idInt)

	if detailGroupingVesselLnErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete grouping vessel ln",
			"error":   "record not found",
		})
	}

	_, deleteGroupingVesselLnErr := h.historyService.DeleteGroupingVesselLn(idInt, uint(claims["id"].(float64)))

	if deleteGroupingVesselLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["grouping_vessel_ln_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteGroupingVesselLnErr.Error(),
		})

		groupingVesselLnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselLnId: &groupingVesselLnId,
			Input:              inputJson,
			Message:            messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteGroupingVesselLnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete grouping vessel ln",
			"error":   deleteGroupingVesselLnErr.Error(),
		})
	}

	if detailGroupingVesselLn.Detail.PebDocumentLink != "" ||
		detailGroupingVesselLn.Detail.InsuranceDocumentLink != "" ||
		detailGroupingVesselLn.Detail.LsExportDocumentLink != "" ||
		detailGroupingVesselLn.Detail.NavyDocumentLink != "" ||
		detailGroupingVesselLn.Detail.SkaCooDocumentLink != "" ||
		detailGroupingVesselLn.Detail.CoaCowDocumentLink != "" ||
		detailGroupingVesselLn.Detail.BlMvDocumentLink != "" {
		fileName := fmt.Sprintf("GML/%s/", *detailGroupingVesselLn.Detail.IdNumber)
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["grouping_vessel_ln_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			groupingVesselLnId := uint(idInt)
			createdErrLog := logs.Logs{
				GroupingVesselLnId: &groupingVesselLnId,
				Input:              inputJson,
				Message:            messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete grouping vessel ln aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete grouping vessel ln",
	})
}

func (h *groupingVesselLnHandler) ListGroupingVesselLn(c *fiber.Ctx) error {
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

	var sortAndFilter groupingvesselln.SortFilterGroupingVesselLn
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

	listGroupingVesselLn, listGroupingVesselLnErr := h.groupingVesselLnService.ListGroupingVesselLn(pageNumber, sortAndFilter)

	if listGroupingVesselLnErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listGroupingVesselLnErr.Error(),
		})
	}

	return c.Status(200).JSON(listGroupingVesselLn)
}

func (h *groupingVesselLnHandler) ListLnWithoutGroup(c *fiber.Ctx) error {
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

	listLnWithoutGroup, listLnWithoutGroupErr := h.transactionService.ListDataLnWithoutGroup()

	if listLnWithoutGroupErr != nil {

		status := 400

		if listLnWithoutGroupErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": listLnWithoutGroupErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"transactions": listLnWithoutGroup,
	})
}
