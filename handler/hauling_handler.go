package handler

import (
	"ajebackend/model/haulingsynchronize"
	"ajebackend/model/logs"
	"ajebackend/model/transactionshauling"
	"ajebackend/model/useriupopk"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type haulingHandler struct {
	haulingSynchronizeService  haulingsynchronize.Service
	transactionsHaulingService transactionshauling.Service
	userIupopkService          useriupopk.Service
	logService                 logs.Service
}

func NewHaulingHandler(
	haulingSynchronizeService haulingsynchronize.Service,
	transactionsHaulingService transactionshauling.Service,
	userIupopkService useriupopk.Service,
	logService logs.Service,
) *haulingHandler {
	return &haulingHandler{
		haulingSynchronizeService,
		transactionsHaulingService,
		userIupopkService,
		logService,
	}
}

func (h *haulingHandler) SyncHaulingDataIsp(c *fiber.Ctx) error {

	haulingDataInput := new(haulingsynchronize.SynchronizeInputTransactionIsp)

	// Binds the request body to the Person struct
	if err := c.BodyParser(haulingDataInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed synchronize data isp",
		})
	}
	iupopkId := &haulingDataInput.IupopkId

	if *iupopkId > 0 {
		_, syncTransactionErr := h.haulingSynchronizeService.SynchronizeTransactionIsp(*haulingDataInput)

		if syncTransactionErr != nil {
			inputJson, _ := json.Marshal(haulingDataInput)

			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": syncTransactionErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"error":   syncTransactionErr.Error(),
				"message": "failed synchronize data isp",
			})
		}
	}

	syncTime := &haulingDataInput.SynchronizeTime

	var idEmpty uint = 0

	if &iupopkId == nil || *iupopkId == 0 {
		iupopkId = &idEmpty
	}

	getData, getDataErr := h.haulingSynchronizeService.GetSyncMasterDataIsp(*iupopkId)

	if getDataErr != nil {
		inputJson, _ := json.Marshal(haulingDataInput)

		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": getDataErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   getDataErr.Error(),
			"message": "failed synchronize get data master",
		})
	}

	if *iupopkId > 0 {
		_, updDataErr := h.haulingSynchronizeService.UpdateSyncMasterIsp(*iupopkId, *syncTime)

		if updDataErr != nil {
			inputJson, _ := json.Marshal(haulingDataInput)

			messageJson, _ := json.Marshal(map[string]interface{}{
				"error":     updDataErr.Error(),
				"sync_time": syncTime,
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"error":   updDataErr.Error(),
				"message": "failed update sync hauling synchronize isp",
			})
		}
	}

	return c.Status(200).JSON(getData)
}

func (h *haulingHandler) SyncHaulingDataJetty(c *fiber.Ctx) error {

	haulingDataInput := new(haulingsynchronize.SynchronizeInputTransactionJetty)

	fmt.Println(haulingDataInput)
	// Binds the request body to the Person struct
	if err := c.BodyParser(haulingDataInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed synchronize data jetty",
		})
	}
	iupopkId := &haulingDataInput.IupopkId

	if *iupopkId > 0 {
		_, syncTransactionErr := h.haulingSynchronizeService.SynchronizeTransactionJetty(*haulingDataInput)

		if syncTransactionErr != nil {
			inputJson, _ := json.Marshal(haulingDataInput)

			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": syncTransactionErr.Error(),
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"error":   syncTransactionErr.Error(),
				"message": "failed synchronize data jetty",
			})
		}

	}

	syncTime := &haulingDataInput.SynchronizeTime

	var idEmpty uint = 0

	if &iupopkId == nil || *iupopkId == 0 {
		iupopkId = &idEmpty
	}

	getData, getDataErr := h.haulingSynchronizeService.GetSyncMasterDataJetty(*iupopkId)

	if getDataErr != nil {
		inputJson, _ := json.Marshal(haulingDataInput)

		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": getDataErr.Error(),
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   getDataErr.Error(),
			"message": "failed synchronize get data master",
		})
	}

	if *iupopkId > 0 {
		_, updDataErr := h.haulingSynchronizeService.UpdateSyncMasterJetty(*iupopkId, *syncTime)

		if updDataErr != nil {
			inputJson, _ := json.Marshal(haulingDataInput)

			messageJson, _ := json.Marshal(map[string]interface{}{
				"error":     updDataErr.Error(),
				"sync_time": syncTime,
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"error":   updDataErr.Error(),
				"message": "failed update sync hauling synchronize jetty",
			})
		}
	}

	return c.Status(200).JSON(getData)
}

func (h *haulingHandler) ListStockRom(c *fiber.Ctx) error {
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

	listData, listDataErr := h.transactionsHaulingService.ListStockRom(pageNumber, iupopkIdInt)

	if listDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDataErr.Error(),
		})
	}

	return c.Status(200).JSON(listData)
}

func (h *haulingHandler) ListTransactionHauling(c *fiber.Ctx) error {
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

	listData, listDataErr := h.transactionsHaulingService.ListTransactionHauling(pageNumber, iupopkIdInt)

	if listDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listDataErr.Error(),
		})
	}

	return c.Status(200).JSON(listData)
}

func (h *haulingHandler) DetailStockRom(c *fiber.Ctx) error {
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

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
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

	detailData, detailDataErr := h.transactionsHaulingService.DetailStockRom(iupopkIdInt, idInt)

	if detailDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": detailDataErr.Error(),
		})
	}

	return c.Status(200).JSON(detailData)
}

func (h *haulingHandler) DetailTransactionHauling(c *fiber.Ctx) error {
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

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "iupopk record not found",
		})
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

	detailData, detailDataErr := h.transactionsHaulingService.DetailTransactionHauling(iupopkIdInt, idInt)

	if detailDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": detailDataErr.Error(),
		})
	}

	return c.Status(200).JSON(detailData)
}

func (h *haulingHandler) SummaryJettyTransactionPerDay(c *fiber.Ctx) error {
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

	summary, summaryErr := h.transactionsHaulingService.SummaryJettyTransactionPerDay(iupopkIdInt)

	if summaryErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": summaryErr.Error(),
		})
	}

	return c.Status(200).JSON(summary)
}

func (h *haulingHandler) SummaryInventoryStockRom(c *fiber.Ctx) error {
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

	c.Query("shipping_end")

	startDate := ""
	endDate := ""

	if c.Query("start_date") != "" {
		startDate = c.Query("start_date")
	}

	if c.Query("end_date") != "" {
		endDate = c.Query("end_date")
	}

	summary, summaryErr := h.transactionsHaulingService.SummaryInventoryStockRom(iupopkIdInt, startDate, endDate)

	if summaryErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": summaryErr.Error(),
		})
	}

	return c.Status(200).JSON(summary)
}
