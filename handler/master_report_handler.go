package handler

import (
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"ajebackend/model/transactionrequestreport"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xuri/excelize/v2"
)

type masterReportHandler struct {
	masterReportService      masterreport.Service
	userIupopkService        useriupopk.Service
	v                        *validator.Validate
	allMasterService         allmaster.Service
	transactionRequestReport transactionrequestreport.Service
}

func NewMasterReportHandler(masterReportService masterreport.Service, userIupopkService useriupopk.Service, v *validator.Validate, allMasterService allmaster.Service, transactionRequestReport transactionrequestreport.Service) *masterReportHandler {
	return &masterReportHandler{
		masterReportService,
		userIupopkService,
		v,
		allMasterService,
		transactionRequestReport,
	}
}

func (h *masterReportHandler) RecapDmo(c *fiber.Ctx) error {
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

	year := c.Params("year")

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	reportRecap, reportRecapErr := h.masterReportService.RecapDmo(year, iupopkIdInt)

	if reportRecapErr != nil {
		status := 400

		if reportRecapErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": reportRecapErr.Error(),
		})
	}

	return c.Status(200).JSON(reportRecap)
}

func (h *masterReportHandler) DownloadRecapDmo(c *fiber.Ctx) error {
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

	year := c.Params("year")

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

	reportRecap, reportRecapErr := h.masterReportService.RecapDmo(year, iupopkIdInt)

	if reportRecapErr != nil {
		status := 400

		if reportRecapErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": reportRecapErr.Error(),
		})
	}

	excelFile := excelize.NewFile()
	excelFile.NewSheet("Rekapitulasi")

	defer func() {
		err := os.Remove("./Book1.xlsx")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	excelFile, errRecap := h.masterReportService.CreateReportRecapDmo(year, reportRecap, iupopkData, excelFile, "Rekapitulasi")

	if errRecap != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errRecap.Error(),
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

	fileName := fmt.Sprintf("rekapitulasi-%s-%s", iupopkData.Code, year)

	return c.Status(200).Download("./Book1.xlsx", fileName)
}

func (h *masterReportHandler) RealizationReport(c *fiber.Ctx) error {
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

	year := c.Params("year")

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	realizationReport, realizationReportErr := h.masterReportService.RealizationReport(year, iupopkIdInt)

	if realizationReportErr != nil {
		status := 400

		if realizationReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": realizationReportErr.Error(),
		})
	}

	return c.Status(200).JSON(realizationReport)
}

func (h *masterReportHandler) DownloadRealizationReport(c *fiber.Ctx) error {
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

	year := c.Params("year")

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

	realizationReport, realizationReportErr := h.masterReportService.RealizationReport(year, iupopkIdInt)

	if realizationReportErr != nil {
		status := 400

		if realizationReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": realizationReportErr.Error(),
		})
	}

	excelFile := excelize.NewFile()
	excelFile.NewSheet("JAN")
	excelFile.NewSheet("FEB")
	excelFile.NewSheet("MAR")
	excelFile.NewSheet("APR")
	excelFile.NewSheet("MEI")
	excelFile.NewSheet("JUN")
	excelFile.NewSheet("JUL")
	excelFile.NewSheet("AGU")
	excelFile.NewSheet("SEP")
	excelFile.NewSheet("OKT")
	excelFile.NewSheet("NOV")
	excelFile.NewSheet("DES")

	defer func() {
		err := os.Remove("./Book1.xlsx")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	excelFile, errRealization := h.masterReportService.CreateReportRealization(year, realizationReport, iupopkData, excelFile)

	if errRealization != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errRealization.Error(),
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

	fileName := fmt.Sprintf("realisasi-%s-%s", iupopkData.Code, year)

	return c.Status(200).Download("./Book1.xlsx", fileName)
}

func (h *masterReportHandler) SaleDetailReport(c *fiber.Ctx) error {
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

	year := c.Params("year")

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	saleDetailReport, saleDetailReportErr := h.masterReportService.SaleDetailReport(year, iupopkIdInt)

	if saleDetailReportErr != nil {
		status := 400

		if saleDetailReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": saleDetailReportErr.Error(),
		})
	}

	return c.Status(200).JSON(saleDetailReport)
}

func (h *masterReportHandler) DownloadSaleDetailReport(c *fiber.Ctx) error {
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

	year := c.Params("year")

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

	saleDetailReport, saleDetailReportErr := h.masterReportService.SaleDetailReport(year, iupopkIdInt)

	if saleDetailReportErr != nil {
		status := 400

		if saleDetailReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": saleDetailReportErr.Error(),
		})
	}

	excelFile := excelize.NewFile()
	excelFile.NewSheet("Detail")

	excelFile.NewSheet("Chart")

	defer func() {
		err := os.Remove("./Book1.xlsx")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	excelFile, errSaleDetail := h.masterReportService.CreateReportSalesDetail(year, saleDetailReport, iupopkData, excelFile, "Detail", "Chart")

	if errSaleDetail != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errSaleDetail.Error(),
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

	fileName := fmt.Sprintf("detail-%s-%s", iupopkData.Code, year)

	return c.Status(200).Download("./Book1.xlsx", fileName)
}

func (h *masterReportHandler) GetTransactionReport(c *fiber.Ctx) error {
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

	iupopkData, iupopkDataErr := h.allMasterService.FindIupopk(iupopkIdInt)

	if iupopkDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkDataErr.Error(),
		})
	}

	typeTransaction := c.Params("type")

	if typeTransaction != "dn" && typeTransaction != "ln" {
		return c.Status(404).JSON(fiber.Map{
			"error": "type transaction record not found",
		})
	}

	checkUser, checkUserErr := h.userIupopkService.FindUser(uint(claims["id"].(float64)), iupopkIdInt)

	if checkUserErr != nil || checkUser.IsActive == false {
		return c.Status(401).JSON(responseUnauthorized)
	}

	transactionReportInput := new(masterreport.TransactionReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(transactionReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*transactionReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	transactionReport, transactionReportErr := h.masterReportService.GetTransactionReport(iupopkIdInt, *transactionReportInput, typeTransaction)

	if transactionReportErr != nil {
		status := 400

		if transactionReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": transactionReportErr.Error(),
		})
	}

	file, err := excelize.OpenFile("./assets/template/Template Transaction.xlsx")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer func() {
		err := os.Remove("./Book1.xlsx")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	fileTransaction, fileTransactionErr := h.masterReportService.CreateTransactionReport(file, "Sheet1", iupopkData, transactionReport)

	if fileTransactionErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fileTransactionErr.Error(),
		})
	}

	if err := fileTransaction.SaveAs("Book1.xlsx"); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "failed to get report",
			"error":   err.Error(),
		})
	}

	fileName := fmt.Sprintf("Data Transaksi %s", iupopkData.Code)

	return c.Status(200).Download("./Book1.xlsx", fileName)
}

func (h *masterReportHandler) CreateTransactionRequestReport(c *fiber.Ctx) error {
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

	transactionReportInput := new(masterreport.TransactionReportInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(transactionReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*transactionReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createTransactionReqReport, createTransactionReqReportErr := h.transactionRequestReport.CreateTransactionRequestReport(*transactionReportInput, iupopkIdInt)

	if createTransactionReqReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createTransactionReqReportErr.Error(),
		})
	}
}
