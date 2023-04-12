package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/notificationuser"
	"ajebackend/model/rkab"
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

type rkabHandler struct {
	rkabService             rkab.Service
	logService              logs.Service
	userIupopkService       useriupopk.Service
	historyService          history.Service
	notificationUserService notificationuser.Service
	v                       *validator.Validate
	allMasterService        allmaster.Service
}

func NewRkabHandler(
	rkabService rkab.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
	allMasterService allmaster.Service,
) *rkabHandler {
	return &rkabHandler{
		rkabService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
		allMasterService,
	}
}

func (h *rkabHandler) CreateRkab(c *fiber.Ctx) error {
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

	iupopkData, iupopkDataErr := h.allMasterService.FindIupopk(iupopkIdInt)

	if iupopkDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkDataErr.Error(),
		})
	}

	inputCreateRkab := new(rkab.RkabInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputCreateRkab)

	if errParsing != nil {
		fmt.Println(errParsing)
	}

	errors := h.v.Struct(*inputCreateRkab)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateRkab
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	file, errFormFile := c.FormFile("rkab_document")
	responseErr := fiber.Map{
		"message": "failed to create rkab",
	}

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	createRkab, createRkabErr := h.historyService.CreateRkab(*inputCreateRkab, iupopkIdInt, uint(claims["id"].(float64)))

	if createRkabErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateRkab
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createRkabErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createRkabErr.Error(),
		})
	}

	fileName := fmt.Sprintf("%s/RKB/%v/%v_rkab.pdf", iupopkData.Code, createRkab.IdNumber, createRkab.IdNumber)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateRkab

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			RkabId:  &createRkab.ID,
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["message"] = "failed to create rkab upload"
		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	updateRkab, updateRkabErr := h.historyService.UploadDocumentRkab(createRkab.ID, up.Location, uint(claims["id"].(float64)), iupopkIdInt)

	if updateRkabErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateRkab
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           updateRkabErr.Error(),
			"upload_response": up,
		})

		createdErrLog := logs.Logs{
			RkabId:  &createRkab.ID,
			Input:   inputJson,
			Message: messageJson,
		}
		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = updateRkabErr.Error()

		status := 400

		if updateRkabErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(responseErr)
	}

	return c.Status(201).JSON(updateRkab)
}

func (h *rkabHandler) ListRkab(c *fiber.Ctx) error {
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

	var filterRkab rkab.SortFilterRkab

	filterRkab.Field = c.Query("field")
	filterRkab.Sort = c.Query("sort")
	filterRkab.DateOfIssue = c.Query("date_of_issue")
	filterRkab.Year = c.Query("year")
	filterRkab.ProductionQuota = c.Query("production_quota")
	filterRkab.Status = c.Query("status")
	listRkab, listRkabErr := h.rkabService.ListRkab(pageNumber, filterRkab, iupopkIdInt)

	if listRkabErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listRkabErr.Error(),
		})
	}

	return c.Status(200).JSON(listRkab)
}

func (h *rkabHandler) DeleteRkab(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	iupopkId := c.Params("iupopk_id")
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailRkab, detailRkabErr := h.rkabService.DetailRkabWithId(idInt, iupopkIdInt)

	if detailRkabErr != nil {
		status := 400

		if detailRkabErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete rkab",
			"error":   detailRkabErr.Error(),
		})
	}

	_, isDeletedRkabErr := h.historyService.DeleteRkab(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedRkabErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["rkab_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedRkabErr.Error(),
		})

		createdErrLog := logs.Logs{
			RkabId:  &detailRkab.ID,
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedRkabErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete rkab",
			"error":   isDeletedRkabErr.Error(),
		})
	}

	if detailRkab.RkabDocumentLink != "" {
		documentLink := detailRkab.RkabDocumentLink

		documentLinkSplit := strings.Split(documentLink, "/")

		fileName := ""

		for i, v := range documentLinkSplit {
			if i == 3 {
				fileName += v + "/"
			}

			if i == 4 {
				fileName += v + "/"
			}

			if i == 5 {
				fileName += v
			}
		}
		_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

		if deleteAwsErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["user_id"] = claims["id"]
			inputMap["rkab_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			createdErrLog := logs.Logs{
				RkabId:  &detailRkab.ID,
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete rkab aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete rkab",
	})
}

func (h *rkabHandler) DetailRkab(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := fiber.Map{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	year := c.Params("year")
	iupopkId := c.Params("iupopk_id")

	yearInt, err := strconv.Atoi(year)
	iupopkIdInt, iupopkErr := strconv.Atoi(iupopkId)

	if err != nil || iupopkErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailRkab, detailRkabErr := h.rkabService.DetailRkabWithYear(yearInt, iupopkIdInt)

	if detailRkabErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailRkabErr.Error(),
		})
	}

	return c.Status(200).JSON(detailRkab)
}
