package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/cafassignment"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
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

type cafAssignmentHandler struct {
	cafAssignmentService        cafassignment.Service
	cafAssignmentEndUserService cafassignmentenduser.Service
	logService                  logs.Service
	userIupopkService           useriupopk.Service
	historyService              history.Service
	notificationUserService     notificationuser.Service
	v                           *validator.Validate
	allMasterService            allmaster.Service
}

func NewCafAssignmentHandler(
	cafAssignmentService cafassignment.Service,
	cafAssignmentEndUserService cafassignmentenduser.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
	allMasterService allmaster.Service,
) *cafAssignmentHandler {
	return &cafAssignmentHandler{
		cafAssignmentService,
		cafAssignmentEndUserService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
		allMasterService,
	}
}

func (h *cafAssignmentHandler) ListCafAssignment(c *fiber.Ctx) error {
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

	var filterCafAssignment cafassignment.SortFilterCafAssignment

	filterCafAssignment.Quantity = c.Query("quantity")
	filterCafAssignment.Year = c.Query("year")
	filterCafAssignment.Field = c.Query("field")
	filterCafAssignment.Sort = c.Query("sort")

	listCafAssignment, listCafAssignmentErr := h.cafAssignmentService.ListCafAssignment(pageNumber, filterCafAssignment, iupopkIdInt)

	if listCafAssignmentErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listCafAssignmentErr.Error(),
		})
	}

	return c.Status(200).JSON(listCafAssignment)
}

func (h *cafAssignmentHandler) CreateCafAssignment(c *fiber.Ctx) error {
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

	inputCreateCafAssignment := new(cafassignmentenduser.CreateCafAssignmentInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputCreateCafAssignment)

	if errParsing == nil {
		formPart, errFormPart := c.MultipartForm()
		if errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}
		if len(inputCreateCafAssignment.ListCafAssignment) == 0 {
			for _, value := range formPart.Value["list_caf_assignment"] {

				var tempCafAssignment []cafassignmentenduser.CafAssignmentInput

				errUnmarshal := json.Unmarshal([]byte(value), &tempCafAssignment)

				fmt.Println(errUnmarshal)
				inputCreateCafAssignment.ListCafAssignment = tempCafAssignment
			}

			if len(inputCreateCafAssignment.ListCafAssignment) == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "please check there is no caf assignment data",
				})
			}
		}
	}

	errors := h.v.Struct(*inputCreateCafAssignment)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCafAssignment
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

	createCafAssignment, createCafAssignmentErr := h.historyService.CreateCafAssignment(*inputCreateCafAssignment, uint(claims["id"].(float64)), iupopkIdInt)

	if createCafAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createCafAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createCafAssignmentErr.Error(),
		})
	}

	file, errFormFile := c.FormFile("letter_assignment")

	responseErr := fiber.Map{
		"message": "failed to create caf assignment",
	}

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	fileName := fmt.Sprintf("%s/SPS/%v/%v_letter_assignment.pdf", iupopkData.Code, createCafAssignment.IdNumber, createCafAssignment.IdNumber)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCafAssignment

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			CafAssignmentId: &createCafAssignment.ID,
			Input:           inputJson,
			Message:         messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["message"] = "failed to create caf assigment upload"
		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	updateCafAssignment, updateCafAssignmentErr := h.historyService.UploadCreateDocumentCafAssignment(createCafAssignment.ID, up.Location, uint(claims["id"].(float64)), iupopkIdInt)

	if updateCafAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           updateCafAssignmentErr.Error(),
			"upload_response": up,
		})

		createdErrLog := logs.Logs{
			CafAssignmentId: &createCafAssignment.ID,
			Input:           inputJson,
			Message:         messageJson,
		}
		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = updateCafAssignmentErr.Error()

		status := 400

		if updateCafAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(responseErr)
	}

	var createNotif notification.InputNotification

	createNotif.Type = "caf"
	createNotif.Status = "membuat"
	createNotif.Period = updateCafAssignment.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = updateCafAssignment.ID
		inputMap["caf_assignment"] = updateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationErr.Error(),
		})
	}

	return c.Status(201).JSON(updateCafAssignment)
}

func (h *cafAssignmentHandler) DetailCafAssignment(c *fiber.Ctx) error {
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

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	detailCafAssignment, detailCafAssignmentErr := h.cafAssignmentEndUserService.DetailCafAssignment(idInt, iupopkIdInt)

	if detailCafAssignmentErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailCafAssignmentErr.Error(),
		})
	}

	return c.Status(200).JSON(detailCafAssignment)
}

func (h *cafAssignmentHandler) UpdateCafAssignment(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
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

	inputUpdateCafAssignment := new(cafassignmentenduser.UpdateCafAssignmentInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputUpdateCafAssignment)

	if errParsing == nil {
		formPart, errFormPart := c.MultipartForm()
		if errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}
		if len(inputUpdateCafAssignment.ListCafAssignment) == 0 {
			for _, value := range formPart.Value["list_caf_assignment"] {

				var tempCafAssignment []cafassignmentenduser.CafAssignmentEndUser

				errUnmarshal := json.Unmarshal([]byte(value), &tempCafAssignment)

				fmt.Println(errUnmarshal)

				inputUpdateCafAssignment.ListCafAssignment = tempCafAssignment
			}

			if len(inputUpdateCafAssignment.ListCafAssignment) == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "please check there is no caf assignment data",
				})
			}

		}
	}

	errors := h.v.Struct(*inputUpdateCafAssignment)
	idUint := uint(idInt)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = idInt
		inputMap["input"] = inputUpdateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			CafAssignmentId: &idUint,
			Input:           inputJson,
			Message:         messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateCafAssignment, updateCafAssignmentErr := h.historyService.UpdateCafAssignment(idInt, *inputUpdateCafAssignment, uint(claims["id"].(float64)), iupopkIdInt)

	if updateCafAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = idInt
		inputMap["input"] = inputUpdateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateCafAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			CafAssignmentId: &idUint,
			Input:           inputJson,
			Message:         messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateCafAssignmentErr.Error(),
		})
	}

	file, errFormFile := c.FormFile("revision_letter_assignment")

	if errFormFile == nil {

		responseErr := fiber.Map{
			"message": "failed to create caf assignment",
		}

		if errFormFile != nil {
			responseErr["error"] = errFormFile.Error()
			return c.Status(400).JSON(responseErr)
		}

		if !strings.Contains(file.Filename, ".pdf") {
			responseErr["error"] = "document must be pdf"
			return c.Status(400).JSON(responseErr)
		}

		fileName := fmt.Sprintf("%s/SPS/%v/%v_revision_letter_assignment.pdf", iupopkData.Code, updateCafAssignment.IdNumber, updateCafAssignment.IdNumber)

		up, uploadErr := awshelper.UploadDocument(file, fileName)

		if uploadErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["file"] = file
			inputMap["user_id"] = claims["id"]
			inputMap["caf_assignment_id"] = idInt
			inputMap["input"] = inputUpdateCafAssignment

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": uploadErr.Error(),
			})

			createdErrLog := logs.Logs{
				CafAssignmentId: &idUint,
				Input:           inputJson,
				Message:         messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			responseErr["message"] = "failed to create caf assigment upload"
			responseErr["error"] = uploadErr.Error()
			return c.Status(400).JSON(responseErr)
		}

		if updateCafAssignment.RevisionAssignmentLetterLink == "" {
			updateDocCafAssignment, updateDocCafAssignmentErr := h.historyService.UploadUpdateDocumentCafAssignment(updateCafAssignment.ID, up.Location, uint(claims["id"].(float64)), iupopkIdInt)

			if updateDocCafAssignmentErr != nil {
				inputMap := make(map[string]interface{})
				inputMap["file"] = file
				inputMap["user_id"] = claims["id"]
				inputMap["caf_assignment_id"] = idInt
				inputMap["input"] = inputUpdateCafAssignment
				inputJson, _ := json.Marshal(inputMap)
				messageJson, _ := json.Marshal(map[string]interface{}{
					"error":           updateDocCafAssignmentErr.Error(),
					"upload_response": up,
				})

				createdErrLog := logs.Logs{
					CafAssignmentId: &idUint,
					Input:           inputJson,
					Message:         messageJson,
				}
				h.logService.CreateLogs(createdErrLog)

				responseErr["error"] = updateDocCafAssignmentErr.Error()

				status := 400

				if updateDocCafAssignmentErr.Error() == "record not found" {
					status = 404
				}

				return c.Status(status).JSON(responseErr)
			}
			return c.Status(200).JSON(updateDocCafAssignment)
		}
	}

	var createNotif notification.InputNotification

	createNotif.Type = "caf"
	createNotif.Status = "mengedit"
	createNotif.Period = updateCafAssignment.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = updateCafAssignment.ID
		inputMap["caf_assignment"] = updateCafAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(updateCafAssignment)
}

func (h *cafAssignmentHandler) DeleteCafAssignment(c *fiber.Ctx) error {
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

	detailCafAssignment, detailCafAssignmentErr := h.cafAssignmentEndUserService.DetailCafAssignment(idInt, iupopkIdInt)

	if detailCafAssignmentErr != nil {
		status := 400

		if detailCafAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete caf assigment",
			"error":   detailCafAssignmentErr.Error(),
		})
	}

	_, isDeletedCafAssignmentErr := h.historyService.DeleteCafAssignment(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedCafAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedCafAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			CafAssignmentId: &detailCafAssignment.Detail.ID,
			Input:           inputJson,
			Message:         messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedCafAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete caf assignment",
			"error":   isDeletedCafAssignmentErr.Error(),
		})
	}

	if detailCafAssignment.Detail.AssignmentLetterLink != "" || detailCafAssignment.Detail.RevisionAssignmentLetterLink != "" {
		documentLink := detailCafAssignment.Detail.AssignmentLetterLink

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
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete caf assignment aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	var createNotif notification.InputNotification

	createNotif.Type = "caf"
	createNotif.Status = "menghapus"
	createNotif.Period = detailCafAssignment.Detail.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["caf_assignment_id"] = detailCafAssignment.Detail.ID
		inputMap["caf_assignment"] = detailCafAssignment.Detail
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationErr.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete caf assignment",
	})
}
