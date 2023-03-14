package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/insw"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
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

type inswHandler struct {
	transactionService      transaction.Service
	historyService          history.Service
	v                       *validator.Validate
	logService              logs.Service
	groupingVesselLnService groupingvesselln.Service
	notificationUserService notificationuser.Service
	inswService             insw.Service
	userIupopkService       useriupopk.Service
}

func NewInswHandler(transactionService transaction.Service, historyService history.Service, v *validator.Validate, logService logs.Service, groupingVesselLnService groupingvesselln.Service, notificationUserService notificationuser.Service, inswService insw.Service, userIupopkService useriupopk.Service) *inswHandler {
	return &inswHandler{
		transactionService,
		historyService,
		v,
		logService,
		groupingVesselLnService,
		notificationUserService,
		inswService,
		userIupopkService,
	}
}

func (h *inswHandler) CreateInsw(c *fiber.Ctx) error {
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

	inputCreateInsw := new(insw.InputCreateInsw)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateInsw); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateInsw)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_month"] = inputCreateInsw.Month
		inputMap["insw_year"] = inputCreateInsw.Year
		inputMap["input"] = inputCreateInsw
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

	createInsw, createInswErr := h.historyService.CreateInsw(inputCreateInsw.Month, inputCreateInsw.Year, uint(claims["id"].(float64)), iupopkIdInt)

	if createInswErr != nil {
		fmt.Println(createInswErr)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_month"] = inputCreateInsw.Month
		inputMap["insw_year"] = inputCreateInsw.Year
		inputMap["input"] = inputCreateInsw
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createInswErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createInswErr.Error(),
		})
	}

	return c.Status(201).JSON(createInsw)
}

func (h *inswHandler) ListGroupingVesselLnWithPeriod(c *fiber.Ctx) error {
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

	inputCreateInsw := new(insw.InputCreateInsw)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputCreateInsw); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*inputCreateInsw)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	listGroupingVesselLnWithPeriod, listGroupingVesselLnWithPeriodErr := h.groupingVesselLnService.ListGroupingVesselLnWithPeriod(inputCreateInsw.Month, inputCreateInsw.Year, iupopkIdInt)

	if listGroupingVesselLnWithPeriodErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listGroupingVesselLnWithPeriodErr.Error(),
		})
	}

	return c.Status(200).JSON(listGroupingVesselLnWithPeriod)
}

func (h *inswHandler) DetailInsw(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt, iupopkIdInt)

	if detailInswErr != nil {

		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	return c.Status(200).JSON(detailInsw)
}

func (h *inswHandler) DeleteInsw(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt, iupopkIdInt)

	if detailInswErr != nil {

		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	_, deleteInswErr := h.historyService.DeleteInsw(idInt, uint(claims["id"].(float64)), iupopkIdInt)

	if deleteInswErr != nil {

		status := 400

		if deleteInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": deleteInswErr.Error(),
		})
	}

	if detailInsw.Detail.InswDocumentLink != "" {
		documentLinkSplit := strings.Split(detailInsw.Detail.InswDocumentLink, "/")

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
			inputMap["insw_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": deleteAwsErr.Error(),
			})

			inswId := uint(idInt)
			createdErrLog := logs.Logs{
				InswId:  &inswId,
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			status := 400

			if deleteAwsErr.Error() == "record not found" {
				status = 404
			}

			return c.Status(status).JSON(fiber.Map{
				"message": "failed to delete insw ln aws",
				"error":   deleteAwsErr.Error(),
			})
		}

	}

	var createNotif notification.InputNotification

	createNotif.Type = "insw"
	createNotif.Status = "menghapus"
	createNotif.Period = fmt.Sprintf("%s %v", detailInsw.Detail.Month, detailInsw.Detail.Year)

	_, createNotificationDeleteInsw := h.notificationUserService.CreateNotification(createNotif, uint(claims["id"].(float64)), iupopkIdInt)

	if createNotificationDeleteInsw != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["insw_id"] = idInt
		inputMap["month"] = detailInsw.Detail.Month
		inputMap["year"] = detailInsw.Detail.Year
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createNotificationDeleteInsw.Error(),
		})
		inswId := uint(idInt)
		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			InswId:  &inswId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error": createNotificationDeleteInsw.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete insw",
	})
}

func (h *inswHandler) UpdateDocumentInsw(c *fiber.Ctx) error {
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

	inputUpdateInsw := new(insw.InputUpdateDocumentInsw)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateInsw); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update insw",
		})
	}

	errors := h.v.Struct(*inputUpdateInsw)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update insw",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateInsw
		inputMap["insw_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		inswId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			InswId:  &inswId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt, iupopkIdInt)

	if detailInswErr != nil {
		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	updateInsw, updateInswErr := h.historyService.UpdateDocumentInsw(idInt, *inputUpdateInsw, uint(claims["id"].(float64)), iupopkIdInt)

	if updateInswErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateInsw
		inputMap["insw_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateInswErr.Error(),
		})

		inswId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			InswId:  &inswId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateInswErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateInswErr.Error(),
			"message": "failed to update insw",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "insw"
	inputNotification.Status = "membuat"
	inputNotification.Period = fmt.Sprintf("%v %v", detailInsw.Detail.Month, detailInsw.Detail.Year)
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateInsw
		inputMap["insw_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		inswId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
			InswId:  &inswId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update insw",
		})
	}

	return c.Status(200).JSON(updateInsw)
}

func (h *inswHandler) RequestCreateExcelInsw(c *fiber.Ctx) error {
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
			"error": "record not found",
		})
	}

	header := c.GetReqHeaders()

	detailInsw, detailInswErr := h.groupingVesselLnService.DetailInsw(idInt, iupopkIdInt)

	if detailInswErr != nil {
		status := 400

		if detailInswErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailInswErr.Error(),
		})
	}

	var inputRequestCreateExcel groupingvesselln.InputRequestCreateUploadInsw
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.Insw = detailInsw.Detail
	inputRequestCreateExcel.GroupingVesselLn = detailInsw.ListGroupingVesselLn

	hitJob, hitJobErr := h.groupingVesselLnService.RequestCreateExcelInsw(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *inswHandler) ListInsw(c *fiber.Ctx) error {
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

	var filterInsw insw.SortFilterInsw

	filterInsw.Field = c.Query("field")
	filterInsw.Sort = c.Query("sort")
	filterInsw.Month = c.Query("month")
	filterInsw.Year = c.Query("year")

	listInsw, listInswErr := h.inswService.ListInsw(pageNumber, filterInsw, iupopkIdInt)

	if listInswErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listInswErr.Error(),
		})
	}

	return c.Status(200).JSON(listInsw)
}
