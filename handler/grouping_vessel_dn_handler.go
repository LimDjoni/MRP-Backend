package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/destination"
	"ajebackend/model/transaction"
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

type groupingVesselDnHandler struct {
	transactionService      transaction.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	groupingVesselDnService groupingvesseldn.Service
	destinationService      destination.Service
	userIupopkService       useriupopk.Service
}

func NewGroupingVesselDnHandler(transactionService transaction.Service, historyService history.Service, v *validator.Validate, logService logs.Service, groupingVesselDnService groupingvesseldn.Service, destinationService destination.Service, userIupopkService useriupopk.Service) *groupingVesselDnHandler {
	return &groupingVesselDnHandler{
		transactionService,
		historyService,
		v,
		logService,
		groupingVesselDnService,
		destinationService,
		userIupopkService,
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

	var sortAndFilter groupingvesseldn.SortFilterGroupingVesselDn
	page := c.Query("page")
	sortAndFilter.Field = c.Query("field")
	sortAndFilter.Sort = c.Query("sort")
	sortAndFilter.BlDateStart = c.Query("bl_date_start")
	sortAndFilter.BlDateEnd = c.Query("bl_date_end")
	sortAndFilter.Quantity = c.Query("quantity")
	sortAndFilter.VesselId = c.Query("vessel_id")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if page == "" {
		pageNumber = 1
	}

	listGroupingVesselDn, listGroupingVesselDnErr := h.groupingVesselDnService.ListGroupingVesselDn(pageNumber, sortAndFilter, iupopkIdInt)

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

	createdGroupingVesselDn, createdGroupingVesselDnErr := h.historyService.CreateGroupingVesselDN(*groupingVesselDnInput, uint(claims["id"].(float64)), iupopkIdInt)

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

func (h *groupingVesselDnHandler) EditGroupingVesselDn(c *fiber.Ctx) error {
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

	editGroupingVesselDnInput := new(groupingvesseldn.InputEditGroupingVesselDn)

	// Binds the request body to the Person struct
	if err := c.BodyParser(editGroupingVesselDnInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*editGroupingVesselDnInput)

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

	editGroupingVesselDn, editGroupingVesselDnErr := h.historyService.EditGroupingVesselDn(idInt, *editGroupingVesselDnInput, uint(claims["id"].(float64)), iupopkIdInt)

	if editGroupingVesselDnErr != nil {
		inputJson, _ := json.Marshal(editGroupingVesselDnInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": editGroupingVesselDnErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": editGroupingVesselDnErr.Error(),
		})
	}

	return c.Status(200).JSON(editGroupingVesselDn)
}

func (h *groupingVesselDnHandler) UploadDocumentGroupingVesselDn(c *fiber.Ctx) error {
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
	case "coa_cow",
		"bl_mv", "skab":
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

	detailGroupingVesselDn, detailGroupingVesselDnErr := h.transactionService.GetDetailGroupingVesselDn(idInt, iupopkIdInt)

	if detailGroupingVesselDnErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to upload document",
			"error":   detailGroupingVesselDnErr.Error(),
		})
	}

	var fileName string

	fileName = detailGroupingVesselDn.Detail.Iupopk.Code

	fileName += fmt.Sprintf("/GMD/%v/%v_%s.pdf", *detailGroupingVesselDn.Detail.IdNumber, *detailGroupingVesselDn.Detail.IdNumber, documentType)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["grouping_vessel_dn_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		groupingVesselDnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselDnId: &groupingVesselDnId,
			Input:              inputJson,
			Message:            messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	editDocument, editDocumentErr := h.historyService.UploadDocumentGroupingVesselDn(detailGroupingVesselDn.Detail.ID, up.Location, uint(claims["id"].(float64)), documentType, iupopkIdInt)

	if editDocumentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["document_type"] = documentType
		inputMap["user_id"] = claims["id"]
		inputMap["grouping_vessel_dn_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           editDocumentErr.Error(),
			"upload_response": up,
		})

		groupingVesselDnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselDnId: &groupingVesselDnId,
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

func (h *groupingVesselDnHandler) GetDetailGroupingVesselDn(c *fiber.Ctx) error {
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
	detailGroupingVesselDn, detailGroupingVesselDnErr := h.transactionService.GetDetailGroupingVesselDn(idInt, iupopkIdInt)

	if detailGroupingVesselDnErr != nil {

		status := 400
		if detailGroupingVesselDnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailGroupingVesselDnErr.Error(),
		})
	}

	return c.Status(200).JSON(detailGroupingVesselDn)
}

func (h *groupingVesselDnHandler) DeleteGroupingVesselDn(c *fiber.Ctx) error {
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

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete grouping vessel dn",
			"error":   "record not found",
		})
	}

	detailGroupingVesselDn, detailGroupingVesselDnErr := h.transactionService.GetDetailGroupingVesselDn(idInt, iupopkIdInt)

	if detailGroupingVesselDnErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to delete grouping vessel dn",
			"error":   "record not found",
		})
	}

	_, deleteGroupingVesselDnErr := h.historyService.DeleteGroupingVesselDn(idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if deleteGroupingVesselDnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["grouping_vessel_dn_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": deleteGroupingVesselDnErr.Error(),
		})

		groupingVesselDnId := uint(idInt)
		createdErrLog := logs.Logs{
			GroupingVesselDnId: &groupingVesselDnId,
			Input:              inputJson,
			Message:            messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if deleteGroupingVesselDnErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete grouping vessel dn",
			"error":   deleteGroupingVesselDnErr.Error(),
		})
	}

	if detailGroupingVesselDn.Detail.CoaCowDocumentLink != nil ||
		detailGroupingVesselDn.Detail.BlMvDocumentLink != nil {
		var fileName string
		fileName = detailGroupingVesselDn.Detail.Iupopk.Code
		fileName += fmt.Sprintf("/GMD/%s/", *detailGroupingVesselDn.Detail.IdNumber)
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["grouping_vessel_dn_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			groupingVesselDnId := uint(idInt)
			createdErrLog := logs.Logs{
				GroupingVesselDnId: &groupingVesselDnId,
				Input:              inputJson,
				Message:            messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete grouping vessel dn aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete grouping vessel dn",
	})
}

func (h *groupingVesselDnHandler) ListDnWithoutGroup(c *fiber.Ctx) error {
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

	listDnWithoutGroup, listDnWithoutGroupErr := h.transactionService.ListDataDnWithoutGroup(iupopkIdInt)

	if listDnWithoutGroupErr != nil {

		status := 400

		if listDnWithoutGroupErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": listDnWithoutGroupErr.Error(),
		})
	}

	return c.Status(200).JSON(listDnWithoutGroup)
}
