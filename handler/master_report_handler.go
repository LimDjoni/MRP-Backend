package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/royaltyrecon"
	"ajebackend/model/royaltyreport"
	"ajebackend/model/useriupopk"
	"ajebackend/validatorfunc"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"ajebackend/model/transactionrequestreport"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xuri/excelize/v2"
)

type masterReportHandler struct {
	masterReportService             masterreport.Service
	userIupopkService               useriupopk.Service
	v                               *validator.Validate
	allMasterService                allmaster.Service
	transactionRequestReportService transactionrequestreport.Service
	logService                      logs.Service
	royaltyReconService             royaltyrecon.Service
	royaltyReportService            royaltyreport.Service
	historyService                  history.Service
	notificationUserService         notificationuser.Service
}

func NewMasterReportHandler(
	masterReportService masterreport.Service,
	userIupopkService useriupopk.Service,
	v *validator.Validate,
	allMasterService allmaster.Service,
	transactionRequestReportService transactionrequestreport.Service,
	logService logs.Service,
	royaltyReconService royaltyrecon.Service,
	royaltyReportService royaltyreport.Service,
	historyService history.Service,
	notificationUserService notificationuser.Service,
) *masterReportHandler {
	return &masterReportHandler{
		masterReportService,
		userIupopkService,
		v,
		allMasterService,
		transactionRequestReportService,
		logService,
		royaltyReconService,
		royaltyReportService,
		historyService,
		notificationUserService,
	}
}

// Report Recap DMO
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

// Report Realisasi DMO
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

// Report Detail Penjualan DMO
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

// Request Transaksi Report (All Data)
func (h *masterReportHandler) PreviewTransactionReport(c *fiber.Ctx) error {
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

	transactionReport, transactionReportErr := h.transactionRequestReportService.PreviewTransactionReport(*transactionReportInput, iupopkIdInt)

	if transactionReportErr != nil {
		status := 400

		if transactionReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": transactionReportErr.Error(),
		})
	}

	return c.Status(200).JSON(transactionReport)
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

	iupopkData, iupopkDataErr := h.allMasterService.FindIupopk(iupopkIdInt)

	if iupopkDataErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": iupopkDataErr.Error(),
		})
	}

	header := c.GetReqHeaders()

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

	createTransactionReqReport, createTransactionReqReportErr := h.transactionRequestReportService.CreateTransactionReport(*transactionReportInput, iupopkIdInt)

	if createTransactionReqReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createTransactionReqReportErr.Error(),
		})
	}

	var input transactionrequestreport.InputRequestJobReportTransaction

	input.Authorization = header["Authorization"]
	input.Id = createTransactionReqReport.ID
	input.Iupopk = iupopkData

	reqJob, reqJobErr := h.transactionRequestReportService.ReqJobCreateTransactionReport(input)

	if reqJobErr != nil {

		_, updErr := h.transactionRequestReportService.UpdateTransactionReportError(int(createTransactionReqReport.ID), iupopkIdInt)

		maps := make(map[string]interface{})
		maps["error"] = reqJobErr.Error()
		maps["error_update"] = nil
		if updErr != nil {
			maps["error_update"] = updErr.Error()
		}
		maps["message"] = "failed to create job"

		return c.Status(400).JSON(maps)
	}

	return c.Status(201).JSON(reqJob)
}

func (h *masterReportHandler) UpdateJobTransactionRequestReport(c *fiber.Ctx) error {
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

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "id record not found",
		})
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

	detailTransactionReportRequest, detailTransactionReportRequestErr := h.transactionRequestReportService.DetailTransactionReport(idInt, iupopkIdInt)

	if detailTransactionReportRequestErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": detailTransactionReportRequestErr.Error(),
		})
	}
	var sheetName = "Sheet1"

	fileDn, err := excelize.OpenFile("./assets/template/Template Transaction.xlsx")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fileLn, err := excelize.OpenFile("./assets/template/Template Transaction.xlsx")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	dnReport, dnReportErr := h.masterReportService.CreateTransactionReport(fileDn, sheetName, iupopkData, detailTransactionReportRequest.ListDnTransactions)

	if dnReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": dnReportErr.Error(),
		})
	}

	lnReport, lnReportErr := h.masterReportService.CreateTransactionReport(fileLn, sheetName, iupopkData, detailTransactionReportRequest.ListLnTransactions)

	if lnReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": lnReportErr.Error(),
		})
	}

	dnFile, errDnFile := dnReport.WriteToBuffer()

	if errDnFile != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errDnFile.Error(),
		})
	}

	lnFile, errLnFile := lnReport.WriteToBuffer()

	if errLnFile != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": errLnFile.Error(),
		})
	}

	fileNameDn := fmt.Sprintf("%s/TRR/%s/%s_%s.xlsx", iupopkData.Code, *detailTransactionReportRequest.Detail.IdNumber, *detailTransactionReportRequest.Detail.IdNumber, "transaction_dn")

	upDn, uploadDnErr := awshelper.UploadDocumentFromExcelize(dnFile, fileNameDn)

	if uploadDnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["document_type"] = "dn"
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_report_request_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadDnErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionRequestReportId: &transactionId,
			Input:                      inputJson,
			Message:                    messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		_, updTransErr := h.transactionRequestReportService.UpdateTransactionReportError(idInt, iupopkIdInt)

		maps := make(map[string]interface{})

		maps["error"] = uploadDnErr.Error()
		maps["error_update"] = nil

		if updTransErr != nil {
			maps["error_update"] = updTransErr.Error()
		}

		return c.Status(400).JSON(maps)
	}

	fileNameLn := fmt.Sprintf("%s/TRR/%s/%s_%s.xlsx", iupopkData.Code, *detailTransactionReportRequest.Detail.IdNumber, *detailTransactionReportRequest.Detail.IdNumber, "transaction_ln")

	upLn, uploadLnErr := awshelper.UploadDocumentFromExcelize(lnFile, fileNameLn)

	if uploadLnErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["document_type"] = "ln"
		inputMap["user_id"] = claims["id"]
		inputMap["transaction_report_request_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": uploadLnErr.Error(),
		})

		transactionId := uint(idInt)
		createdErrLog := logs.Logs{
			TransactionRequestReportId: &transactionId,
			Input:                      inputJson,
			Message:                    messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		_, updTransErr := h.transactionRequestReportService.UpdateTransactionReportError(idInt, iupopkIdInt)

		maps := make(map[string]interface{})

		maps["error"] = uploadDnErr.Error()
		maps["error_update"] = nil

		if updTransErr != nil {
			maps["error_update"] = updTransErr.Error()
		}

		return c.Status(400).JSON(maps)
	}

	finishUpdTransaction, finishUpdTransactionErr := h.transactionRequestReportService.UpdateTransactionReport(upDn.Location, upLn.Location, idInt, iupopkIdInt)

	if finishUpdTransactionErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": finishUpdTransactionErr.Error(),
		})
	}

	return c.Status(200).JSON(finishUpdTransaction)
}

func (h *masterReportHandler) DetailTransactionReport(c *fiber.Ctx) error {
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

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "id record not found",
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

	detailTransactionReport, detailTransactionReportErr := h.transactionRequestReportService.DetailTransactionReport(idInt, iupopkIdInt)

	if detailTransactionReportErr != nil {
		status := 400

		if detailTransactionReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailTransactionReportErr.Error(),
		})
	}

	return c.Status(200).JSON(detailTransactionReport)
}

func (h *masterReportHandler) DeleteTransactionReportById(c *fiber.Ctx) error {
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

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "id record not found",
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

	detail, detailErr := h.transactionRequestReportService.DetailTransactionReport(idInt, iupopkIdInt)

	if detailErr != nil {
		status := 400

		if detailErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailErr.Error(),
		})
	}

	_, isDeletedErr := h.transactionRequestReportService.DeleteTransactionReportById(idInt, iupopkIdInt)

	if isDeletedErr != nil {
		status := 400
		if isDeletedErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": isDeletedErr.Error(),
		})
	}

	if detail.Detail.DocumentDnLink != "" || detail.Detail.DocumentDnLink != "" {
		documentLink := detail.Detail.DocumentDnLink

		if documentLink == "" {
			documentLink = detail.Detail.DocumentLnLink
		}

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
			inputMap["transaction_request_report_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error":     deleteAwsErr.Error(),
				"id_number": detail.Detail.IdNumber,
			})

			createdErrLog := logs.Logs{
				Input:   inputJson,
				Message: messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete transaction request report aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success deleted transaction request report",
	})
}

func (h *masterReportHandler) DeleteTransactionReport(c *fiber.Ctx) error {
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

	listDeleted, listDeletedErr := h.transactionRequestReportService.ListDeletedTransactionReport()

	if listDeletedErr != nil {
		status := 400
		if listDeletedErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": listDeletedErr.Error(),
		})
	}

	_, isDeletedErr := h.transactionRequestReportService.DeleteTransactionReport()

	if isDeletedErr != nil {
		status := 400
		if isDeletedErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": isDeletedErr.Error(),
		})
	}

	for _, v := range listDeleted {
		if v.DocumentDnLink != "" || v.DocumentDnLink != "" {
			documentLink := v.DocumentDnLink

			if documentLink == "" {
				documentLink = v.DocumentLnLink
			}

			documentLinkSplit := strings.Split(documentLink, "/")

			fileName := ""
			for index, value := range documentLinkSplit {
				if index == 3 {
					fileName += value + "/"
				}

				if index == 4 {
					fileName += value + "/"
				}

				if index == 5 {
					fileName += value
				}
			}

			_, deleteAwsErr := awshelper.DeleteDocumentBatch(fileName)

			if deleteAwsErr != nil {
				inputMap := make(map[string]interface{})
				inputMap["user_id"] = claims["id"]

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
					"message": "failed to delete transaction request report aws",
					"error":   deleteAwsErr.Error(),
				})
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success deleted transaction request report",
	})
}

func (h *masterReportHandler) ListTransactionReport(c *fiber.Ctx) error {
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

	list, listErr := h.transactionRequestReportService.ListTransactionReport(pageNumber, iupopkIdInt)

	if listErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listErr.Error(),
		})
	}

	return c.Status(200).JSON(list)
}

// Report Royalty Recon
func (h *masterReportHandler) PreviewRoyaltyRecon(c *fiber.Ctx) error {
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

	royaltyReconInput := new(royaltyrecon.InputRoyaltyRecon)

	// Binds the request body to the Person struct
	if err := c.BodyParser(royaltyReconInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*royaltyReconInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	royaltyRecon, royaltyReconErr := h.royaltyReconService.GetTransactionRoyaltyRecon(royaltyReconInput.DateFrom, royaltyReconInput.DateTo, iupopkIdInt)

	if royaltyReconErr != nil {
		status := 400

		if royaltyReconErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": royaltyReconErr.Error(),
		})
	}

	return c.Status(200).JSON(royaltyRecon)
}

func (h *masterReportHandler) CreateRoyaltyRecon(c *fiber.Ctx) error {
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

	royaltyReconInput := new(royaltyrecon.InputRoyaltyRecon)

	// Binds the request body to the Person struct
	if err := c.BodyParser(royaltyReconInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*royaltyReconInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	royaltyRecon, royaltyReconErr := h.historyService.CreateRoyaltyRecon(royaltyReconInput.DateFrom, royaltyReconInput.DateTo, iupopkIdInt, uint(claims["id"].(float64)))

	if royaltyReconErr != nil {
		inputJson, _ := json.Marshal(royaltyReconInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":   royaltyReconErr.Error(),
			"message": "create royalty recon err",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": royaltyReconErr.Error(),
		})
	}

	return c.Status(201).JSON(royaltyRecon)
}

func (h *masterReportHandler) DeleteRoyaltyRecon(c *fiber.Ctx) error {
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

	detailRoyaltyRecon, detailRoyaltyReconErr := h.royaltyReconService.GetDetailTransactionRoyaltyRecon(idInt, iupopkIdInt)

	if detailRoyaltyReconErr != nil {
		status := 400

		if detailRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete royalty recon",
			"error":   detailRoyaltyReconErr.Error(),
		})
	}

	_, isDeletedRoyaltyReconErr := h.historyService.DeleteRoyaltyRecon(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedRoyaltyReconErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["royalty_recon_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedRoyaltyReconErr.Error(),
		})

		royaltyReconId := uint(idInt)
		createdErrLog := logs.Logs{
			RoyaltyReconId: &royaltyReconId,
			Input:          inputJson,
			Message:        messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete royalty recon",
			"error":   isDeletedRoyaltyReconErr.Error(),
		})
	}

	if detailRoyaltyRecon.Detail.RoyaltyReconDocumentLink != "" {
		documentLink := detailRoyaltyRecon.Detail.RoyaltyReconDocumentLink

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
			inputMap["royalty_recon_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": isDeletedRoyaltyReconErr.Error(),
			})

			royaltyReconId := uint(idInt)
			createdErrLog := logs.Logs{
				RoyaltyReconId: &royaltyReconId,
				Input:          inputJson,
				Message:        messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete royalty recon aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete royalty recon",
	})
}

func (h *masterReportHandler) DetailRoyaltyRecon(c *fiber.Ctx) error {
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

	detailRoyaltyRecon, detailRoyaltyReconErr := h.royaltyReconService.GetDetailTransactionRoyaltyRecon(idInt, iupopkIdInt)

	if detailRoyaltyReconErr != nil {
		status := 400

		if detailRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReconErr.Error(),
		})
	}

	return c.Status(200).JSON(detailRoyaltyRecon)
}

func (h *masterReportHandler) UpdateDocumentRoyaltyRecon(c *fiber.Ctx) error {
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

	inputUpdateRoyaltyRecon := new(royaltyrecon.InputUpdateDocumentRoyaltyRecon)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateRoyaltyRecon); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update royalty recon",
		})
	}

	errors := h.v.Struct(*inputUpdateRoyaltyRecon)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update royalty recon",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyRecon
		inputMap["royalty_recon_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		royaltyReconId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:          inputJson,
			Message:        messageJson,
			RoyaltyReconId: &royaltyReconId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailRoyaltyRecon, detailRoyaltyReconErr := h.royaltyReconService.GetDetailTransactionRoyaltyRecon(idInt, iupopkIdInt)

	if detailRoyaltyReconErr != nil {
		status := 400

		if detailRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReconErr.Error(),
		})
	}

	updateRoyaltyRecon, updateRoyaltyReconErr := h.historyService.UpdateDocumentRoyaltyRecon(idInt, *inputUpdateRoyaltyRecon, uint(claims["id"].(float64)), iupopkIdInt)

	if updateRoyaltyReconErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyRecon
		inputMap["royalty_recon_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateRoyaltyReconErr.Error(),
		})

		royaltyReconId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:          inputJson,
			Message:        messageJson,
			RoyaltyReconId: &royaltyReconId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateRoyaltyReconErr.Error(),
			"message": "failed to update royalty recon",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "royalty recon"
	inputNotification.Period = fmt.Sprintf("%v/%v", detailRoyaltyRecon.Detail.DateFrom, detailRoyaltyRecon.Detail.DateTo)
	inputNotification.Status = "membuat dokumen"
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyRecon
		inputMap["royalty_recon_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		royaltyReconId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:          inputJson,
			Message:        messageJson,
			RoyaltyReconId: &royaltyReconId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update royalty recon",
		})
	}

	return c.Status(200).JSON(updateRoyaltyRecon)
}

func (h *masterReportHandler) RequestCreateExcelRoyaltyRecon(c *fiber.Ctx) error {
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

	detailRoyaltyRecon, detailRoyaltyReconErr := h.royaltyReconService.GetDetailTransactionRoyaltyRecon(idInt, iupopkIdInt)

	if detailRoyaltyReconErr != nil {
		status := 400

		if detailRoyaltyReconErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReconErr.Error(),
		})
	}

	var inputRequestCreateExcel royaltyrecon.InputRequestCreateUploadRoyaltyRecon
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.RoyaltyRecon = detailRoyaltyRecon.Detail
	inputRequestCreateExcel.ListTransaction = detailRoyaltyRecon.ListTransaction
	inputRequestCreateExcel.Iupopk = detailRoyaltyRecon.Detail.Iupopk
	hitJob, hitJobErr := h.royaltyReconService.RequestCreateExcelRoyaltyRecon(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *masterReportHandler) ListRoyaltyRecon(c *fiber.Ctx) error {
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

	var filterRoyaltyRecon royaltyrecon.SortFilterRoyaltyRecon

	filterRoyaltyRecon.Field = c.Query("field")
	filterRoyaltyRecon.Sort = c.Query("sort")
	filterRoyaltyRecon.DateStart = c.Query("date_start")
	filterRoyaltyRecon.DateEnd = c.Query("date_end")

	listRoyaltyRecon, listRoyaltyReconErr := h.royaltyReconService.ListRoyaltyRecon(pageNumber, filterRoyaltyRecon, iupopkIdInt)

	if listRoyaltyReconErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listRoyaltyReconErr.Error(),
		})
	}

	return c.Status(200).JSON(listRoyaltyRecon)
}

// Report Royalty Report
func (h *masterReportHandler) PreviewRoyaltyReport(c *fiber.Ctx) error {
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

	royaltyReportInput := new(royaltyreport.InputRoyaltyReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(royaltyReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*royaltyReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	royaltyReport, royaltyReportErr := h.royaltyReportService.GetTransactionRoyaltyReport(royaltyReportInput.DateFrom, royaltyReportInput.DateTo, iupopkIdInt)

	if royaltyReportErr != nil {
		status := 400

		if royaltyReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": royaltyReportErr.Error(),
		})
	}

	return c.Status(200).JSON(royaltyReport)
}

func (h *masterReportHandler) CreateRoyaltyReport(c *fiber.Ctx) error {
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

	royaltyReportInput := new(royaltyreport.InputRoyaltyReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(royaltyReportInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*royaltyReportInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	royaltyReport, royaltyReportErr := h.historyService.CreateRoyaltyReport(royaltyReportInput.DateFrom, royaltyReportInput.DateTo, iupopkIdInt, uint(claims["id"].(float64)))

	if royaltyReportErr != nil {
		inputJson, _ := json.Marshal(royaltyReportInput)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error":   royaltyReportErr.Error(),
			"message": "create royalty report err",
		})

		createdErrLog := logs.Logs{
			Input:   inputJson,
			Message: messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": royaltyReportErr.Error(),
		})
	}

	return c.Status(201).JSON(royaltyReport)
}

func (h *masterReportHandler) DeleteRoyaltyReport(c *fiber.Ctx) error {
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

	detailRoyaltyReport, detailRoyaltyReportErr := h.royaltyReportService.GetDetailTransactionRoyaltyReport(idInt, iupopkIdInt)

	if detailRoyaltyReportErr != nil {
		status := 400

		if detailRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete royalty report",
			"error":   detailRoyaltyReportErr.Error(),
		})
	}

	_, isDeletedRoyaltyReportErr := h.historyService.DeleteRoyaltyReport(idInt, iupopkIdInt, uint(claims["id"].(float64)))

	if isDeletedRoyaltyReportErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["royalty_report_id"] = idInt

		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": isDeletedRoyaltyReportErr.Error(),
		})

		royaltyReportId := uint(idInt)
		createdErrLog := logs.Logs{
			RoyaltyReportId: &royaltyReportId,
			Input:           inputJson,
			Message:         messageJson,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400

		if isDeletedRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete royalty report",
			"error":   isDeletedRoyaltyReportErr.Error(),
		})
	}

	if detailRoyaltyReport.Detail.RoyaltyReportDocumentLink != "" {
		documentLink := detailRoyaltyReport.Detail.RoyaltyReportDocumentLink

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
			inputMap["royalty_report_id"] = idInt

			inputJson, _ := json.Marshal(inputMap)
			messageJson, _ := json.Marshal(map[string]interface{}{
				"error": isDeletedRoyaltyReportErr.Error(),
			})

			royaltyReportId := uint(idInt)
			createdErrLog := logs.Logs{
				RoyaltyReportId: &royaltyReportId,
				Input:           inputJson,
				Message:         messageJson,
			}

			h.logService.CreateLogs(createdErrLog)

			return c.Status(400).JSON(fiber.Map{
				"message": "failed to delete royalty report aws",
				"error":   deleteAwsErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success delete royalty report",
	})
}

func (h *masterReportHandler) DetailRoyaltyReport(c *fiber.Ctx) error {
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

	detailRoyaltyReport, detailRoyaltyReportErr := h.royaltyReportService.GetDetailTransactionRoyaltyReport(idInt, iupopkIdInt)

	if detailRoyaltyReportErr != nil {
		status := 400

		if detailRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReportErr.Error(),
		})
	}

	return c.Status(200).JSON(detailRoyaltyReport)
}

func (h *masterReportHandler) UpdateDocumentRoyaltyReport(c *fiber.Ctx) error {
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

	inputUpdateRoyaltyReport := new(royaltyreport.InputUpdateDocumentRoyaltyReport)

	// Binds the request body to the Person struct
	if err := c.BodyParser(inputUpdateRoyaltyReport); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to update royalty report",
		})
	}

	errors := h.v.Struct(*inputUpdateRoyaltyReport)

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "failed to update royalty report",
			"error":   "record not found",
		})
	}

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyReport
		inputMap["royalty_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"errors": dataErrors,
		})

		royaltyReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:           inputJson,
			Message:         messageJson,
			RoyaltyReportId: &royaltyReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	detailRoyaltyReport, detailRoyaltyReportErr := h.royaltyReportService.GetDetailTransactionRoyaltyReport(idInt, iupopkIdInt)

	if detailRoyaltyReportErr != nil {
		status := 400

		if detailRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReportErr.Error(),
		})
	}

	updateRoyaltyReport, updateRoyaltyReportErr := h.historyService.UpdateDocumentRoyaltyReport(idInt, *inputUpdateRoyaltyReport, uint(claims["id"].(float64)), iupopkIdInt)

	if updateRoyaltyReportErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyReport
		inputMap["royalty_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": updateRoyaltyReportErr.Error(),
		})

		royaltyReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:           inputJson,
			Message:         messageJson,
			RoyaltyReportId: &royaltyReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		status := 400
		if updateRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}

		return c.Status(status).JSON(fiber.Map{
			"error":   updateRoyaltyReportErr.Error(),
			"message": "failed to update royalty report",
		})
	}

	var inputNotification notification.InputNotification
	inputNotification.Type = "royalty report"
	inputNotification.Period = fmt.Sprintf("%v/%v", detailRoyaltyReport.Detail.DateFrom, detailRoyaltyReport.Detail.DateTo)
	inputNotification.Status = "membuat dokumen"
	_, createdNotificationErr := h.notificationUserService.CreateNotification(inputNotification, uint(claims["id"].(float64)), iupopkIdInt)

	if createdNotificationErr != nil {
		inputMap := make(map[string]interface{})
		inputMap["user_id"] = claims["id"]
		inputMap["input"] = inputUpdateRoyaltyReport
		inputMap["royalty_report_id"] = idInt
		inputJson, _ := json.Marshal(inputMap)
		messageJson, _ := json.Marshal(map[string]interface{}{
			"error": createdNotificationErr.Error(),
		})
		royaltyReportId := uint(idInt)

		createdErrLog := logs.Logs{
			Input:           inputJson,
			Message:         messageJson,
			RoyaltyReportId: &royaltyReportId,
		}

		h.logService.CreateLogs(createdErrLog)

		return c.Status(400).JSON(fiber.Map{
			"error":   createdNotificationErr.Error(),
			"message": "failed to create notification update royalty report",
		})
	}

	return c.Status(200).JSON(updateRoyaltyReport)
}

func (h *masterReportHandler) RequestCreateExcelRoyaltyReport(c *fiber.Ctx) error {
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

	detailRoyaltyReport, detailRoyaltyReportErr := h.royaltyReportService.GetDetailTransactionRoyaltyReport(idInt, iupopkIdInt)

	if detailRoyaltyReportErr != nil {
		status := 400

		if detailRoyaltyReportErr.Error() == "record not found" {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error": detailRoyaltyReportErr.Error(),
		})
	}

	var inputRequestCreateExcel royaltyreport.InputRequestCreateUploadRoyaltyReport
	inputRequestCreateExcel.Authorization = header["Authorization"]
	inputRequestCreateExcel.RoyaltyReport = detailRoyaltyReport.Detail
	inputRequestCreateExcel.ListTransaction = detailRoyaltyReport.ListTransaction
	inputRequestCreateExcel.Iupopk = detailRoyaltyReport.Detail.Iupopk
	hitJob, hitJobErr := h.royaltyReportService.RequestCreateExcelRoyaltyReport(inputRequestCreateExcel)

	if hitJobErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": hitJobErr.Error(),
		})
	}

	return c.Status(200).JSON(hitJob)
}

func (h *masterReportHandler) ListRoyaltyReport(c *fiber.Ctx) error {
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

	var filterRoyaltyReport royaltyreport.SortFilterRoyaltyReport

	filterRoyaltyReport.Field = c.Query("field")
	filterRoyaltyReport.Sort = c.Query("sort")
	filterRoyaltyReport.DateStart = c.Query("date_start")
	filterRoyaltyReport.DateEnd = c.Query("date_end")

	listRoyaltyReport, listRoyaltyReportErr := h.royaltyReportService.ListRoyaltyReport(pageNumber, filterRoyaltyReport, iupopkIdInt)

	if listRoyaltyReportErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listRoyaltyReportErr.Error(),
		})
	}

	return c.Status(200).JSON(listRoyaltyReport)
}
