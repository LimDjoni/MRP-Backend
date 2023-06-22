package handler

import (
	"ajebackend/model/awshelper"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
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
}

func NewMasterReportHandler(masterReportService masterreport.Service, userIupopkService useriupopk.Service, v *validator.Validate, allMasterService allmaster.Service, transactionRequestReportService transactionrequestreport.Service, logService logs.Service) *masterReportHandler {
	return &masterReportHandler{
		masterReportService,
		userIupopkService,
		v,
		allMasterService,
		transactionRequestReportService,
		logService,
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

		fmt.Println("finish 3")

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
