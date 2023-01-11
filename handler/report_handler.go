package handler

import (
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xuri/excelize/v2"
)

type reportHandler struct {
	transactionService transaction.Service
	userService        user.Service
	v                  *validator.Validate
	logService         logs.Service
}

func NewReportHandler(transactionService transaction.Service, userService user.Service, v *validator.Validate, logService logs.Service) *reportHandler {
	return &reportHandler{
		transactionService,
		userService,
		v,
		logService,
	}
}

func (h *reportHandler) Report(c *fiber.Ctx) error {
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

	reportInput := new(transaction.InputRequestGetReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(reportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*reportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	if reportInput.Year == 0 {
		year, _, _ := time.Now().Date()
		reportInput.Year = year
	}

	reportRecap, reportDetail, reportErr := h.transactionService.GetReport(reportInput.Year)

	if reportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   reportErr.Error(),
		})
	}

	reportRecap.Year = reportInput.Year
	reportRecap.ProductionPlan = reportInput.ProductionPlan
	reportRecap.PercentageProductionObligation = reportInput.PercentageProductionObligation
	reportRecap.ProductionObligation = reportInput.ProductionPlan * reportInput.PercentageProductionObligation / 100

	percentage := reportRecap.Total / reportRecap.ProductionObligation * 100
	reportRecap.FulfillmentPercentageProductionObligation = fmt.Sprintf("%.2f%%", percentage)
	reportRecap.ProrateProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total/reportRecap.ProductionPlan*100)
	reportRecap.FulfillmentOfProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total/reportRecap.ProductionPlan*100)

	return c.Status(200).JSON(fiber.Map{
		"detail": reportDetail,
		"recap":  reportRecap,
	})
}

func (h *reportHandler) DownloadReport(c *fiber.Ctx) error {
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

	queryYear := c.Query("year")
	var queryYearInt int
	queryProductionPlan := c.Query("production_plan")

	queryProductionPlanFloat, errProductionPlan := strconv.ParseFloat(queryProductionPlan, 64)
	if errProductionPlan != nil {
		fmt.Println(errProductionPlan, 1)
	}

	queryPercentageProductionObligation := c.Query("percentage_production_obligation")

	queryPercentageProductionObligationFloat, errPercentageProductionObligationFloat := strconv.ParseFloat(queryPercentageProductionObligation, 64)
	if errPercentageProductionObligationFloat != nil {
		fmt.Println(errPercentageProductionObligationFloat, 2)
	}

	if queryYear == "" {
		year, _, _ := time.Now().Date()
		queryYearInt = year
	} else {
		intYear, errParse := strconv.Atoi(queryYear)

		if errParse != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "failed to get report",
				"error":   "mohon masukkan tahun yang valid",
			})
		}
		queryYearInt = intYear
	}

	reportRecap, reportDetail, reportErr := h.transactionService.GetReport(queryYearInt)

	if reportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   reportErr.Error(),
		})
	}

	reportRecap.Year = queryYearInt
	reportRecap.ProductionPlan = queryProductionPlanFloat
	reportRecap.PercentageProductionObligation = queryPercentageProductionObligationFloat
	reportRecap.ProductionObligation = queryProductionPlanFloat * queryPercentageProductionObligationFloat / 100

	percentage := reportRecap.Total / reportRecap.ProductionObligation * 100
	reportRecap.FulfillmentPercentageProductionObligation = fmt.Sprintf("%.2f%%", percentage)
	reportRecap.ProrateProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total/reportRecap.ProductionPlan*100)
	reportRecap.FulfillmentOfProductionPlan = fmt.Sprintf("%.2f%%", reportRecap.Total/reportRecap.ProductionPlan*100)

	var mapDataElectricity map[string]map[string]float64
	marshalDataElectricity, _ := json.Marshal(reportDetail.Electricity)
	json.Unmarshal(marshalDataElectricity, &mapDataElectricity)

	var mapDataNonElectricity map[string]map[string]float64
	marshalDataNonElectricity, _ := json.Marshal(reportDetail.NonElectricity)
	json.Unmarshal(marshalDataNonElectricity, &mapDataNonElectricity)

	var mapRecap map[string]map[string]map[string]float64
	marshalRecap, _ := json.Marshal(reportDetail)
	json.Unmarshal(marshalRecap, &mapRecap)

	var mapRecapProduction map[string]map[string]float64
	marshalRecapProduction, _ := json.Marshal(reportDetail)
	json.Unmarshal(marshalRecapProduction, &mapRecapProduction)

	var mapRecapReport map[string]interface{}
	marshalRecapReport, _ := json.Marshal(reportRecap)
	json.Unmarshal(marshalRecapReport, &mapRecapReport)

	excelFile := excelize.NewFile()

	excelFile.NewSheet("Rekapitulasi")
	excelFile, errRecap := helper.CreateReportRecap("Rekapitulasi", excelFile, mapRecap, mapRecapProduction, mapRecapReport)
	if errRecap != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   errRecap.Error(),
		})
	}

	excelFile.NewSheet("Kelistrikan")
	excelFile, err := helper.CreateReportDetailCompany(reportDetail.ElectricityCompany, "Kelistrikan", excelFile, mapDataElectricity, "Realisasi Penjualan Batu Bara Kelistrikan")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   err.Error(),
		})
	}

	excelFile.NewSheet("NonKelistrikan")
	excelFile, errNonElectric := helper.CreateReportDetailCompany(reportDetail.NonElectricityCompany, "NonKelistrikan", excelFile, mapDataNonElectricity, "Realisasi Penjualan Batu Bara Non-Kelistrikan")

	if errNonElectric != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   errNonElectric.Error(),
		})
	}

	excelFile.DeleteSheet("Sheet1")
	excelFile.SetActiveSheet(0)

	if err := excelFile.SaveAs("Book1.xlsx"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   err.Error(),
		})
	}

	fileName := fmt.Sprintf("report-%v", queryYearInt)

	return c.Status(200).Download("./Book1.xlsx", fileName)
}
