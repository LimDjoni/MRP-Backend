package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/electricassignment"
	"ajebackend/model/electricassignmentenduser"
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

type electrictAssignmentHandler struct {
	electricAssignmentService        electricassignment.Service
	electricAssignmentEndUserService electricassignmentenduser.Service
	logService                       logs.Service
	userIupopkService                useriupopk.Service
	historyService                   history.Service
	notificationUserService          notificationuser.Service
	v                                *validator.Validate
	allMasterService                 allmaster.Service
}

func NewElectrictAssignmentHandler(
	electricAssignmentService electricassignment.Service,
	electricAssignmentEndUserService electricassignmentenduser.Service,
	logService logs.Service,
	userIupopkService useriupopk.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
	v *validator.Validate,
	allMasterService allmaster.Service,
) *electrictAssignmentHandler {
	return &electrictAssignmentHandler{
		electricAssignmentService,
		electricAssignmentEndUserService,
		logService,
		userIupopkService,
		historyService,
		notificationUserService,
		v,
		allMasterService,
	}
}

func (h *electrictAssignmentHandler) ListElectricAssignment(c *fiber.Ctx) error {
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

	var filterElectricAssignment electricassignment.SortFilterElectricAssignment

	filterElectricAssignment.Quantity = c.Query("quantity")
	filterElectricAssignment.Year = c.Query("year")
	filterElectricAssignment.Field = c.Query("field")
	filterElectricAssignment.Sort = c.Query("sort")

	listElectricAssignment, listElectricAssignmentErr := h.electricAssignmentService.ListElectricAssignment(pageNumber, filterElectricAssignment, iupopkIdInt)

	if listElectricAssignmentErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listElectricAssignmentErr.Error(),
		})
	}

	return c.Status(200).JSON(listElectricAssignment)
}

func (h *electrictAssignmentHandler) CreateElectricAssignment(c *fiber.Ctx) error {
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

	inputCreateElectricAssignment := new(electricassignmentenduser.CreateElectricAssignmentInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputCreateElectricAssignment)

	if errParsing == nil {
		formPart, errFormPart := c.MultipartForm()
		if errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}
		if len(inputCreateElectricAssignment.ListElectricAssignment) == 0 {
			for _, value := range formPart.Value["list_electric_assignment"] {

				var tempElectricAssignment []electricassignmentenduser.ElectricAssignmentInput

				errUnmarshal := json.Unmarshal([]byte(value), &tempElectricAssignment)

				fmt.Println(errUnmarshal)
				inputCreateElectricAssignment.ListElectricAssignment = tempElectricAssignment
			}

			if len(inputCreateElectricAssignment.ListElectricAssignment) == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "please check there is no electric assignment data",
				})
			}
		}
	}

	errors := h.v.Struct(*inputCreateElectricAssignment)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateElectricAssignment
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

	file, errFormFile := c.FormFile("letter_assignment")

	responseErr := fiber.Map{
		"message": "failed to create electric assignment",
	}

	if errFormFile != nil {
		responseErr["error"] = errFormFile.Error()
		return c.Status(400).JSON(responseErr)
	}

	if !strings.Contains(file.Filename, ".pdf") {
		responseErr["error"] = "document must be pdf"
		return c.Status(400).JSON(responseErr)
	}

	createElectricAssignment, createElectricAssignmentErr := h.historyService.CreateElectricAssignment(*inputCreateElectricAssignment, uint(claims["id"].(float64)), iupopkIdInt)

	if createElectricAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateElectricAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createElectricAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createElectricAssignmentErr.Error(),
		})
	}

	fileName := fmt.Sprintf("%s/SPK/%v/%v_letter_assignment.pdf", iupopkData.Code, createElectricAssignment.IdNumber, createElectricAssignment.IdNumber)

	up, uploadErr := awshelper.UploadDocument(file, fileName)

	if uploadErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateElectricAssignment

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadErr.Error(),
		})

		createdErrLog := logs.Logs{
			ElectricAssignmentId: &createElectricAssignment.ID,
			Input:                inputJson,
			Message:              messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		responseErr["message"] = "failed to create electric assigment upload"
		responseErr["error"] = uploadErr.Error()
		return c.Status(400).JSON(responseErr)
	}

	updateElectricAssignment, updateElectricAssignmentErr := h.historyService.UploadCreateDocumentElectricAssignment(createElectricAssignment.ID, up.Location, uint(claims["id"].(float64)), iupopkIdInt)

	if updateElectricAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["file"] = file
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputCreateElectricAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":           updateElectricAssignmentErr.Error(),
			"upload_response": up,
		})

		createdErrLog := logs.Logs{
			ElectricAssignmentId: &createElectricAssignment.ID,
			Input:                inputJson,
			Message:              messageJson,
		}
		h.logService.CreateLogs(createdErrLog)

		responseErr["error"] = updateElectricAssignmentErr.Error()

		status := 400

		if updateElectricAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(responseErr)
	}

	var createNotif notification.InputNotification

	createNotif.Type = "electric assignment"
	createNotif.Status = "membuat"
	createNotif.Period = updateElectricAssignment.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = updateElectricAssignment.ID
		inputMap["electric_assignment"] = updateElectricAssignment
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

	return c.Status(201).JSON(updateElectricAssignment)
}

func (h *electrictAssignmentHandler) DetailElectricAssignment(c *fiber.Ctx) error {
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

	detailElectricAssignment, detailElectricAssignmentErr := h.electricAssignmentEndUserService.DetailElectricAssignment(idInt, iupopkIdInt)

	if detailElectricAssignmentErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": detailElectricAssignmentErr.Error(),
		})
	}

	return c.Status(200).JSON(detailElectricAssignment)
}

func (h *electrictAssignmentHandler) UpdateElectricAssignment(c *fiber.Ctx) error {
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

	inputUpdateElectricAssignment := new(electricassignmentenduser.UpdateElectricAssignmentInput)

	// Binds the request body to the Person struct
	errParsing := c.BodyParser(inputUpdateElectricAssignment)

	if errParsing == nil {
		formPart, errFormPart := c.MultipartForm()
		if errFormPart != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "please check there is no data",
			})
		}
		if len(inputUpdateElectricAssignment.ListElectricAssignment) == 0 {
			for _, value := range formPart.Value["list_electric_assignment"] {

				var tempElectricAssignment []electricassignmentenduser.ElectricAssignmentEndUser

				errUnmarshal := json.Unmarshal([]byte(value), &tempElectricAssignment)

				fmt.Println(errUnmarshal)

				inputUpdateElectricAssignment.ListElectricAssignment = tempElectricAssignment
			}

			if len(inputUpdateElectricAssignment.ListElectricAssignment) == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "please check there is no electric assignment data",
				})
			}

		}
	}

	errors := h.v.Struct(*inputUpdateElectricAssignment)
	idUint := uint(idInt)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = idInt
		inputMap["input"] = inputUpdateElectricAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		createdErrLog := logs.Logs{
			ElectricAssignmentId: &idUint,
			Input:                inputJson,
			Message:              messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	updateElectricAssignment, updateElectricAssignmentErr := h.historyService.UpdateElectricAssignment(idInt, *inputUpdateElectricAssignment, uint(claims["id"].(float64)), iupopkIdInt)

	if updateElectricAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = idInt
		inputMap["input"] = inputUpdateElectricAssignment
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateElectricAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			ElectricAssignmentId: &idUint,
			Input:                inputJson,
			Message:              messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": updateElectricAssignmentErr.Error(),
		})
	}

	file, errFormFile := c.FormFile("revision_letter_assignment")

	if errFormFile == nil {

		responseErr := fiber.Map{
			"message": "failed to create electric assignment",
		}

		if errFormFile != nil {
			responseErr["error"] = errFormFile.Error()
			return c.Status(400).JSON(responseErr)
		}

		if !strings.Contains(file.Filename, ".pdf") {
			responseErr["error"] = "document must be pdf"
			return c.Status(400).JSON(responseErr)
		}

		fileName := fmt.Sprintf("%s/SPK/%v/%v_revision_letter_assignment.pdf", iupopkData.Code, updateElectricAssignment.IdNumber, updateElectricAssignment.IdNumber)

		up, uploadErr := awshelper.UploadDocument(file, fileName)

		if uploadErr != nil {
			inputMap := make(map[string]interface{})
			inputMap["file"] = file
			inputMap["user_id"] = claims["id"]
			inputMap["electric_assignment_id"] = idInt
			inputMap["input"] = inputUpdateElectricAssignment

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": uploadErr.Error(),
			})

			createdErrLog := logs.Logs{
				ElectricAssignmentId: &idUint,
				Input:                inputJson,
				Message:              messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			responseErr["message"] = "failed to create electric assigment upload"
			responseErr["error"] = uploadErr.Error()
			return c.Status(400).JSON(responseErr)
		}

		if updateElectricAssignment.RevisionAssignmentLetterLink == "" {
			updateDocElectricAssignment, updateDocElectricAssignmentErr := h.historyService.UploadUpdateDocumentElectricAssignment(updateElectricAssignment.ID, up.Location, uint(claims["id"].(float64)), iupopkIdInt)

			if updateDocElectricAssignmentErr != nil {
				inputMap := make(map[string]interface{})
				inputMap["file"] = file
				inputMap["user_id"] = claims["id"]
				inputMap["electric_assignment_id"] = idInt
				inputMap["input"] = inputUpdateElectricAssignment
				inputJson, _ := json.Marshal(inputMap)
				messageJson, _ := json.Marshal(map[string]interface{}{
					"error":           updateDocElectricAssignmentErr.Error(),
					"upload_response": up,
				})

				createdErrLog := logs.Logs{
					ElectricAssignmentId: &idUint,
					Input:                inputJson,
					Message:              messageJson,
				}
				h.logService.CreateLogs(createdErrLog)

				responseErr["error"] = updateDocElectricAssignmentErr.Error()

				status := 400

				if updateDocElectricAssignmentErr.Error() == "record not found" {
					status = 404
				}

				return c.Status(status).JSON(responseErr)
			}
			return c.Status(200).JSON(updateDocElectricAssignment)
		}
	}

	var createNotif notification.InputNotification

	createNotif.Type = "electric assignment"
	createNotif.Status = "mengedit"
	createNotif.Period = updateElectricAssignment.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = updateElectricAssignment.ID
		inputMap["electric_assignment"] = updateElectricAssignment
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

	return c.Status(200).JSON(updateElectricAssignment)
}

func (h *electrictAssignmentHandler) DeleteElectricAssignment(c *fiber.Ctx) error {
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

	detailElectricAssignment, detailElectricAssignmentErr := h.electricAssignmentEndUserService.DetailElectricAssignment(idInt, iupopkIdInt)

	if detailElectricAssignmentErr != nil {
		status := 400

		if detailElectricAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete electric assigment",
			"error":   detailElectricAssignmentErr.Error(),
		})
	}

	_, isDeletedElectricAssignmentErr := h.historyService.DeleteElectricAssignment(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedElectricAssignmentErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedElectricAssignmentErr.Error(),
		})

		createdErrLog := logs.Logs{
			ElectricAssignmentId: &detailElectricAssignment.Detail.ID,
			Input:                inputJson,
			Message:              messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedElectricAssignmentErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete electric assignment",
			"error":   isDeletedElectricAssignmentErr.Error(),
		})
	}

	if detailElectricAssignment.Detail.AssignmentLetterLink != "" || detailElectricAssignment.Detail.RevisionAssignmentLetterLink != "" {
		documentLink := detailElectricAssignment.Detail.AssignmentLetterLink

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
				"message": "failed to delete electric assignment aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	var createNotif notification.InputNotification

	createNotif.Type = "caf"
	createNotif.Status = "menghapus"
	createNotif.Period = detailElectricAssignment.Detail.Year

	_, createNotificationErr := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["electric_assignment_id"] = detailElectricAssignment.Detail.ID
		inputMap["electric_assignment"] = detailElectricAssignment.Detail
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
		"message": "success delete electric assignment",
	})
}
