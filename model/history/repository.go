package history

import (
	"ajebackend/helper"
	"ajebackend/model/coareport"
	"ajebackend/model/counter"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/master/currency"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/navycompany"
	"ajebackend/model/master/navyship"
	"ajebackend/model/master/salessystem"
	"ajebackend/model/master/vessel"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/rkab"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint, transactionType string, iupopId int) (bool, error)
	UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error)
	CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	CreateMinerba(period string, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error)
	UpdateMinerba(id int, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error)
	DeleteMinerba(idMinerba int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint, iupopkId int) (minerba.Minerba, error)
	CreateDmo(dmoInput dmo.CreateDmoInput, userId uint, iupopkId int) (dmo.Dmo, error)
	DeleteDmo(idDmo int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string, iupopkId int) (dmo.Dmo, error)
	UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	CreateProduction(input production.InputCreateProduction, userId uint, iupopkId int) (production.Production, error)
	UpdateProduction(input production.InputCreateProduction, productionId int, userId uint, iupopkId int) (production.Production, error)
	DeleteProduction(productionId int, userId uint, iupopkId int) (bool, error)
	CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	DeleteGroupingVesselDn(id int, userId uint, iupopkId int) (bool, error)
	UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	DeleteGroupingVesselLn(id int, userId uint, iupopkId int) (bool, error)
	CreateMinerbaLn(period string, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	UpdateMinerbaLn(id int, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	DeleteMinerbaLn(idMinerbaLn int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	CreateInsw(month string, year int, userId uint, iupopkId int) (insw.Insw, error)
	DeleteInsw(idInsw int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint, iupopkId int) (insw.Insw, error)
	CreateReportDmo(input reportdmo.InputCreateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	DeleteReportDmo(idReportDmo int, userId uint, iupopkId int) (bool, error)
	CreateCoaReport(dateFrom string, dateTo string, iupopkId int, userId uint) (coareport.CoaReport, error)
	DeleteCoaReport(id int, iupopkId int, userId uint) (bool, error)
	UpdateDocumentCoaReport(id int, documentLink coareport.InputUpdateDocumentCoaReport, userId uint, iupopkId int) (coareport.CoaReport, error)
	CreateRkab(input rkab.RkabInput, iupopkId int, userId uint) (rkab.Rkab, error)
	DeleteRkab(id int, iupopkId int, userId uint) (bool, error)
	UploadDocumentRkab(id uint, urlS3 string, userId uint, iupopkId int) (rkab.Rkab, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func createIdNumber(model string, id uint) string {
	year, month, _ := time.Now().Date()

	monthNumber := strconv.Itoa(int(month))

	if len([]rune(monthNumber)) < 2 {
		monthNumber = "0" + monthNumber
	}

	idNumber := fmt.Sprintf("%s-%v-%v-%v", model, monthNumber, year%1e2, helper.CreateIdNumber(int(id)))

	return idNumber
}

// Transaction

func (r *repository) CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	var createdTransaction transaction.Transaction

	tx := r.db.Begin()

	var curr currency.Currency

	currencyErr := tx.Where("code = ?", "IDR").First(&curr).Error

	if currencyErr != nil {
		tx.Rollback()
		return createdTransaction, currencyErr
	}

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdTransaction, findIupErr
	}

	createdTransaction.DmoId = nil
	createdTransaction.TransactionType = "DN"
	createdTransaction.SellerId = &iup.ID
	createdTransaction.DestinationCountryId = inputTransactionDN.DestinationCountryId
	createdTransaction.ShippingDate = inputTransactionDN.ShippingDate

	createdTransaction.Quantity = math.Round(inputTransactionDN.Quantity*1000) / 1000
	createdTransaction.QuantityUnloading = math.Round(inputTransactionDN.QuantityUnloading*1000) / 1000
	createdTransaction.TugboatId = inputTransactionDN.TugboatId
	createdTransaction.BargeId = inputTransactionDN.BargeId
	createdTransaction.VesselId = inputTransactionDN.VesselId
	createdTransaction.CustomerId = inputTransactionDN.CustomerId
	createdTransaction.LoadingPortId = inputTransactionDN.LoadingPortId
	createdTransaction.UnloadingPortId = inputTransactionDN.UnloadingPortId
	createdTransaction.DmoDestinationPortId = inputTransactionDN.DmoDestinationPortId
	createdTransaction.SkbDate = inputTransactionDN.SkbDate
	createdTransaction.SkbNumber = strings.ToUpper(inputTransactionDN.SkbNumber)
	createdTransaction.SkabDate = inputTransactionDN.SkabDate
	createdTransaction.SkabNumber = strings.ToUpper(inputTransactionDN.SkabNumber)
	createdTransaction.BillOfLadingDate = inputTransactionDN.BillOfLadingDate
	createdTransaction.BillOfLadingNumber = strings.ToUpper(inputTransactionDN.BillOfLadingNumber)
	createdTransaction.RoyaltyRate = inputTransactionDN.RoyaltyRate
	createdTransaction.DpRoyaltyPrice = inputTransactionDN.DpRoyaltyPrice
	createdTransaction.DpRoyaltyCurrencyId = &curr.ID
	createdTransaction.DpRoyaltyDate = inputTransactionDN.DpRoyaltyDate
	if inputTransactionDN.DpRoyaltyNtpn != nil {
		dpNtpn := strings.ToUpper(*inputTransactionDN.DpRoyaltyNtpn)
		createdTransaction.DpRoyaltyNtpn = &dpNtpn
	}

	if inputTransactionDN.DpRoyaltyBillingCode != nil {
		dpBillingCode := strings.ToUpper(*inputTransactionDN.DpRoyaltyBillingCode)
		createdTransaction.DpRoyaltyBillingCode = &dpBillingCode
	}

	createdTransaction.DpRoyaltyTotal = inputTransactionDN.DpRoyaltyTotal
	createdTransaction.PaymentDpRoyaltyPrice = inputTransactionDN.PaymentDpRoyaltyPrice
	createdTransaction.PaymentDpRoyaltyCurrencyId = &curr.ID
	createdTransaction.PaymentDpRoyaltyDate = inputTransactionDN.PaymentDpRoyaltyDate
	if inputTransactionDN.PaymentDpRoyaltyNtpn != nil {
		paymentDpNtpn := strings.ToUpper(*inputTransactionDN.PaymentDpRoyaltyNtpn)
		createdTransaction.PaymentDpRoyaltyNtpn = &paymentDpNtpn
	}

	if inputTransactionDN.PaymentDpRoyaltyBillingCode != nil {
		paymentDpBillingCode := strings.ToUpper(*inputTransactionDN.PaymentDpRoyaltyBillingCode)
		createdTransaction.PaymentDpRoyaltyBillingCode = &paymentDpBillingCode
	}

	createdTransaction.PaymentDpRoyaltyTotal = inputTransactionDN.PaymentDpRoyaltyTotal
	createdTransaction.LhvDate = inputTransactionDN.LhvDate
	createdTransaction.LhvNumber = strings.ToUpper(inputTransactionDN.LhvNumber)
	createdTransaction.SurveyorId = inputTransactionDN.SurveyorId
	createdTransaction.CowDate = inputTransactionDN.CowDate
	createdTransaction.CowNumber = strings.ToUpper(inputTransactionDN.CowNumber)
	createdTransaction.CoaDate = inputTransactionDN.CoaDate
	createdTransaction.CoaNumber = strings.ToUpper(inputTransactionDN.CoaNumber)
	createdTransaction.QualityTmAr = inputTransactionDN.QualityTmAr
	createdTransaction.QualityImAdb = inputTransactionDN.QualityImAdb
	createdTransaction.QualityAshAr = inputTransactionDN.QualityAshAr
	createdTransaction.QualityAshAdb = inputTransactionDN.QualityAshAdb
	createdTransaction.QualityVmAdb = inputTransactionDN.QualityVmAdb
	createdTransaction.QualityFcAdb = inputTransactionDN.QualityFcAdb
	createdTransaction.QualityTsAr = inputTransactionDN.QualityTsAr
	createdTransaction.QualityTsAdb = inputTransactionDN.QualityTsAdb
	createdTransaction.QualityCaloriesAr = inputTransactionDN.QualityCaloriesAr
	createdTransaction.QualityCaloriesAdb = inputTransactionDN.QualityCaloriesAdb
	createdTransaction.BargingDistance = inputTransactionDN.BargingDistance
	createdTransaction.SalesSystemId = inputTransactionDN.SalesSystemId
	createdTransaction.InvoiceDate = inputTransactionDN.InvoiceDate
	createdTransaction.InvoiceNumber = strings.ToUpper(inputTransactionDN.InvoiceNumber)
	createdTransaction.InvoicePriceUnit = inputTransactionDN.InvoicePriceUnit
	createdTransaction.InvoicePriceTotal = inputTransactionDN.InvoicePriceTotal
	createdTransaction.ContractDate = inputTransactionDN.ContractDate
	createdTransaction.ContractNumber = strings.ToUpper(inputTransactionDN.ContractNumber)
	createdTransaction.DmoBuyerId = inputTransactionDN.DmoBuyerId
	createdTransaction.IsCoaFinish = inputTransactionDN.IsCoaFinish
	createdTransaction.IsRoyaltyFinalFinish = inputTransactionDN.IsRoyaltyFinalFinish
	createdTransaction.DestinationId = inputTransactionDN.DestinationId
	createTransactionErr := tx.Create(&createdTransaction).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, createTransactionErr
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, findCounterTransactionErr
	}

	code := "TDN-"

	code += iup.Code

	idNumber := createIdNumber(code, uint(counterTransaction.TransactionDn))

	updateTransactionsErr := tx.Model(&createdTransaction).Where("id = ?", createdTransaction.ID).Update("id_number", idNumber).Error

	if updateTransactionsErr != nil {
		tx.Rollback()
		return createdTransaction, updateTransactionsErr
	}

	var history History

	history.TransactionId = &createdTransaction.ID
	history.Status = "Created"
	history.UserId = userId
	history.IupopkId = createdTransaction.SellerId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdTransaction, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("transaction_dn", counterTransaction.TransactionDn+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdTransaction, updateCounterErr
	}

	tx.Commit()
	return createdTransaction, nil
}

func (r *repository) DeleteTransaction(id int, userId uint, transactionType string, iupopId int) (bool, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ? AND transaction_type = ? AND seller_id = ?", id, transactionType, iupopId).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := tx.Unscoped().Where("id = ? AND transaction_type = ? AND seller_id = ?", id, transactionType, iupopId).Delete(&transaction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Transaction %v with id number %s and id %v and iupop %v", transactionType, *transaction.IdNumber, transaction.ID, iupopId)
	history.UserId = userId
	history.IupopkId = transaction.SellerId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, createHistoryErr
}

func (r *repository) UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ? AND transaction_type = ? AND seller_id = ?", idTransaction, "DN", iupopkId).First(&transaction).Error

	if errFind != nil {
		tx.Rollback()
		return transaction, errFind
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorBeforeDataJsonMarshal
	}
	inputEditTransactionDN.SkbNumber = strings.ToUpper(inputEditTransactionDN.SkbNumber)
	inputEditTransactionDN.SkabNumber = strings.ToUpper(inputEditTransactionDN.SkabNumber)
	inputEditTransactionDN.BillOfLadingNumber = strings.ToUpper(inputEditTransactionDN.BillOfLadingNumber)
	inputEditTransactionDN.LhvNumber = strings.ToUpper(inputEditTransactionDN.LhvNumber)
	inputEditTransactionDN.CowNumber = strings.ToUpper(inputEditTransactionDN.CowNumber)
	inputEditTransactionDN.CoaNumber = strings.ToUpper(inputEditTransactionDN.CoaNumber)
	inputEditTransactionDN.InvoiceNumber = strings.ToUpper(inputEditTransactionDN.InvoiceNumber)
	inputEditTransactionDN.ContractNumber = strings.ToUpper(inputEditTransactionDN.ContractNumber)
	inputEditTransactionDN.Quantity = math.Round(inputEditTransactionDN.Quantity*1000) / 1000
	inputEditTransactionDN.QuantityUnloading = math.Round(inputEditTransactionDN.QuantityUnloading*1000) / 1000
	dataInput, errorMarshal := json.Marshal(inputEditTransactionDN)

	if errorMarshal != nil {
		tx.Rollback()
		return transaction, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	delete(dataInputMapString, "vessel_name")
	delete(dataInputMapString, "dmo_destination_port_ln_name")
	if errorUnmarshal != nil {
		tx.Rollback()
		return transaction, errorUnmarshal
	}

	updateErr := tx.Model(&transaction).Updates(dataInputMapString).Error

	if updateErr != nil {
		tx.Rollback()
		return transaction, updateErr
	}

	afterData, errorAfterDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorAfterDataJsonMarshal
	}

	var history History

	history.TransactionId = &transaction.ID
	history.UserId = userId
	history.Status = "Updated"
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = transaction.SellerId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return transaction, createHistoryErr
	}

	tx.Commit()
	return transaction, nil
}

func (r *repository) UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error) {
	var uploadedTransaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ? AND transaction_type = ? AND seller_id = ?", idTransaction, transactionType, iupopkId).First(&uploadedTransaction).Error

	if errFind != nil {
		return uploadedTransaction, errFind
	}

	var isReupload = false
	editData := make(map[string]interface{})

	switch documentType {
	case "skb":
		if uploadedTransaction.SkbDocumentLink != "" {
			isReupload = true
		}
		editData["skb_document_link"] = urlS3
	case "skab":
		if uploadedTransaction.SkabDocumentLink != "" {
			isReupload = true
		}
		editData["skab_document_link"] = urlS3
	case "bl":
		if uploadedTransaction.BLDocumentLink != "" {
			isReupload = true
		}
		editData["bl_document_link"] = urlS3
	case "royalti_provision":
		if uploadedTransaction.RoyaltiProvisionDocumentLink != "" {
			isReupload = true
		}
		editData["royalti_provision_document_link"] = urlS3
	case "royalti_final":
		if uploadedTransaction.RoyaltiFinalDocumentLink != "" {
			isReupload = true
		}
		editData["royalti_final_document_link"] = urlS3
	case "cow":
		if uploadedTransaction.COWDocumentLink != "" {
			isReupload = true
		}
		editData["cow_document_link"] = urlS3
	case "coa":
		if uploadedTransaction.COADocumentLink != "" {
			isReupload = true
		}
		editData["coa_document_link"] = urlS3
	case "invoice":
		if uploadedTransaction.InvoiceAndContractDocumentLink != "" {
			isReupload = true
		}
		editData["invoice_and_contract_document_link"] = urlS3
	case "lhv":
		if uploadedTransaction.LHVDocumentLink != "" {
			isReupload = true
		}
		editData["lhv_document_link"] = urlS3
	}

	errEdit := tx.Model(&uploadedTransaction).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return uploadedTransaction, errEdit
	}
	var history History

	history.TransactionId = &uploadedTransaction.ID
	history.UserId = userId
	if isReupload == false {
		history.Status = fmt.Sprintf("Uploaded %s document", documentType)
	}

	if isReupload == true {
		history.Status = fmt.Sprintf("Reupload %s document", documentType)
	}

	history.IupopkId = uploadedTransaction.SellerId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return uploadedTransaction, createHistoryErr
	}

	tx.Commit()
	return uploadedTransaction, nil
}

// Transaction LN

func (r *repository) CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	var createdTransactionLn transaction.Transaction

	tx := r.db.Begin()

	var curr currency.Currency

	currencyErr := tx.Where("code = ?", "USD").First(&curr).Error

	if currencyErr != nil {
		tx.Rollback()
		return createdTransactionLn, currencyErr
	}

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdTransactionLn, findIupErr
	}

	var createVesselMaster vessel.Vessel
	if inputTransactionLN.VesselName != "" {

		createVesselMaster.Name = inputTransactionLN.VesselName

		errCreateVesselMaster := tx.FirstOrCreate(&createVesselMaster, createVesselMaster).Error

		if errCreateVesselMaster != nil {
			tx.Rollback()
			return createdTransactionLn, errCreateVesselMaster
		}
	} else {
		return createdTransactionLn, errors.New("vessel name is required")
	}

	createdTransactionLn.TransactionType = "LN"
	createdTransactionLn.SellerId = &iup.ID
	createdTransactionLn.Quantity = math.Round(inputTransactionLN.Quantity*1000) / 1000
	createdTransactionLn.QuantityUnloading = math.Round(inputTransactionLN.QuantityUnloading*1000) / 1000
	createdTransactionLn.DestinationCountryId = inputTransactionLN.DestinationCountryId
	createdTransactionLn.ShippingDate = inputTransactionLN.ShippingDate
	createdTransactionLn.TugboatId = inputTransactionLN.TugboatId
	createdTransactionLn.BargeId = inputTransactionLN.BargeId
	createdTransactionLn.VesselId = &createVesselMaster.ID
	createdTransactionLn.CustomerId = inputTransactionLN.CustomerId
	createdTransactionLn.LoadingPortId = inputTransactionLN.LoadingPortId
	createdTransactionLn.UnloadingPortId = inputTransactionLN.UnloadingPortId
	createdTransactionLn.DmoDestinationPortId = inputTransactionLN.DmoDestinationPortId
	createdTransactionLn.SkbDate = inputTransactionLN.SkbDate
	createdTransactionLn.SkbNumber = strings.ToUpper(inputTransactionLN.SkbNumber)
	createdTransactionLn.SkabDate = inputTransactionLN.SkabDate
	createdTransactionLn.SkabNumber = strings.ToUpper(inputTransactionLN.SkabNumber)
	createdTransactionLn.BillOfLadingDate = inputTransactionLN.BillOfLadingDate
	createdTransactionLn.BillOfLadingNumber = strings.ToUpper(inputTransactionLN.BillOfLadingNumber)
	createdTransactionLn.RoyaltyRate = inputTransactionLN.RoyaltyRate
	createdTransactionLn.DpRoyaltyPrice = inputTransactionLN.DpRoyaltyPrice
	createdTransactionLn.DpRoyaltyCurrencyId = &curr.ID
	createdTransactionLn.DpRoyaltyDate = inputTransactionLN.DpRoyaltyDate
	if inputTransactionLN.DpRoyaltyNtpn != nil {
		dpNtpn := strings.ToUpper(*inputTransactionLN.DpRoyaltyNtpn)
		createdTransactionLn.DpRoyaltyNtpn = &dpNtpn
	}

	if inputTransactionLN.DpRoyaltyBillingCode != nil {
		dpBillingCode := strings.ToUpper(*inputTransactionLN.DpRoyaltyBillingCode)
		createdTransactionLn.DpRoyaltyBillingCode = &dpBillingCode
	}

	createdTransactionLn.DpRoyaltyTotal = inputTransactionLN.DpRoyaltyTotal
	createdTransactionLn.PaymentDpRoyaltyPrice = inputTransactionLN.PaymentDpRoyaltyPrice
	createdTransactionLn.PaymentDpRoyaltyCurrencyId = &curr.ID

	createdTransactionLn.PaymentDpRoyaltyDate = inputTransactionLN.PaymentDpRoyaltyDate
	if inputTransactionLN.PaymentDpRoyaltyNtpn != nil {
		paymentDpNtpn := strings.ToUpper(*inputTransactionLN.PaymentDpRoyaltyNtpn)
		createdTransactionLn.PaymentDpRoyaltyNtpn = &paymentDpNtpn
	}

	if inputTransactionLN.PaymentDpRoyaltyBillingCode != nil {
		paymentDpBillingCode := strings.ToUpper(*inputTransactionLN.PaymentDpRoyaltyBillingCode)
		createdTransactionLn.PaymentDpRoyaltyBillingCode = &paymentDpBillingCode
	}

	createdTransactionLn.PaymentDpRoyaltyTotal = inputTransactionLN.PaymentDpRoyaltyTotal
	createdTransactionLn.LhvDate = inputTransactionLN.LhvDate
	createdTransactionLn.LhvNumber = strings.ToUpper(inputTransactionLN.LhvNumber)
	createdTransactionLn.SurveyorId = inputTransactionLN.SurveyorId
	createdTransactionLn.CowDate = inputTransactionLN.CowDate
	createdTransactionLn.CowNumber = strings.ToUpper(inputTransactionLN.CowNumber)
	createdTransactionLn.CoaDate = inputTransactionLN.CoaDate
	createdTransactionLn.CoaNumber = strings.ToUpper(inputTransactionLN.CoaNumber)
	createdTransactionLn.QualityTmAr = inputTransactionLN.QualityTmAr
	createdTransactionLn.QualityImAdb = inputTransactionLN.QualityImAdb
	createdTransactionLn.QualityAshAr = inputTransactionLN.QualityAshAr
	createdTransactionLn.QualityAshAdb = inputTransactionLN.QualityAshAdb
	createdTransactionLn.QualityVmAdb = inputTransactionLN.QualityVmAdb
	createdTransactionLn.QualityFcAdb = inputTransactionLN.QualityFcAdb
	createdTransactionLn.QualityTsAr = inputTransactionLN.QualityTsAr
	createdTransactionLn.QualityTsAdb = inputTransactionLN.QualityTsAdb
	createdTransactionLn.QualityCaloriesAr = inputTransactionLN.QualityCaloriesAr
	createdTransactionLn.QualityCaloriesAdb = inputTransactionLN.QualityCaloriesAdb
	createdTransactionLn.BargingDistance = inputTransactionLN.BargingDistance
	createdTransactionLn.SalesSystemId = inputTransactionLN.SalesSystemId
	createdTransactionLn.InvoiceDate = inputTransactionLN.InvoiceDate
	createdTransactionLn.InvoiceNumber = strings.ToUpper(inputTransactionLN.InvoiceNumber)
	createdTransactionLn.InvoicePriceUnit = inputTransactionLN.InvoicePriceUnit
	createdTransactionLn.InvoicePriceTotal = inputTransactionLN.InvoicePriceTotal
	createdTransactionLn.ContractDate = inputTransactionLN.ContractDate
	createdTransactionLn.ContractNumber = strings.ToUpper(inputTransactionLN.ContractNumber)
	createdTransactionLn.DmoBuyerId = inputTransactionLN.DmoBuyerId
	createdTransactionLn.DmoDestinationPortLnName = inputTransactionLN.DmoDestinationPortLnName

	createTransactionErr := tx.Create(&createdTransactionLn).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransactionLn, createTransactionErr
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdTransactionLn, findCounterTransactionErr
	}

	code := "TLN-"

	code += iup.Code

	idNumber := createIdNumber(code, uint(counterTransaction.TransactionLn))

	updateTransactionsErr := tx.Model(&createdTransactionLn).Update("id_number", idNumber).Error

	if updateTransactionsErr != nil {
		tx.Rollback()
		return createdTransactionLn, updateTransactionsErr
	}

	var history History

	history.TransactionId = &createdTransactionLn.ID
	history.Status = "Created LN"
	history.UserId = userId
	history.IupopkId = createdTransactionLn.SellerId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdTransactionLn, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("transaction_ln", counterTransaction.TransactionLn+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdTransactionLn, updateCounterErr
	}

	tx.Commit()
	return createdTransactionLn, createHistoryErr
}

func (r *repository) UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	var updatedTransactionLn transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ? AND transaction_type = ? AND seller_id = ?", id, "LN", iupopkId).First(&updatedTransactionLn).Error

	if errFind != nil {
		tx.Rollback()
		return updatedTransactionLn, errFind
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(updatedTransactionLn)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedTransactionLn, errorBeforeDataJsonMarshal
	}

	inputTransactionLN.SkbNumber = strings.ToUpper(inputTransactionLN.SkbNumber)
	inputTransactionLN.SkabNumber = strings.ToUpper(inputTransactionLN.SkabNumber)
	inputTransactionLN.BillOfLadingNumber = strings.ToUpper(inputTransactionLN.BillOfLadingNumber)
	inputTransactionLN.LhvNumber = strings.ToUpper(inputTransactionLN.LhvNumber)
	inputTransactionLN.CowNumber = strings.ToUpper(inputTransactionLN.CowNumber)
	inputTransactionLN.CoaNumber = strings.ToUpper(inputTransactionLN.CoaNumber)
	inputTransactionLN.InvoiceNumber = strings.ToUpper(inputTransactionLN.InvoiceNumber)
	inputTransactionLN.ContractNumber = strings.ToUpper(inputTransactionLN.ContractNumber)
	inputTransactionLN.Quantity = math.Round(inputTransactionLN.Quantity*1000) / 1000
	inputTransactionLN.QuantityUnloading = math.Round(inputTransactionLN.QuantityUnloading*1000) / 1000

	editTransaction, errorMarshal := json.Marshal(inputTransactionLN)

	if errorMarshal != nil {
		tx.Rollback()
		return updatedTransactionLn, errorMarshal
	}

	var editTransactionInput map[string]interface{}

	errorUnmarshalTransaction := json.Unmarshal(editTransaction, &editTransactionInput)

	if inputTransactionLN.VesselName != "" {
		var createVesselMaster vessel.Vessel

		createVesselMaster.Name = inputTransactionLN.VesselName

		errCreateVesselMaster := tx.FirstOrCreate(&createVesselMaster, createVesselMaster).Error

		if errCreateVesselMaster != nil {
			tx.Rollback()
			return updatedTransactionLn, errCreateVesselMaster
		}

		editTransactionInput["vessel_id"] = createVesselMaster.ID
	}

	delete(editTransactionInput, "vessel_name")
	if errorUnmarshalTransaction != nil {
		tx.Rollback()
		return updatedTransactionLn, errorUnmarshalTransaction
	}

	updateTransactionErr := tx.Preload(clause.Associations).Model(&updatedTransactionLn).Updates(editTransactionInput).Error

	if updateTransactionErr != nil {
		tx.Rollback()
		return updatedTransactionLn, updateTransactionErr
	}

	afterData, errorAfterDataJsonMarshal := json.Marshal(updatedTransactionLn)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedTransactionLn, errorAfterDataJsonMarshal
	}

	var history History

	history.TransactionId = &updatedTransactionLn.ID
	history.Status = "Updated LN"
	history.UserId = userId
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = updatedTransactionLn.SellerId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedTransactionLn, createHistoryErr
	}

	tx.Commit()
	return updatedTransactionLn, createHistoryErr
}

// Minerba

func (r *repository) CreateMinerba(period string, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error) {
	var createdMinerba minerba.Minerba

	tx := r.db.Begin()
	createdMinerba.Period = period
	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ? AND seller_id = ?", updateTransaction, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdMinerba, findTransactionsErr
	}

	if len(transactions) != len(updateTransaction) {
		tx.Rollback()
		return createdMinerba, errors.New("please check some of transactions not found")
	}

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdMinerba, findIupErr
	}

	var tempQuantity float64
	for _, v := range transactions {
		tempQuantity += v.QuantityUnloading
	}

	stringTempQuantity := fmt.Sprintf("%.3f", tempQuantity)
	parseTempQuantity, _ := strconv.ParseFloat(stringTempQuantity, 64)

	createdMinerba.Quantity = math.Round(parseTempQuantity*1000) / 1000

	createdMinerba.IupopkId = iup.ID

	errCreateMinerba := tx.Create(&createdMinerba).Error

	if errCreateMinerba != nil {
		tx.Rollback()
		return createdMinerba, errCreateMinerba
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdMinerba, findCounterTransactionErr
	}

	periodSplit := strings.Split(period, " ")

	idNumber := fmt.Sprintf("LSD-%s-%s-%s-", iup.Code, helper.MonthStringToNumberString(periodSplit[0]), periodSplit[1][len(periodSplit[1])-2:])
	idNumber += helper.CreateIdNumber(counterTransaction.Sp3medn)

	updateMinerbaErr := tx.Model(&createdMinerba).Update("id_number", idNumber).Error

	if updateMinerbaErr != nil {
		tx.Rollback()
		return createdMinerba, updateMinerbaErr
	}

	updateTransactionErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", updateTransaction, iupopkId).Update("minerba_id", createdMinerba.ID).Error

	if updateTransactionErr != nil {
		tx.Rollback()
		return createdMinerba, updateTransactionErr
	}

	var history History

	history.MinerbaId = &createdMinerba.ID
	history.Status = "Created Minerba Report"
	history.UserId = userId
	history.IupopkId = &createdMinerba.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdMinerba, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("sp3medn", counterTransaction.Sp3medn+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdMinerba, updateCounterErr
	}

	tx.Commit()
	return createdMinerba, nil
}

func (r *repository) UpdateMinerba(id int, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error) {
	var updatedMinerba minerba.Minerba
	var quantityMinerba float64

	historyBefore := make(map[string]interface{})
	historyAfter := make(map[string]interface{})
	tx := r.db.Begin()

	findMinerbaErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&updatedMinerba).Error

	if findMinerbaErr != nil {
		return updatedMinerba, findMinerbaErr
	}

	historyBefore["minerba"] = updatedMinerba

	var beforeTransaction []transaction.Transaction
	findTransactionBeforeErr := tx.Where("minerba_id = ? AND seller_id = ?", id, iupopkId).Find(&beforeTransaction).Error

	if findTransactionBeforeErr != nil {
		return updatedMinerba, findTransactionBeforeErr
	}

	var transactionBefore []uint

	for _, v := range beforeTransaction {
		transactionBefore = append(transactionBefore, v.ID)
	}

	historyBefore["transactions"] = transactionBefore

	errUpdMinerbaNil := tx.Model(&beforeTransaction).Where("minerba_id = ? AND seller_id = ?", id, iupopkId).Update("minerba_id", nil).Error

	if errUpdMinerbaNil != nil {
		tx.Rollback()
		return updatedMinerba, errUpdMinerbaNil
	}

	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ? AND seller_id = ?", updateTransaction, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return updatedMinerba, findTransactionsErr
	}

	if len(transactions) != len(updateTransaction) {
		tx.Rollback()
		return updatedMinerba, errors.New("please check some of transactions not found")
	}

	for _, v := range transactions {
		quantityMinerba += v.QuantityUnloading
	}

	quantityMinerba = math.Round(quantityMinerba*1000) / 1000

	errUpdateMinerba := tx.Model(&updatedMinerba).Update("quantity", quantityMinerba).Error

	if errUpdateMinerba != nil {
		tx.Rollback()
		return updatedMinerba, errUpdateMinerba
	}

	historyAfter["minerba"] = updatedMinerba
	historyAfter["transactions"] = updateTransaction

	updateTransactionErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", updateTransaction, iupopkId).Update("minerba_id", id).Error

	if updateTransactionErr != nil {
		tx.Rollback()
		return updatedMinerba, updateTransactionErr
	}

	var history History
	beforeData, errorBeforeDataJsonMarshal := json.Marshal(historyBefore)
	afterData, errorAfterDataJsonMarshal := json.Marshal(historyAfter)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedMinerba, errorBeforeDataJsonMarshal
	}

	if errorAfterDataJsonMarshal != nil {
		tx.Rollback()
		return updatedMinerba, errorAfterDataJsonMarshal
	}

	iupId := uint(iupopkId)
	history.MinerbaId = &updatedMinerba.ID
	history.Status = "Updated Minerba Report"
	history.UserId = userId
	history.AfterData = afterData
	history.BeforeData = beforeData
	history.IupopkId = &iupId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedMinerba, createHistoryErr
	}

	tx.Commit()
	return updatedMinerba, nil
}

func (r *repository) DeleteMinerba(idMinerba int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var minerba minerba.Minerba

	findMinerbaErr := tx.Where("id = ? AND iupopk_id = ?", idMinerba, iupopkId).First(&minerba).Error

	if findMinerbaErr != nil {
		tx.Rollback()
		return false, findMinerbaErr
	}

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", idMinerba, iupopkId).Delete(&minerba).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Minerba Report with id number %s and id %v", *minerba.IdNumber, minerba.ID)
	history.UserId = userId
	history.IupopkId = &minerba.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint, iupopkId int) (minerba.Minerba, error) {
	tx := r.db.Begin()
	var minerba minerba.Minerba

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&minerba).Error

	if errFind != nil {
		tx.Rollback()
		return minerba, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "sp3medn") {
				editData["sp3_medn_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&minerba).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return minerba, errEdit
	}

	var history History

	history.MinerbaId = &minerba.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document minerba with id = %v", minerba.ID)
	history.IupopkId = &minerba.IupopkId
	dataInput, _ := json.Marshal(documentLink)

	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return minerba, createHistoryErr
	}

	tx.Commit()
	return minerba, nil
}

// Dmo

func (r *repository) CreateDmo(dmoInput dmo.CreateDmoInput, userId uint, iupopkId int) (dmo.Dmo, error) {
	var createdDmo dmo.Dmo
	var transactionBarge []transaction.Transaction
	var groupingVessel []groupingvesseldn.GroupingVesselDn

	barge := false
	vessel := false
	tx := r.db.Begin()

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdDmo, findIupErr
	}

	createdDmo.IupopkId = iup.ID
	createdDmo.Period = dmoInput.Period
	createdDmo.IsDocumentCustom = dmoInput.IsDocumentCustom
	createdDmo.DocumentDate = dmoInput.DocumentDate
	if len(dmoInput.TransactionBarge) > 0 {
		var bargeQuantity float64
		barge = true
		findTransactionBargeErr := tx.Where("id IN ? AND dmo_id IS NULL AND seller_id = ?", dmoInput.TransactionBarge, iupopkId).Find(&transactionBarge).Error

		if findTransactionBargeErr != nil {
			tx.Rollback()
			return createdDmo, findTransactionBargeErr
		}

		if len(transactionBarge) != len(dmoInput.TransactionBarge) {
			tx.Rollback()
			return createdDmo, errors.New("Ada transaksi yang sudah digunakan")
		}

		for _, v := range transactionBarge {
			bargeQuantity += v.Quantity
		}

		createdDmo.BargeTotalQuantity = math.Round(bargeQuantity*1000) / 1000
		createdDmo.BargeGrandTotalQuantity = math.Round(bargeQuantity*1000) / 1000
	}

	if len(dmoInput.GroupingVessel) > 0 {
		var vesselQuantity float64
		var vesselAdjustment float64
		var vesselGrandTotalQuantity float64
		vessel = true
		var checkDmoGrouping []dmovessel.DmoVessel

		findCheckDmoGroupingErr := tx.Where("grouping_vessel_dn_id IN ?", dmoInput.GroupingVessel).Find(&checkDmoGrouping).Error

		if findCheckDmoGroupingErr != nil {
			tx.Rollback()
			return createdDmo, findCheckDmoGroupingErr
		}

		if len(checkDmoGrouping) > 0 {
			tx.Rollback()
			return createdDmo, errors.New("Ada grouping vessel yang sudah digunakan")
		}

		findGroupingVesselErr := tx.Where("id IN ? AND iupopk_id = ?", dmoInput.GroupingVessel, iupopkId).Find(&groupingVessel).Error

		if findGroupingVesselErr != nil {
			tx.Rollback()
			return createdDmo, findGroupingVesselErr
		}

		for _, v := range groupingVessel {
			vesselQuantity += v.Quantity
			vesselAdjustment += v.Adjustment
			vesselGrandTotalQuantity += v.GrandTotalQuantity
		}

		createdDmo.VesselAdjustment = math.Round(vesselAdjustment*1000) / 1000
		createdDmo.VesselTotalQuantity = math.Round(vesselQuantity*1000) / 1000
		createdDmo.VesselGrandTotalQuantity = math.Round(vesselGrandTotalQuantity*1000) / 1000
	}

	if len(transactionBarge) != len(dmoInput.TransactionBarge) && len(groupingVessel) != len(dmoInput.GroupingVessel) {
		tx.Rollback()
		return createdDmo, errors.New("please check some of transactions not found")
	}

	if barge && vessel {
		createdDmo.Type = "Combination"
	}

	if barge && !vessel {
		createdDmo.Type = "Barge"
	}

	if !barge && vessel {
		createdDmo.Type = "Vessel"
	}

	if dmoInput.IsDocumentCustom {
		createdDmo.IsReconciliationLetterDownloaded = true
		createdDmo.IsReconciliationLetterSigned = true
	}

	createdDmoErr := tx.Preload(clause.Associations).Create(&createdDmo).Error

	if createdDmoErr != nil {
		tx.Rollback()
		return createdDmo, createdDmoErr
	}

	findDmoErr := tx.Preload(clause.Associations).Where("id = ?", createdDmo.ID).First(&createdDmo).Error

	if findDmoErr != nil {
		tx.Rollback()
		return createdDmo, findDmoErr
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdDmo, findCounterTransactionErr
	}
	code := "LBU-"

	code += iup.Code

	periodSplit := strings.Split(dmoInput.Period, " ")

	idNumber := fmt.Sprintf("LBU-%s-%s-%s-", iup.Code, helper.MonthStringToNumberString(periodSplit[0]), periodSplit[1][len(periodSplit[1])-2:])
	idNumber += helper.CreateIdNumber(counterTransaction.BaEndUser)

	updateDmoErr := tx.Model(&createdDmo).Update("id_number", idNumber).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return createdDmo, updateDmoErr
	}

	if len(dmoInput.TransactionBarge) > 0 {
		updateTransactionBargeErr := tx.Table("transactions").Where("id IN ?", dmoInput.TransactionBarge).Update("dmo_id", createdDmo.ID).Error

		if updateTransactionBargeErr != nil {
			tx.Rollback()
			return createdDmo, updateTransactionBargeErr
		}
	}

	if len(dmoInput.GroupingVessel) > 0 {
		var transactionGroupVessel []transaction.Transaction

		var listIdTransactionGroupVessel []uint
		findTransactionGroupVesselErr := tx.Where("grouping_vessel_dn_id IN ? AND dmo_id IS NULL", dmoInput.GroupingVessel).Find(&transactionGroupVessel).Error

		if findTransactionGroupVesselErr != nil {
			tx.Rollback()
			return createdDmo, findTransactionGroupVesselErr
		}

		for _, v := range transactionGroupVessel {
			listIdTransactionGroupVessel = append(listIdTransactionGroupVessel, v.ID)
		}

		updateTransactionGroupVesselErr := tx.Table("transactions").Where("id IN ?", listIdTransactionGroupVessel).Update("dmo_id", createdDmo.ID).Error

		if updateTransactionGroupVesselErr != nil {
			tx.Rollback()
			return createdDmo, updateTransactionGroupVesselErr
		}

		var dmoVessels []dmovessel.DmoVessel

		if len(dmoInput.GroupingVessel) > 0 {

			for _, v := range dmoInput.GroupingVessel {
				var dmoVesselDummy dmovessel.DmoVessel
				dmoVesselDummy.DmoId = createdDmo.ID
				dmoVesselDummy.GroupingVesselDnId = uint(v)
				dmoVessels = append(dmoVessels, dmoVesselDummy)
			}
			createDmoVesselsErr := tx.Create(&dmoVessels).Error

			if createDmoVesselsErr != nil {
				tx.Rollback()
				return createdDmo, createDmoVesselsErr
			}
		}
	}

	var traderDmo []traderdmo.TraderDmo

	var lastCount = 0
	for idx, value := range dmoInput.Trader {
		var traderDummy traderdmo.TraderDmo

		traderDummy.DmoId = createdDmo.ID
		traderDummy.TraderId = uint(value)
		traderDummy.Order = idx + 1
		lastCount = idx + 1
		traderDmo = append(traderDmo, traderDummy)
	}

	var traderEndUser traderdmo.TraderDmo
	traderEndUser.DmoId = createdDmo.ID
	traderEndUser.TraderId = uint(dmoInput.EndUser)
	traderEndUser.Order = lastCount + 1
	traderEndUser.IsEndUser = true
	traderDmo = append(traderDmo, traderEndUser)

	createTraderDmoErr := tx.Create(&traderDmo).Error

	if createTraderDmoErr != nil {
		tx.Rollback()
		return createdDmo, createTraderDmoErr
	}

	var history History

	history.DmoId = &createdDmo.ID
	history.Status = "Created Dmo"
	history.UserId = userId
	history.IupopkId = &iup.ID
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdDmo, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("ba_end_user", counterTransaction.BaEndUser+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdDmo, updateCounterErr
	}

	tx.Commit()
	return createdDmo, nil
}

func (r *repository) DeleteDmo(idDmo int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var dmo dmo.Dmo

	findDmoErr := tx.Where("id = ? AND iupopk_id = ?", idDmo, iupopkId).First(&dmo).Error

	if findDmoErr != nil {
		tx.Rollback()
		return false, findDmoErr
	}

	errDelete := tx.Unscoped().Where("id = ?", idDmo).Delete(&dmo).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	iupId := uint(iupopkId)
	history.Status = fmt.Sprintf("Deleted Dmo with id number %s and id %v", *dmo.IdNumber, dmo.ID)
	history.UserId = userId
	history.IupopkId = &iupId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint, iupopkId int) (dmo.Dmo, error) {
	tx := r.db.Begin()
	var dmoUpdate dmo.Dmo

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&dmoUpdate).Error

	if errFind != nil {
		tx.Rollback()
		return dmoUpdate, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "bast") {
				editData["bast_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "berita_acara.pdf") {
				if dmoUpdate.IsDocumentCustom {
					editData["signed_reconciliation_letter_document_link"] = value["Location"]
				}
				editData["reconciliation_letter_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "surat_pernyataan") {
				editData["statement_letter_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "berita_acara_pengguna_akhir") {
				editData["reconciliation_letter_end_user_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "recapdmo") {
				editData["recap_dmo_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "detaildmo") {
				editData["detail_dmo_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&dmoUpdate).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return dmoUpdate, errEdit
	}

	var history History

	history.DmoId = &dmoUpdate.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document dmo with id = %v", dmoUpdate.ID)
	history.IupopkId = &dmoUpdate.IupopkId
	dataInput, _ := json.Marshal(documentLink)
	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return dmoUpdate, errorBeforeDataJsonMarshal
	}

	var field string

	if isBast {
		field = "is_bast_document_downloaded"
	}

	if isStatementLetter {
		field = "is_statement_letter_downloaded"
	}

	if isReconciliationLetter {
		field = "is_reconciliation_letter_downloaded"
	}

	if isReconciliationLetterEndUser {
		field = "is_reconciliation_letter_end_user_downloaded"
	}

	updateErr := tx.Model(&dmoUpdate).Where("id = ? AND iupopk_id = ?", id, iupopkId).Update(field, true).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, updateErr
	}

	var history History

	history.DmoId = &dmoUpdate.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update %v dmo", field)

	afterData, _ := json.Marshal(dmoUpdate)
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = &dmoUpdate.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string, iupopkId int) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return dmoUpdate, errorBeforeDataJsonMarshal
	}

	var field string
	updatesDmo := make(map[string]interface{})

	if isBast {
		field = "bast"
		updatesDmo["is_bast_document_signed"] = true
		updatesDmo["signed_bast_document_link"] = location
	}

	if isStatementLetter {
		field = "statement_letter"
		updatesDmo["is_statement_letter_signed"] = true
		updatesDmo["signed_statement_letter_document_link"] = location
	}

	if isReconciliationLetter {
		field = "reconciliation_letter"
		updatesDmo["is_reconciliation_letter_signed"] = true
		updatesDmo["signed_reconciliation_letter_document_link"] = location
	}

	if isReconciliationLetterEndUser {
		field = "reconciliation_letter_end_user"
		updatesDmo["is_reconciliation_letter_end_user_signed"] = true
		updatesDmo["signed_reconciliation_letter_end_user_document_link"] = location
	}

	updateErr := tx.Model(&dmoUpdate).Where("id = ? AND iupopk_id = ?", id, iupopkId).Updates(updatesDmo).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, updateErr
	}

	var history History

	history.DmoId = &dmoUpdate.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update %v dmo", field)

	afterData, _ := json.Marshal(dmoUpdate)
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = &dmoUpdate.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return dmoUpdate, errorBeforeDataJsonMarshal
	}

	var field string
	updatesDmo := make(map[string]interface{})

	if isBast {
		field = "signed_bast false"
		updatesDmo["is_bast_document_signed"] = false
		updatesDmo["signed_bast_document_link"] = nil
	}

	if isStatementLetter {
		field = "signed_statement_letter false"
		updatesDmo["is_statement_letter_signed"] = false
		updatesDmo["signed_statement_letter_document_link"] = nil
	}

	if isReconciliationLetter {
		field = "signed_reconciliation_letter false"
		updatesDmo["is_reconciliation_letter_signed"] = false
		updatesDmo["signed_reconciliation_letter_document_link"] = nil
	}

	if isReconciliationLetterEndUser {
		field = "signed_reconciliation_letter_end_user false"
		updatesDmo["is_reconciliation_letter_end_user_signed"] = false
		updatesDmo["signed_reconciliation_letter_end_user_document_link"] = nil
	}

	updateErr := tx.Model(&dmoUpdate).Where("id = ?", id).Updates(updatesDmo).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, updateErr
	}

	var history History

	history.DmoId = &dmoUpdate.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update %v dmo", field)

	afterData, _ := json.Marshal(dmoUpdate)
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = &dmoUpdate.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	var updatedDmo dmo.Dmo
	var transactionBarge []transaction.Transaction
	var groupingVessel []groupingvesseldn.GroupingVesselDn

	barge := false
	vessel := false
	tx := r.db.Begin()
	beforeUpdate := make(map[string]interface{})
	afterUpdate := make(map[string]interface{})
	errFind := tx.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&updatedDmo).Error

	if errFind != nil {
		tx.Rollback()
		return updatedDmo, errFind
	}

	beforeUpdate["dmo"] = updatedDmo

	var transactionBefore []transaction.Transaction

	errFindTransactionBefore := tx.Where("dmo_id = ? AND grouping_vessel_dn_id IS NULL AND seller_id = ?", id, iupopkId).Find(&transactionBefore).Error

	if errFindTransactionBefore != nil {
		tx.Rollback()
		return updatedDmo, errFindTransactionBefore
	}

	beforeUpdate["transaction"] = transactionBefore

	errUpdateTransactionBefore := tx.Table("transactions").Where("dmo_id = ? AND seller_id = ?", id, iupopkId).Update("dmo_id", nil).Error

	if errUpdateTransactionBefore != nil {
		tx.Rollback()
		return updatedDmo, errUpdateTransactionBefore
	}

	var dmoVesselBefore []dmovessel.DmoVessel
	var groupingVesselIdBefore []uint
	errFindDmoVessel := tx.Where("dmo_id = ?", id).Find(&dmoVesselBefore).Error

	if errFindDmoVessel != nil {
		tx.Rollback()
		return updatedDmo, errFindDmoVessel
	}

	for _, v := range dmoVesselBefore {
		groupingVesselIdBefore = append(groupingVesselIdBefore, v.ID)
	}

	var groupingVesselBefore []groupingvesseldn.GroupingVesselDn

	errFindGroupingVesselBefore := tx.Where("id IN ? AND iupopk_id = ?", groupingVesselIdBefore, iupopkId).Find(&groupingVesselBefore).Error

	if errFindGroupingVesselBefore != nil {
		tx.Rollback()
		return updatedDmo, errFindGroupingVesselBefore
	}

	beforeUpdate["grouping_vessel"] = groupingVesselBefore

	errDeleteDmoVessel := tx.Table("dmo_vessels").Unscoped().Where("dmo_id = ?", id).Delete(&dmoVesselBefore).Error

	if errDeleteDmoVessel != nil {
		tx.Rollback()
		return updatedDmo, errDeleteDmoVessel
	}

	updatedMap := make(map[string]interface{})

	if len(dmoUpdateInput.TransactionBarge) > 0 {
		var bargeQuantity float64
		barge = true
		findTransactionBargeErr := tx.Where("id IN ? AND dmo_id IS NULL AND seller_id = ?", dmoUpdateInput.TransactionBarge, iupopkId).Find(&transactionBarge).Error

		if findTransactionBargeErr != nil {
			tx.Rollback()
			return updatedDmo, findTransactionBargeErr
		}

		if len(transactionBarge) != len(dmoUpdateInput.TransactionBarge) {
			tx.Rollback()
			return updatedDmo, errors.New("Ada transaksi yang sudah digunakan")
		}

		for _, v := range transactionBarge {
			bargeQuantity += v.QuantityUnloading
		}

		updatedMap["barge_total_quantity"] = math.Round(bargeQuantity*1000) / 1000
		updatedMap["barge_grand_total_quantity"] = math.Round(bargeQuantity*1000) / 1000
	}

	if len(dmoUpdateInput.GroupingVessel) > 0 {
		var vesselQuantity float64
		var vesselAdjustment float64
		var vesselGrandTotalQuantity float64
		vessel = true
		var checkDmoGrouping []dmovessel.DmoVessel

		findCheckDmoGroupingErr := tx.Where("grouping_vessel_dn_id IN ?", dmoUpdateInput.GroupingVessel).Find(&checkDmoGrouping).Error

		if findCheckDmoGroupingErr != nil {
			tx.Rollback()
			return updatedDmo, findCheckDmoGroupingErr
		}

		if len(checkDmoGrouping) > 0 {
			tx.Rollback()
			return updatedDmo, errors.New("Ada grouping vessel yang sudah digunakan")
		}

		findGroupingVesselErr := tx.Where("id IN ? AND iupopk_id = ?", dmoUpdateInput.GroupingVessel, iupopkId).Find(&groupingVessel).Error

		if findGroupingVesselErr != nil {
			tx.Rollback()
			return updatedDmo, findGroupingVesselErr
		}

		for _, v := range groupingVessel {
			vesselQuantity += v.Quantity
			vesselAdjustment += v.Adjustment
			vesselGrandTotalQuantity += v.GrandTotalQuantity
		}

		updatedMap["vessel_adjustment"] = math.Round(vesselAdjustment*1000) / 1000
		updatedMap["vessel_total_quantity"] = math.Round(vesselQuantity*1000) / 1000
		updatedMap["vessel_grand_total_quantity"] = math.Round(vesselGrandTotalQuantity*1000) / 1000
	}

	if barge && vessel {
		updatedMap["type"] = "Combination"
	}

	if barge && !vessel {
		updatedMap["type"] = "Barge"
	}

	if !barge && vessel {
		updatedMap["type"] = "Vessel"
	}

	updateDmoErr := tx.Model(&updatedDmo).Where("id = ? AND iupopk_id = ?", id, iupopkId).Updates(updatedMap).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return updatedDmo, updateDmoErr
	}

	afterUpdate["dmo"] = updatedDmo

	if len(dmoUpdateInput.TransactionBarge) > 0 {
		var transactionAfter []transaction.Transaction

		updateTransactionBargeErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", dmoUpdateInput.TransactionBarge, iupopkId).Update("dmo_id", updatedDmo.ID).Error

		if updateTransactionBargeErr != nil {
			tx.Rollback()
			return updatedDmo, updateTransactionBargeErr
		}

		errFindTransactionAfter := tx.Where("id IN ? AND seller_id = ?", dmoUpdateInput.TransactionBarge, iupopkId).Find(&transactionAfter).Error

		if errFindTransactionAfter != nil {
			tx.Rollback()
			return updatedDmo, errFindTransactionAfter
		}

		afterUpdate["transaction"] = transactionAfter
	}

	if len(dmoUpdateInput.GroupingVessel) > 0 {
		var transactionGroupVessel []transaction.Transaction

		var listIdTransactionGroupVessel []uint
		findTransactionGroupVesselErr := tx.Where("grouping_vessel_dn_id IN ? AND dmo_id IS NULL AND seller_id = ?", dmoUpdateInput.GroupingVessel, iupopkId).Find(&transactionGroupVessel).Error

		if findTransactionGroupVesselErr != nil {
			tx.Rollback()
			return updatedDmo, findTransactionGroupVesselErr
		}

		for _, v := range transactionGroupVessel {
			listIdTransactionGroupVessel = append(listIdTransactionGroupVessel, v.ID)
		}

		updateTransactionGroupVesselErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", listIdTransactionGroupVessel, iupopkId).Update("dmo_id", updatedDmo.ID).Error

		if updateTransactionGroupVesselErr != nil {
			tx.Rollback()
			return updatedDmo, updateTransactionGroupVesselErr
		}

		var dmoVessels []dmovessel.DmoVessel

		if len(dmoUpdateInput.GroupingVessel) > 0 {

			for _, v := range dmoUpdateInput.GroupingVessel {
				var dmoVesselDummy dmovessel.DmoVessel
				dmoVesselDummy.DmoId = updatedDmo.ID
				dmoVesselDummy.GroupingVesselDnId = uint(v)
				dmoVessels = append(dmoVessels, dmoVesselDummy)
			}
			createDmoVesselsErr := tx.Create(&dmoVessels).Error

			if createDmoVesselsErr != nil {
				tx.Rollback()
				return updatedDmo, createDmoVesselsErr
			}
		}

		var groupingVesselAfter []groupingvesseldn.GroupingVesselDn

		errFindGroupingVesselAfter := tx.Where("id IN ? AND iupopk_id = ?", dmoUpdateInput.GroupingVessel, iupopkId).Find(&groupingVesselAfter).Error

		if errFindGroupingVesselAfter != nil {
			tx.Rollback()
			return updatedDmo, errFindGroupingVesselAfter
		}

		afterUpdate["grouping_vessel"] = groupingVesselAfter
	}

	beforeUpdateJson, errorBeforeDataJsonMarshal := json.Marshal(beforeUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedDmo, errorBeforeDataJsonMarshal
	}

	afterUpdateJson, errorAfterDataJsonMarshal := json.Marshal(afterUpdate)

	if errorAfterDataJsonMarshal != nil {
		tx.Rollback()
		return updatedDmo, errorAfterDataJsonMarshal
	}

	var history History

	history.DmoId = &updatedDmo.ID
	history.Status = "Updated Dmo"
	history.UserId = userId
	history.BeforeData = beforeUpdateJson
	history.AfterData = afterUpdateJson
	history.IupopkId = &updatedDmo.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedDmo, createHistoryErr
	}

	tx.Commit()
	return updatedDmo, nil
}

// Production
func (r *repository) CreateProduction(input production.InputCreateProduction, userId uint, iupopkId int) (production.Production, error) {
	var createdProduction production.Production

	createdProduction.ProductionDate = input.ProductionDate
	createdProduction.Quantity = math.Round(input.Quantity*1000) / 1000
	createdProduction.IupopkId = uint(iupopkId)
	tx := r.db.Begin()

	errCreateProduction := tx.Create(&createdProduction).Error

	if errCreateProduction != nil {
		tx.Rollback()
		return createdProduction, errCreateProduction
	}

	var history History

	history.ProductionId = &createdProduction.ID
	history.Status = "Created Production"
	history.UserId = userId
	history.IupopkId = &createdProduction.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdProduction, createHistoryErr
	}

	tx.Commit()
	return createdProduction, nil
}

func (r *repository) UpdateProduction(input production.InputCreateProduction, productionId int, userId uint, iupopkId int) (production.Production, error) {
	var updatedProduction production.Production

	tx := r.db.Begin()

	errFindProduction := tx.Where("id = ? AND iupopk_id = ?", productionId, iupopkId).First(&updatedProduction).Error

	if errFindProduction != nil {
		tx.Rollback()
		return updatedProduction, errFindProduction
	}

	beforeData, _ := json.Marshal(updatedProduction)

	editData := make(map[string]interface{})
	editData["production_date"] = input.ProductionDate
	editData["quantity"] = math.Round(input.Quantity*1000) / 1000

	errUpdateProduction := tx.Model(&updatedProduction).Updates(editData).Error

	if errUpdateProduction != nil {
		tx.Rollback()
		return updatedProduction, errUpdateProduction
	}

	afterData, _ := json.Marshal(updatedProduction)

	var history History

	history.ProductionId = &updatedProduction.ID
	history.Status = "Update Production"
	history.UserId = userId
	history.BeforeData = beforeData
	history.AfterData = afterData
	history.IupopkId = &updatedProduction.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedProduction, createHistoryErr
	}

	tx.Commit()
	return updatedProduction, nil
}

func (r *repository) DeleteProduction(productionId int, userId uint, iupopkId int) (bool, error) {
	tx := r.db.Begin()
	var deletedProduction production.Production

	findDeletedProductionErr := tx.Where("id = ? AND iupopk_id = ?", productionId, iupopkId).First(&deletedProduction).Error

	if findDeletedProductionErr != nil {
		tx.Rollback()
		return false, findDeletedProductionErr
	}

	beforeData, _ := json.Marshal(deletedProduction)

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", productionId, iupopkId).Delete(&deletedProduction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Production with id %v", deletedProduction.ID)
	history.UserId = userId
	history.BeforeData = beforeData
	history.IupopkId = &deletedProduction.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

// Grouping Vessel DN
func (r *repository) CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	var createdGroupingVesselDn groupingvesseldn.GroupingVesselDn
	var transactions []transaction.Transaction
	tx := r.db.Begin()
	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, findIupErr
	}

	findTransactionsErr := tx.Where("id IN ? AND dmo_id is NULL AND transaction_type = ? AND grouping_vessel_dn_id is NULL AND is_migration = ? AND is_not_claim = ? AND seller_id = ?", inputGrouping.ListTransactions, "DN", false, false, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, findTransactionsErr
	}

	if len(transactions) != len(inputGrouping.ListTransactions) {
		tx.Rollback()
		return createdGroupingVesselDn, errors.New("please check some of transactions not found or already created in another group or not validate")
	}

	createdGroupingVesselDn.VesselId = inputGrouping.VesselId
	createdGroupingVesselDn.Quantity = math.Round(inputGrouping.Quantity*1000) / 1000
	createdGroupingVesselDn.Adjustment = math.Round(inputGrouping.Adjustment*1000) / 1000
	createdGroupingVesselDn.GrandTotalQuantity = math.Round(inputGrouping.GrandTotalQuantity*1000) / 1000
	createdGroupingVesselDn.BlDate = inputGrouping.BlDate
	createdGroupingVesselDn.IsCoaFinish = inputGrouping.IsCoaFinish
	createdGroupingVesselDn.SalesSystem = inputGrouping.SalesSystem
	createdGroupingVesselDn.DestinationId = inputGrouping.DestinationId
	createdGroupingVesselDn.LoadingPortId = inputGrouping.LoadingPortId
	createdGroupingVesselDn.DestinationCountryId = inputGrouping.DestinationCountryId
	createdGroupingVesselDn.DmoDestinationPortId = inputGrouping.DmoDestinationPortId
	createdGroupingVesselDn.BuyerId = inputGrouping.BuyerId
	createdGroupingVesselDn.CowDate = inputGrouping.CowDate
	if inputGrouping.CowNumber != "" {
		cowUpper := strings.ToUpper(inputGrouping.CowNumber)
		createdGroupingVesselDn.CowNumber = &cowUpper
	} else {
		createdGroupingVesselDn.CowNumber = nil
	}
	createdGroupingVesselDn.CoaDate = inputGrouping.CoaDate

	if inputGrouping.CoaNumber != "" {
		coaUpper := strings.ToUpper(inputGrouping.CoaNumber)
		createdGroupingVesselDn.CoaNumber = &coaUpper
	} else {
		createdGroupingVesselDn.CoaNumber = nil
	}
	createdGroupingVesselDn.SkabDate = inputGrouping.SkabDate

	if inputGrouping.SkabNumber != "" {
		coaUpper := strings.ToUpper(inputGrouping.SkabNumber)
		createdGroupingVesselDn.SkabNumber = &coaUpper
	} else {
		createdGroupingVesselDn.SkabNumber = nil
	}
	createdGroupingVesselDn.QualityTmAr = inputGrouping.QualityTmAr
	createdGroupingVesselDn.QualityImAdb = inputGrouping.QualityImAdb
	createdGroupingVesselDn.QualityAshAr = inputGrouping.QualityAshAr
	createdGroupingVesselDn.QualityAshAdb = inputGrouping.QualityAshAdb
	createdGroupingVesselDn.QualityVmAdb = inputGrouping.QualityVmAdb
	createdGroupingVesselDn.QualityFcAdb = inputGrouping.QualityFcAdb
	createdGroupingVesselDn.QualityTsAr = inputGrouping.QualityTsAr
	createdGroupingVesselDn.QualityTsAdb = inputGrouping.QualityTsAdb
	createdGroupingVesselDn.QualityCaloriesAr = inputGrouping.QualityCaloriesAr
	createdGroupingVesselDn.QualityCaloriesAdb = inputGrouping.QualityCaloriesAdb
	createdGroupingVesselDn.BlNumber = strings.ToUpper(inputGrouping.BlNumber)
	createdGroupingVesselDn.IupopkId = uint(iupopkId)
	errCreatedGroupingVesselDn := tx.Create(&createdGroupingVesselDn).Error

	if errCreatedGroupingVesselDn != nil {
		tx.Rollback()
		return createdGroupingVesselDn, errCreatedGroupingVesselDn
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, findCounterTransactionErr
	}

	code := "GMD-"
	code += iup.Code

	idNumber := createIdNumber(code, uint(counterTransaction.GroupingMvDn))

	updatedGroupingVesselDnErr := tx.Model(&createdGroupingVesselDn).Update("id_number", idNumber).Error

	if updatedGroupingVesselDnErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, updatedGroupingVesselDnErr
	}

	beforeData := make(map[string]interface{})
	beforeData["transactions"] = transactions

	beforeDataJson, _ := json.Marshal(beforeData)

	updateTransactions := make(map[string]interface{})

	updateTransactions["vessel_id"] = inputGrouping.VesselId
	updateTransactions["grouping_vessel_dn_id"] = createdGroupingVesselDn.ID
	updateTransactions["dmo_buyer_id"] = inputGrouping.BuyerId
	updateTransactions["dmo_destination_port_id"] = inputGrouping.DmoDestinationPortId
	updateTransactions["destination_country_id"] = inputGrouping.DestinationCountryId
	updateTransactions["destination_id"] = inputGrouping.DestinationId

	errUpdateTransactions := tx.Table("transactions").Where("id IN ?", inputGrouping.ListTransactions).Updates(updateTransactions).Error

	if errUpdateTransactions != nil {
		tx.Rollback()
		return createdGroupingVesselDn, errUpdateTransactions
	}

	afterData := make(map[string]interface{})
	afterData["transactions"] = transactions
	afterDataJson, _ := json.Marshal(afterData)
	var history History

	history.GroupingVesselDnId = &createdGroupingVesselDn.ID
	history.Status = "Created Grouping Vessel DN"
	history.UserId = userId
	history.BeforeData = beforeDataJson
	history.AfterData = afterDataJson
	history.IupopkId = &createdGroupingVesselDn.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("grouping_mv_dn", counterTransaction.GroupingMvDn+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdGroupingVesselDn, updateCounterErr
	}

	tx.Commit()

	return createdGroupingVesselDn, nil
}

func (r *repository) EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	var updatedGroupingVesselDn groupingvesseldn.GroupingVesselDn
	tx := r.db.Begin()

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&updatedGroupingVesselDn).Error

	if errFind != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errFind
	}

	var transactions []transaction.Transaction
	var listIdTransaction []uint
	errFindTransaction := tx.Where("grouping_vessel_dn_id = ?", id).Find(&transactions).Error

	if errFindTransaction != nil {
		return updatedGroupingVesselDn, errFindTransaction
	}

	for _, trans := range transactions {
		listIdTransaction = append(listIdTransaction, trans.ID)
	}

	beforeData := make(map[string]interface{})
	beforeData["grouping_vessel"] = updatedGroupingVesselDn
	beforeData["transactions"] = transactions
	beforeDataJson, errorBeforeDataJsonMarshal := json.Marshal(beforeData)

	updateTransactions := make(map[string]interface{})

	updateTransactions["vessel_id"] = editGrouping.VesselId
	updateTransactions["dmo_buyer_id"] = editGrouping.BuyerId
	updateTransactions["loading_port_id"] = editGrouping.LoadingPortId
	updateTransactions["dmo_destination_port_id"] = editGrouping.DmoDestinationPortId
	updateTransactions["destination_country_id"] = editGrouping.DestinationCountryId
	updateTransactions["destination_id"] = editGrouping.DestinationId

	errUpdateTransaction := tx.Table("transactions").Where("id IN ?", listIdTransaction).Updates(updateTransactions).Error

	if errUpdateTransaction != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errUpdateTransaction
	}

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errorBeforeDataJsonMarshal
	}

	editGrouping.CowNumber = strings.ToUpper(editGrouping.CowNumber)
	editGrouping.CoaNumber = strings.ToUpper(editGrouping.CoaNumber)
	editGrouping.BlNumber = strings.ToUpper(editGrouping.BlNumber)
	editGrouping.SkabNumber = strings.ToUpper(editGrouping.SkabNumber)
	editGroupingVesselDn, errorMarshal := json.Marshal(editGrouping)

	if errorMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errorMarshal
	}

	var editGroupingVesselDnInput map[string]interface{}

	errorUnmarshalGroupingVesselLn := json.Unmarshal(editGroupingVesselDn, &editGroupingVesselDnInput)

	if errorUnmarshalGroupingVesselLn != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errorUnmarshalGroupingVesselLn
	}

	updateGroupingVesselErr := tx.Model(&updatedGroupingVesselDn).Updates(editGroupingVesselDnInput).Error

	if updateGroupingVesselErr != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, updateGroupingVesselErr
	}

	afterData := make(map[string]interface{})
	afterData["grouping_vessel"] = updatedGroupingVesselDn
	afterData["transactions"] = transactions
	afterDataJson, errorAfterDataJsonMarshal := json.Marshal(afterData)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, errorAfterDataJsonMarshal
	}

	var history History

	history.GroupingVesselDnId = &updatedGroupingVesselDn.ID
	history.Status = "Updated Grouping Vessel DN"
	history.UserId = userId
	history.BeforeData = beforeDataJson
	history.AfterData = afterDataJson

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedGroupingVesselDn, createHistoryErr
	}

	tx.Commit()
	return updatedGroupingVesselDn, nil
}

func (r *repository) DeleteGroupingVesselDn(id int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var groupingVesselDn groupingvesseldn.GroupingVesselDn

	findGroupingVesselDnErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&groupingVesselDn).Error

	if findGroupingVesselDnErr != nil {
		fmt.Println(findGroupingVesselDnErr.Error())
		tx.Rollback()
		return false, findGroupingVesselDnErr
	}

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", id, iupopkId).Delete(&groupingVesselDn).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Grouping Vessel Dn with id number %s and id %v", *groupingVesselDn.IdNumber, groupingVesselDn.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	var uploadedGroupingVesselDn groupingvesseldn.GroupingVesselDn

	tx := r.db.Begin()

	errFind := tx.Where("id = ?", id).First(&uploadedGroupingVesselDn).Error

	if errFind != nil {
		return uploadedGroupingVesselDn, errFind
	}

	var isReupload = false
	editData := make(map[string]interface{})

	switch documentType {
	case "coa_cow":
		if uploadedGroupingVesselDn.CoaCowDocumentLink != nil {
			isReupload = true
		}
		editData["coa_cow_document_link"] = urlS3
	case "bl_mv":
		if uploadedGroupingVesselDn.BlMvDocumentLink != nil {
			isReupload = true
		}
		editData["bl_mv_document_link"] = urlS3
	case "skab":
		if uploadedGroupingVesselDn.SkabDocumentLink != nil {
			isReupload = true
		}
		editData["skab_document_link"] = urlS3
	}

	errEdit := tx.Model(&uploadedGroupingVesselDn).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return uploadedGroupingVesselDn, errEdit
	}
	var history History

	history.GroupingVesselDnId = &uploadedGroupingVesselDn.ID
	history.UserId = userId
	if isReupload == false {
		history.Status = fmt.Sprintf("Uploaded %s document", documentType)
	}

	if isReupload == true {
		history.Status = fmt.Sprintf("Reupload %s document", documentType)
	}

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return uploadedGroupingVesselDn, createHistoryErr
	}

	tx.Commit()
	return uploadedGroupingVesselDn, nil
}

// Grouping Vessel LN
func (r *repository) CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	var createdGroupingVesselLn groupingvesselln.GroupingVesselLn
	var transactions []transaction.Transaction
	tx := r.db.Begin()

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, findIupErr
	}

	findTransactionsErr := tx.Where("id IN ? AND transaction_type = ? AND grouping_vessel_ln_id is NULL AND is_migration = ? AND is_not_claim = ? AND seller_id = ?", inputGrouping.ListTransactions, "LN", false, false, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, findTransactionsErr
	}

	if len(transactions) != len(inputGrouping.ListTransactions) {
		tx.Rollback()
		return createdGroupingVesselLn, errors.New("please check some of transactions not found or already created in another group or not validate")
	}
	var createNavyCompanyMaster navycompany.NavyCompany

	if inputGrouping.NavyCompanyName != "" {

		createNavyCompanyMaster.Name = inputGrouping.NavyCompanyName

		errCreateNavyCompanyMaster := tx.FirstOrCreate(&createNavyCompanyMaster, createNavyCompanyMaster).Error

		if errCreateNavyCompanyMaster != nil {
			tx.Rollback()
			return createdGroupingVesselLn, errCreateNavyCompanyMaster
		}
	}

	var createNavyShipMaster navyship.NavyShip

	if inputGrouping.NavyShipName != "" {

		createNavyShipMaster.Name = inputGrouping.NavyShipName

		errCreateNavyShipMaster := tx.FirstOrCreate(&createNavyShipMaster, createNavyShipMaster).Error

		if errCreateNavyShipMaster != nil {
			tx.Rollback()
			return createdGroupingVesselLn, errCreateNavyShipMaster
		}
	}

	createdGroupingVesselLn.VesselId = inputGrouping.VesselId
	createdGroupingVesselLn.Quantity = math.Round(inputGrouping.Quantity*1000) / 1000
	createdGroupingVesselLn.Adjustment = math.Round(inputGrouping.Adjustment*1000) / 1000
	createdGroupingVesselLn.GrandTotalQuantity = math.Round(inputGrouping.GrandTotalQuantity*1000) / 1000
	createdGroupingVesselLn.DocumentTypeId = inputGrouping.DocumentTypeId
	createdGroupingVesselLn.AjuNumber = strings.ToUpper(inputGrouping.AjuNumber)
	createdGroupingVesselLn.PebRegisterNumber = strings.ToUpper(inputGrouping.PebRegisterNumber)
	createdGroupingVesselLn.PebRegisterDate = inputGrouping.PebRegisterDate
	createdGroupingVesselLn.PabeanOfficeId = inputGrouping.PabeanOfficeId
	createdGroupingVesselLn.SeriesPebGoods = inputGrouping.SeriesPebGoods
	createdGroupingVesselLn.DescriptionOfGoods = inputGrouping.DescriptionOfGoods
	createdGroupingVesselLn.TarifPosHs = strings.ToUpper(inputGrouping.TarifPosHs)
	createdGroupingVesselLn.PebQuantity = inputGrouping.PebQuantity
	createdGroupingVesselLn.PebUnitId = inputGrouping.PebUnitId
	createdGroupingVesselLn.ExportValue = inputGrouping.ExportValue
	createdGroupingVesselLn.CurrencyId = inputGrouping.CurrencyId
	createdGroupingVesselLn.LoadingPortId = inputGrouping.LoadingPortId
	createdGroupingVesselLn.SkaCooNumber = strings.ToUpper(inputGrouping.SkaCooNumber)
	createdGroupingVesselLn.SkaCooDate = inputGrouping.SkaCooDate
	createdGroupingVesselLn.DestinationCountryId = inputGrouping.DestinationCountryId
	createdGroupingVesselLn.LsExportNumber = strings.ToUpper(inputGrouping.LsExportNumber)
	createdGroupingVesselLn.LsExportDate = inputGrouping.LsExportDate
	createdGroupingVesselLn.InsuranceCompanyId = inputGrouping.InsuranceCompanyId
	createdGroupingVesselLn.PolisNumber = strings.ToUpper(inputGrouping.PolisNumber)

	if createNavyCompanyMaster.ID != 0 {
		createdGroupingVesselLn.NavyCompanyId = &createNavyCompanyMaster.ID
	}

	if createNavyShipMaster.ID != 0 {
		createdGroupingVesselLn.NavyShipId = &createNavyShipMaster.ID
	}

	createdGroupingVesselLn.NavyImoNumber = strings.ToUpper(inputGrouping.NavyImoNumber)
	createdGroupingVesselLn.Deadweight = inputGrouping.Deadweight
	createdGroupingVesselLn.IsCoaFinish = inputGrouping.IsCoaFinish
	createdGroupingVesselLn.CowDate = inputGrouping.CowDate
	createdGroupingVesselLn.CowNumber = strings.ToUpper(inputGrouping.CowNumber)
	createdGroupingVesselLn.CoaDate = inputGrouping.CoaDate
	createdGroupingVesselLn.CoaNumber = strings.ToUpper(inputGrouping.CoaNumber)
	createdGroupingVesselLn.QualityTmAr = inputGrouping.QualityTmAr
	createdGroupingVesselLn.QualityImAdb = inputGrouping.QualityImAdb
	createdGroupingVesselLn.QualityAshAr = inputGrouping.QualityAshAr
	createdGroupingVesselLn.QualityAshAdb = inputGrouping.QualityAshAdb
	createdGroupingVesselLn.QualityVmAdb = inputGrouping.QualityVmAdb
	createdGroupingVesselLn.QualityFcAdb = inputGrouping.QualityFcAdb
	createdGroupingVesselLn.QualityTsAr = inputGrouping.QualityTsAr
	createdGroupingVesselLn.QualityTsAdb = inputGrouping.QualityTsAdb
	createdGroupingVesselLn.QualityCaloriesAr = inputGrouping.QualityCaloriesAr
	createdGroupingVesselLn.QualityCaloriesAdb = inputGrouping.QualityCaloriesAdb
	createdGroupingVesselLn.NettQualityCaloriesAr = inputGrouping.NettQualityCaloriesAr
	createdGroupingVesselLn.BlDate = inputGrouping.BlDate
	createdGroupingVesselLn.BlNumber = strings.ToUpper(inputGrouping.BlNumber)
	createdGroupingVesselLn.IupopkId = uint(iupopkId)

	errCreatedGroupingVesselLn := tx.Create(&createdGroupingVesselLn).Error

	if errCreatedGroupingVesselLn != nil {
		tx.Rollback()
		return createdGroupingVesselLn, errCreatedGroupingVesselLn
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, findCounterTransactionErr
	}
	code := "GML-"
	code += iup.Code

	idNumber := createIdNumber(code, uint(counterTransaction.GroupingMvLn))

	updatedGroupingVesselLnErr := tx.Model(&createdGroupingVesselLn).Update("id_number", idNumber).Error

	if updatedGroupingVesselLnErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, updatedGroupingVesselLnErr
	}

	beforeData := make(map[string]interface{})
	beforeData["transactions"] = transactions

	beforeDataJson, _ := json.Marshal(beforeData)

	updateTransactions := make(map[string]interface{})

	updateTransactions["vessel_id"] = inputGrouping.VesselId
	updateTransactions["grouping_vessel_ln_id"] = createdGroupingVesselLn.ID

	errUpdateTransactions := tx.Table("transactions").Where("id IN ?", inputGrouping.ListTransactions).Updates(updateTransactions).Error

	if errUpdateTransactions != nil {
		tx.Rollback()
		return createdGroupingVesselLn, errUpdateTransactions
	}

	afterData := make(map[string]interface{})
	afterData["transactions"] = transactions
	afterDataJson, _ := json.Marshal(afterData)
	var history History

	history.GroupingVesselLnId = &createdGroupingVesselLn.ID
	history.Status = "Created Grouping Vessel LN"
	history.UserId = userId
	history.AfterData = afterDataJson
	history.BeforeData = beforeDataJson
	history.IupopkId = &createdGroupingVesselLn.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("grouping_mv_ln", counterTransaction.GroupingMvLn+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, updateCounterErr
	}

	tx.Commit()

	return createdGroupingVesselLn, nil
}

func (r *repository) EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	var updatedGroupingVesselLn groupingvesselln.GroupingVesselLn
	tx := r.db.Begin()

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&updatedGroupingVesselLn).Error

	if errFind != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errFind
	}
	beforeData := make(map[string]interface{})
	afterData := make(map[string]interface{})
	beforeData["grouping_vessel_ln"] = updatedGroupingVesselLn

	var createNavyCompanyMaster navycompany.NavyCompany

	if editGrouping.NavyCompanyName != "" {

		createNavyCompanyMaster.Name = editGrouping.NavyCompanyName

		errCreateNavyCompanyMaster := tx.FirstOrCreate(&createNavyCompanyMaster, createNavyCompanyMaster).Error

		if errCreateNavyCompanyMaster != nil {
			tx.Rollback()
			return updatedGroupingVesselLn, errCreateNavyCompanyMaster
		}
	}

	var createNavyShipMaster navyship.NavyShip

	if editGrouping.NavyShipName != "" {

		createNavyShipMaster.Name = editGrouping.NavyShipName

		errCreateNavyShipMaster := tx.FirstOrCreate(&createNavyShipMaster, createNavyShipMaster).Error

		if errCreateNavyShipMaster != nil {
			tx.Rollback()
			return updatedGroupingVesselLn, errCreateNavyShipMaster
		}
	}

	if updatedGroupingVesselLn.VesselId != editGrouping.VesselId {
		var transactions []transaction.Transaction
		var listIdTransaction []uint
		errFindTransaction := tx.Where("grouping_vessel_ln_id = ?", id).Find(&transactions).Error

		if errFindTransaction != nil {
			return updatedGroupingVesselLn, errFindTransaction
		}

		beforeData["transactions"] = transactions
		for _, trans := range transactions {
			listIdTransaction = append(listIdTransaction, trans.ID)
		}

		errUpdateTransaction := tx.Table("transactions").Where("id IN ?", listIdTransaction).Update("vessel_id", editGrouping.VesselId).Error

		if errUpdateTransaction != nil {
			tx.Rollback()
			return updatedGroupingVesselLn, errUpdateTransaction
		}

		afterData["transactions"] = transactions
	}

	beforeDataJson, errorBeforeDataJsonMarshal := json.Marshal(beforeData)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errorBeforeDataJsonMarshal
	}

	editGrouping.AjuNumber = strings.ToUpper(editGrouping.AjuNumber)
	editGrouping.PebRegisterNumber = strings.ToUpper(editGrouping.PebRegisterNumber)
	editGrouping.SkaCooNumber = strings.ToUpper(editGrouping.SkaCooNumber)
	editGrouping.LsExportNumber = strings.ToUpper(editGrouping.LsExportNumber)
	editGrouping.PolisNumber = strings.ToUpper(editGrouping.PolisNumber)
	editGrouping.NavyImoNumber = strings.ToUpper(editGrouping.NavyImoNumber)
	editGrouping.CowNumber = strings.ToUpper(editGrouping.CowNumber)
	editGrouping.CoaNumber = strings.ToUpper(editGrouping.CoaNumber)
	editGrouping.BlNumber = strings.ToUpper(editGrouping.BlNumber)
	editGroupingVesselLn, errorMarshal := json.Marshal(editGrouping)

	if errorMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errorMarshal
	}

	var editGroupingVesselLnInput map[string]interface{}

	errorUnmarshalGroupingVesselLn := json.Unmarshal(editGroupingVesselLn, &editGroupingVesselLnInput)

	if errorUnmarshalGroupingVesselLn != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errorUnmarshalGroupingVesselLn
	}
	if createNavyCompanyMaster.ID != 0 {
		editGroupingVesselLnInput["navy_company_id"] = createNavyCompanyMaster.ID
	}

	if createNavyShipMaster.ID != 0 {
		editGroupingVesselLnInput["navy_ship_id"] = createNavyShipMaster.ID
	}

	delete(editGroupingVesselLnInput, "navy_company_name")
	delete(editGroupingVesselLnInput, "navy_ship_name")

	updateGroupingVesselErr := tx.Model(&updatedGroupingVesselLn).Updates(editGroupingVesselLnInput).Error

	if updateGroupingVesselErr != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, updateGroupingVesselErr
	}

	afterData["grouping_vessel_ln"] = updatedGroupingVesselLn
	afterDataJson, errorAfterDataJsonMarshal := json.Marshal(afterData)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errorAfterDataJsonMarshal
	}

	var history History

	history.GroupingVesselLnId = &updatedGroupingVesselLn.ID
	history.Status = "Updated Grouping Vessel LN"
	history.UserId = userId
	history.BeforeData = beforeDataJson
	history.AfterData = afterDataJson

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, createHistoryErr
	}

	tx.Commit()
	return updatedGroupingVesselLn, nil
}

func (r *repository) UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	var uploadedGroupingVesselLn groupingvesselln.GroupingVesselLn

	tx := r.db.Begin()

	errFind := tx.Where("id = ?", id).First(&uploadedGroupingVesselLn).Error

	if errFind != nil {
		return uploadedGroupingVesselLn, errFind
	}

	var isReupload = false
	editData := make(map[string]interface{})

	switch documentType {
	case "peb":
		if uploadedGroupingVesselLn.PebDocumentLink != "" {
			isReupload = true
		}
		editData["peb_document_link"] = urlS3
	case "insurance":
		if uploadedGroupingVesselLn.InsuranceDocumentLink != "" {
			isReupload = true
		}
		editData["insurance_document_link"] = urlS3
	case "ls_export":
		if uploadedGroupingVesselLn.LsExportDocumentLink != "" {
			isReupload = true
		}
		editData["ls_export_document_link"] = urlS3
	case "navy":
		if uploadedGroupingVesselLn.NavyDocumentLink != "" {
			isReupload = true
		}
		editData["navy_document_link"] = urlS3
	case "ska_coo":
		if uploadedGroupingVesselLn.SkaCooDocumentLink != "" {
			isReupload = true
		}
		editData["ska_coo_document_link"] = urlS3
	case "coa_cow":
		if uploadedGroupingVesselLn.CoaCowDocumentLink != "" {
			isReupload = true
		}
		editData["coa_cow_document_link"] = urlS3
	case "bl_mv":
		if uploadedGroupingVesselLn.BlMvDocumentLink != "" {
			isReupload = true
		}
		editData["bl_mv_document_link"] = urlS3
	}

	errEdit := tx.Model(&uploadedGroupingVesselLn).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return uploadedGroupingVesselLn, errEdit
	}
	var history History

	history.GroupingVesselLnId = &uploadedGroupingVesselLn.ID
	history.UserId = userId
	if isReupload == false {
		history.Status = fmt.Sprintf("Uploaded %s document", documentType)
	}

	if isReupload == true {
		history.Status = fmt.Sprintf("Reupload %s document", documentType)
	}

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return uploadedGroupingVesselLn, createHistoryErr
	}

	tx.Commit()
	return uploadedGroupingVesselLn, nil
}

func (r *repository) DeleteGroupingVesselLn(id int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var groupingVesselLn groupingvesselln.GroupingVesselLn

	findGroupingVesselLnErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&groupingVesselLn).Error

	if findGroupingVesselLnErr != nil {
		tx.Rollback()
		return false, findGroupingVesselLnErr
	}

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", id, iupopkId).Delete(&groupingVesselLn).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Grouping Vessel Ln with id number %s and id %v", *groupingVesselLn.IdNumber, groupingVesselLn.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

// Minerba LN

func (r *repository) CreateMinerbaLn(period string, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	var createdMinerbaLn minerbaln.MinerbaLn

	tx := r.db.Begin()
	createdMinerbaLn.Period = period
	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ? AND transaction_type = ? AND minerba_ln_id is NULL AND is_finance_check = ? AND seller_id = ?", listTransactions, "LN", true, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdMinerbaLn, findTransactionsErr
	}

	if len(transactions) != len(listTransactions) {
		tx.Rollback()
		return createdMinerbaLn, errors.New("please check some of transactions not found or already in report")
	}

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdMinerbaLn, findIupErr
	}

	var tempQuantity float64
	for _, v := range transactions {
		tempQuantity += v.QuantityUnloading
	}

	stringTempQuantity := fmt.Sprintf("%.3f", tempQuantity)
	parseTempQuantity, _ := strconv.ParseFloat(stringTempQuantity, 64)

	createdMinerbaLn.Quantity = math.Round(parseTempQuantity*1000) / 1000
	createdMinerbaLn.IupopkId = uint(iupopkId)

	errCreateMinerbaLn := tx.Create(&createdMinerbaLn).Error

	if errCreateMinerbaLn != nil {
		tx.Rollback()
		return createdMinerbaLn, errCreateMinerbaLn
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdMinerbaLn, findCounterTransactionErr
	}

	periodSplit := strings.Split(period, " ")

	idNumber := fmt.Sprintf("LSL-%s-%s-%s-", iup.Code, helper.MonthStringToNumberString(periodSplit[0]), periodSplit[1][len(periodSplit[1])-2:])
	idNumber += helper.CreateIdNumber(counterTransaction.Sp3meln)

	updateMinerbaErr := tx.Model(&createdMinerbaLn).Update("id_number", idNumber).Error

	if updateMinerbaErr != nil {
		tx.Rollback()
		return createdMinerbaLn, updateMinerbaErr
	}

	listTransactionsErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", listTransactions, iupopkId).Update("minerba_ln_id", createdMinerbaLn.ID).Error

	if listTransactionsErr != nil {
		tx.Rollback()
		return createdMinerbaLn, listTransactionsErr
	}

	var history History

	history.MinerbaLnId = &createdMinerbaLn.ID
	history.Status = "Created Minerba LN Report"
	history.UserId = userId
	history.IupopkId = &createdMinerbaLn.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdMinerbaLn, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("sp3meln", counterTransaction.Sp3meln+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdMinerbaLn, updateCounterErr
	}

	tx.Commit()
	return createdMinerbaLn, nil
}

func (r *repository) UpdateMinerbaLn(id int, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	var updatedMinerbaLn minerbaln.MinerbaLn
	var quantityMinerba float64

	historyBefore := make(map[string]interface{})
	historyAfter := make(map[string]interface{})
	tx := r.db.Begin()

	findMinerbaLnErr := tx.Where("id = ?", id).First(&updatedMinerbaLn).Error

	if findMinerbaLnErr != nil {
		return updatedMinerbaLn, findMinerbaLnErr
	}

	historyBefore["minerba_ln"] = updatedMinerbaLn

	var beforeTransaction []transaction.Transaction
	findTransactionBeforeErr := tx.Where("minerba_ln_id = ?", id).Find(&beforeTransaction).Error

	if findTransactionBeforeErr != nil {
		return updatedMinerbaLn, findTransactionBeforeErr
	}

	var transactionBefore []uint

	for _, v := range beforeTransaction {
		transactionBefore = append(transactionBefore, v.ID)
	}

	historyBefore["transactions"] = transactionBefore

	errUpdMinerbaNil := tx.Model(&beforeTransaction).Where("minerba_ln_id = ?", id).Update("minerba_ln_id", nil).Error

	if errUpdMinerbaNil != nil {
		tx.Rollback()
		return updatedMinerbaLn, errUpdMinerbaNil
	}

	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ? AND transaction_type = ? AND minerba_ln_id is NULL AND is_finance_check = ?", listTransactions, "LN", true).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return updatedMinerbaLn, findTransactionsErr
	}

	if len(transactions) != len(listTransactions) {
		tx.Rollback()
		return updatedMinerbaLn, errors.New("please check some of transactions not found")
	}

	for _, v := range transactions {
		quantityMinerba += v.QuantityUnloading
	}

	quantityMinerba = math.Round(quantityMinerba*1000) / 1000
	errUpdateMinerba := tx.Model(&updatedMinerbaLn).Update("quantity", quantityMinerba).Error

	if errUpdateMinerba != nil {
		tx.Rollback()
		return updatedMinerbaLn, errUpdateMinerba
	}

	historyAfter["minerba_ln"] = updatedMinerbaLn
	historyAfter["transactions"] = listTransactions

	listTransactionsErr := tx.Table("transactions").Where("id IN ?", listTransactions).Update("minerba_ln_id", id).Error

	if listTransactionsErr != nil {
		tx.Rollback()
		return updatedMinerbaLn, listTransactionsErr
	}

	var history History
	beforeData, errorBeforeDataJsonMarshal := json.Marshal(historyBefore)
	afterData, errorAfterDataJsonMarshal := json.Marshal(historyAfter)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedMinerbaLn, errorBeforeDataJsonMarshal
	}

	if errorAfterDataJsonMarshal != nil {
		tx.Rollback()
		return updatedMinerbaLn, errorAfterDataJsonMarshal
	}

	history.MinerbaLnId = &updatedMinerbaLn.ID
	history.Status = "Updated Minerba Ln Report"
	history.UserId = userId
	history.AfterData = afterData
	history.BeforeData = beforeData
	history.IupopkId = &updatedMinerbaLn.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedMinerbaLn, createHistoryErr
	}

	tx.Commit()
	return updatedMinerbaLn, nil
}

func (r *repository) DeleteMinerbaLn(idMinerbaLn int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var minerbaLn minerbaln.MinerbaLn

	findMinerbaLnErr := tx.Where("id = ?", idMinerbaLn).First(&minerbaLn).Error

	if findMinerbaLnErr != nil {
		tx.Rollback()
		return false, findMinerbaLnErr
	}

	errDelete := tx.Unscoped().Where("id = ?", idMinerbaLn).Delete(&minerbaLn).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Minerba LN Report with id number %s and id %v", *minerbaLn.IdNumber, minerbaLn.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	tx := r.db.Begin()
	var minerbaLn minerbaln.MinerbaLn

	errFind := tx.Where("id = ?", id).First(&minerbaLn).Error

	if errFind != nil {
		tx.Rollback()
		return minerbaLn, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "sp3meln") {
				editData["sp3_meln_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&minerbaLn).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return minerbaLn, errEdit
	}

	var history History

	history.MinerbaLnId = &minerbaLn.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document minerba ln with id = %v", minerbaLn.ID)

	dataInput, _ := json.Marshal(documentLink)
	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return minerbaLn, createHistoryErr
	}

	tx.Commit()
	return minerbaLn, nil
}

// INSW

func (r *repository) CreateInsw(month string, year int, userId uint, iupopkId int) (insw.Insw, error) {
	tx := r.db.Begin()
	var createInsw insw.Insw

	var checkInsw insw.Insw

	findInswErr := tx.Where("month = ? AND year = ? AND iupopk_id = ?", month, year, iupopkId).First(&checkInsw).Error

	if findInswErr == nil {
		tx.Rollback()
		return createInsw, errors.New("Laporan INSW sudah pernah dibuat")
	}

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createInsw, findIupErr
	}

	firstOfMonth := time.Date(year, time.Month(helper.MonthLongToNumber(month)), 1, 0, 0, 0, 0, time.Local)

	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	var findGroupingVessel []groupingvesselln.GroupingVesselLn

	findGroupingVesselErr := tx.Where("peb_register_date >= ? AND peb_register_date <= ? AND insw_id is NULL AND iupopk_id = ?", firstOfMonth, lastOfMonth, iupopkId).Find(&findGroupingVessel).Error

	if findGroupingVesselErr != nil {
		tx.Rollback()
		return createInsw, findGroupingVesselErr
	}

	var idGroupingVessel []uint

	for _, v := range findGroupingVessel {
		idGroupingVessel = append(idGroupingVessel, v.ID)
	}

	createInsw.Month = month
	createInsw.Year = year
	createInsw.IupopkId = iup.ID
	createInswErr := tx.Create(&createInsw).Error

	if createInswErr != nil {
		tx.Rollback()
		return createInsw, createInswErr
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createInsw, findCounterTransactionErr
	}

	idNumber := fmt.Sprintf("LIW-%s-%s-%s-", iup.Code, helper.MonthStringToNumberString(month), strconv.Itoa(year)[len(strconv.Itoa(year))-2:])
	idNumber += helper.CreateIdNumber(counterTransaction.Insw)

	updateInswErr := tx.Model(&createInsw).Update("id_number", idNumber).Error

	if updateInswErr != nil {
		tx.Rollback()
		return createInsw, updateInswErr
	}

	updateGroupingVesselLnErr := tx.Table("grouping_vessel_lns").Where("id IN ? AND iupopk_id = ?", idGroupingVessel, iupopkId).Update("insw_id", createInsw.ID).Error

	if updateGroupingVesselLnErr != nil {
		tx.Rollback()
		return createInsw, updateGroupingVesselLnErr
	}

	var history History

	history.InswId = &createInsw.ID
	history.Status = "Created Insw Report"
	history.UserId = userId
	history.IupopkId = &createInsw.IupopkId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createInsw, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("insw", counterTransaction.Insw+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createInsw, updateCounterErr
	}

	tx.Commit()
	return createInsw, nil
}

func (r *repository) DeleteInsw(idInsw int, userId uint, iupopkId int) (bool, error) {
	tx := r.db.Begin()
	var inswData insw.Insw

	findInswErr := tx.Where("id = ? AND iupopk_id = ?", idInsw, iupopkId).First(&inswData).Error

	if findInswErr != nil {
		tx.Rollback()
		return false, findInswErr
	}

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", idInsw, iupopkId).Delete(&inswData).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted INSW Report with id number %s and id %v", *inswData.IdNumber, inswData.ID)
	history.UserId = userId
	history.IupopkId = &inswData.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint, iupopkId int) (insw.Insw, error) {
	tx := r.db.Begin()
	var insw insw.Insw

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&insw).Error

	if errFind != nil {
		tx.Rollback()
		return insw, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "insw") {
				editData["insw_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&insw).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return insw, errEdit
	}

	var history History

	history.InswId = &insw.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document insw with id = %v", insw.ID)
	history.IupopkId = &insw.IupopkId
	dataInput, _ := json.Marshal(documentLink)
	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return insw, createHistoryErr
	}

	tx.Commit()
	return insw, nil
}

// Report Dmo

func (r *repository) CreateReportDmo(input reportdmo.InputCreateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	var createdReportDmo reportdmo.ReportDmo

	tx := r.db.Begin()
	var transactionBarge []transaction.Transaction
	var groupingVessel []groupingvesseldn.GroupingVesselDn

	var iup iupopk.Iupopk

	findIupErr := tx.Where("id = ?", iupopkId).First(&iup).Error

	if findIupErr != nil {
		tx.Rollback()
		return createdReportDmo, findIupErr
	}

	createdReportDmo.Period = input.Period
	createdReportDmo.IupopkId = uint(iupopkId)
	if len(input.Transactions) > 0 {
		var bargeQuantity float64
		findTransactionBargeErr := tx.Where("id IN ? AND report_dmo_id IS NULL AND is_finance_check = ? AND transaction_type = ? AND seller_id = ?", input.Transactions, true, "DN", iupopkId).Find(&transactionBarge).Error

		if findTransactionBargeErr != nil {
			tx.Rollback()
			return createdReportDmo, findTransactionBargeErr
		}

		if len(transactionBarge) != len(input.Transactions) {
			tx.Rollback()
			return createdReportDmo, errors.New("Ada transaksi yang sudah digunakan")
		}

		for _, v := range transactionBarge {
			bargeQuantity += v.QuantityUnloading
		}

		createdReportDmo.Quantity = math.Round(bargeQuantity*1000) / 1000
	}

	if len(input.GroupingVessels) > 0 {
		var vesselQuantity float64

		findCheckDmoGroupingErr := tx.Where("id IN ? AND sales_system = ? AND report_dmo_id IS NULL AND iupopk_id = ?", input.GroupingVessels, "Vessel", iupopkId).Find(&groupingVessel).Error

		if findCheckDmoGroupingErr != nil {
			tx.Rollback()
			return createdReportDmo, findCheckDmoGroupingErr
		}

		if len(groupingVessel) != len(input.GroupingVessels) {
			tx.Rollback()
			return createdReportDmo, errors.New("please check grouping vessel is already used")
		}

		for _, v := range groupingVessel {
			vesselQuantity += v.Quantity
		}

		createdReportDmo.Quantity += vesselQuantity
	}

	createdReportDmo.Quantity = math.Round(createdReportDmo.Quantity*1000) / 1000
	if len(transactionBarge) != len(input.Transactions) && len(groupingVessel) != len(input.GroupingVessels) {
		tx.Rollback()
		return createdReportDmo, errors.New("please check some of transactions not found")
	}

	createReportDmoErr := tx.Create(&createdReportDmo).Error

	if createReportDmoErr != nil {
		tx.Rollback()
		return createdReportDmo, createReportDmoErr
	}

	var counterTransaction counter.Counter

	findCounterTransactionErr := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if findCounterTransactionErr != nil {
		tx.Rollback()
		return createdReportDmo, findCounterTransactionErr
	}

	periodSplit := strings.Split(input.Period, " ")

	idNumber := fmt.Sprintf("LDO-%s-%s-%s-", iup.Code, helper.MonthStringToNumberString(periodSplit[0]), periodSplit[1][len(periodSplit[1])-2:])
	idNumber += helper.CreateIdNumber(counterTransaction.Dmo)

	updateDmoErr := tx.Model(&createdReportDmo).Update("id_number", idNumber).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return createdReportDmo, updateDmoErr
	}

	if len(input.Transactions) > 0 {
		updateTransactionBargeErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", input.Transactions, iupopkId).Update("report_dmo_id", createdReportDmo.ID).Error

		if updateTransactionBargeErr != nil {
			tx.Rollback()
			return createdReportDmo, updateTransactionBargeErr
		}
	}

	if len(input.GroupingVessels) > 0 {

		updateGroupVesselErr := tx.Table("grouping_vessel_dns").Where("id IN ? AND iupopk_id = ?", input.GroupingVessels, iupopkId).Update("report_dmo_id", createdReportDmo.ID).Error

		if updateGroupVesselErr != nil {
			tx.Rollback()
			return createdReportDmo, updateGroupVesselErr
		}

		var transactionGroupVessel []transaction.Transaction

		var listIdTransactionGroupVessel []uint
		findTransactionGroupVesselErr := tx.Where("grouping_vessel_dn_id IN ? AND report_dmo_id IS NULL AND seller_id = ?", input.GroupingVessels, iupopkId).Find(&transactionGroupVessel).Error

		if findTransactionGroupVesselErr != nil {
			tx.Rollback()
			return createdReportDmo, findTransactionGroupVesselErr
		}

		for _, v := range transactionGroupVessel {
			listIdTransactionGroupVessel = append(listIdTransactionGroupVessel, v.ID)
		}

		updateTransactionGroupVesselErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", listIdTransactionGroupVessel, iupopkId).Update("report_dmo_id", createdReportDmo.ID).Error

		if updateTransactionGroupVesselErr != nil {
			tx.Rollback()
			return createdReportDmo, updateTransactionGroupVesselErr
		}
	}

	var history History

	history.ReportDmoId = &createdReportDmo.ID
	history.Status = "Created Report Dmo"
	history.UserId = userId
	history.IupopkId = &createdReportDmo.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdReportDmo, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("dmo", counterTransaction.Dmo+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdReportDmo, updateCounterErr
	}

	tx.Commit()
	return createdReportDmo, nil
}

func (r *repository) UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	tx := r.db.Begin()
	var reportDmo reportdmo.ReportDmo

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&reportDmo).Error

	if errFind != nil {
		tx.Rollback()
		return reportDmo, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "recapdmo") {
				editData["recap_dmo_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "detaildmo") {
				editData["detail_dmo_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&reportDmo).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return reportDmo, errEdit
	}

	var history History

	history.ReportDmoId = &reportDmo.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document report dmo with id = %v", reportDmo.ID)
	history.IupopkId = &reportDmo.IupopkId
	dataInput, _ := json.Marshal(documentLink)

	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return reportDmo, createHistoryErr
	}

	tx.Commit()
	return reportDmo, nil
}

func (r *repository) UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	var updateReportDmo reportdmo.ReportDmo
	var quantityReportDmo float64

	historyBefore := make(map[string]interface{})
	historyAfter := make(map[string]interface{})
	tx := r.db.Begin()

	findReportDmoErr := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&updateReportDmo).Error

	if findReportDmoErr != nil {
		return updateReportDmo, findReportDmoErr
	}

	historyBefore["report_dmo"] = updateReportDmo

	var salesSystem []salessystem.SalesSystem
	var salesSystemId []uint

	errFindSalesSystem := r.db.Where("name ILIKE '%Barge'").Find(&salesSystem).Error

	if errFindSalesSystem != nil {
		return updateReportDmo, errFindSalesSystem
	}

	for _, v := range salesSystem {
		salesSystemId = append(salesSystemId, v.ID)
	}

	var beforeTransaction []transaction.Transaction
	findTransactionBeforeErr := tx.Where("report_dmo_id = ? AND sales_system_id IN ? AND seller_id = ?", id, salesSystemId, iupopkId).Find(&beforeTransaction).Error

	if findTransactionBeforeErr != nil {
		return updateReportDmo, findTransactionBeforeErr
	}

	var transactionBefore []uint

	for _, v := range beforeTransaction {
		transactionBefore = append(transactionBefore, v.ID)
	}

	historyBefore["transactions"] = beforeTransaction

	var beforeGroupingVesselDn []groupingvesseldn.GroupingVesselDn

	findGroupingVesselDnErr := tx.Where("report_dmo_id = ? AND iupopk_id = ?", id, iupopkId).Find(&beforeGroupingVesselDn).Error

	if findGroupingVesselDnErr != nil {
		return updateReportDmo, findGroupingVesselDnErr
	}

	historyBefore["grouping_vessel_dn"] = beforeGroupingVesselDn

	errUpdReportDmoNull := tx.Model(&beforeTransaction).Where("report_dmo_id = ? AND seller_id = ?", id, iupopkId).Update("report_dmo_id", nil).Error

	if errUpdReportDmoNull != nil {
		tx.Rollback()
		return updateReportDmo, errUpdReportDmoNull
	}

	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ? AND transaction_type = ? AND report_dmo_id is NULL AND is_finance_check = ? AND sales_system_id IN ? AND seller_id = ?", inputUpdate.Transactions, "DN", true, salesSystemId, iupopkId).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return updateReportDmo, findTransactionsErr
	}

	if len(transactions) != len(inputUpdate.Transactions) {
		tx.Rollback()
		return updateReportDmo, errors.New("please check some of transactions not found")
	}

	for _, v := range transactions {
		quantityReportDmo += v.QuantityUnloading
	}

	listTransactionsErr := tx.Table("transactions").Where("id IN ? AND seller_id = ?", inputUpdate.Transactions, iupopkId).Update("report_dmo_id", id).Error

	if listTransactionsErr != nil {
		tx.Rollback()
		return updateReportDmo, listTransactionsErr
	}

	errUpdGroupingVesselNull := tx.Model(&beforeGroupingVesselDn).Where("report_dmo_id = ? AND iupopk_id = ?", id, iupopkId).Update("report_dmo_id", nil).Error

	if errUpdGroupingVesselNull != nil {
		tx.Rollback()
		return updateReportDmo, errUpdGroupingVesselNull
	}

	var groupingVesselDn []groupingvesseldn.GroupingVesselDn

	findGroupingVesselErr := tx.Where("id IN ? AND report_dmo_id is NULL AND iupopk_id = ?", inputUpdate.GroupingVessels, iupopkId).Find(&groupingVesselDn).Error

	if findGroupingVesselErr != nil {
		tx.Rollback()
		return updateReportDmo, findGroupingVesselErr
	}

	groupingVesselErr := tx.Table("grouping_vessel_dns").Where("id IN ? AND iupopk_id = ?", inputUpdate.GroupingVessels, iupopkId).Update("report_dmo_id", id).Error

	if groupingVesselErr != nil {
		tx.Rollback()
		return updateReportDmo, groupingVesselErr
	}

	for _, v := range groupingVesselDn {
		quantityReportDmo += v.GrandTotalQuantity
	}

	quantityReportDmo = math.Round(quantityReportDmo*1000) / 1000
	errUpdReportDmo := tx.Model(&updateReportDmo).Update("quantity", quantityReportDmo).Error

	if errUpdReportDmo != nil {
		tx.Rollback()
		return updateReportDmo, errUpdReportDmo
	}

	historyAfter["report_dmo"] = updateReportDmo
	historyAfter["transactions"] = transactions
	historyAfter["grouping_vessel_dn"] = groupingVesselDn

	var history History
	beforeData, errorBeforeDataJsonMarshal := json.Marshal(historyBefore)
	afterData, errorAfterDataJsonMarshal := json.Marshal(historyAfter)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updateReportDmo, errorBeforeDataJsonMarshal
	}

	if errorAfterDataJsonMarshal != nil {
		tx.Rollback()
		return updateReportDmo, errorAfterDataJsonMarshal
	}

	history.ReportDmoId = &updateReportDmo.ID
	history.Status = "Updated Report Dmo"
	history.UserId = userId
	history.AfterData = afterData
	history.BeforeData = beforeData
	history.IupopkId = &updateReportDmo.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updateReportDmo, createHistoryErr
	}

	tx.Commit()
	return updateReportDmo, nil
}

func (r *repository) DeleteReportDmo(idReportDmo int, userId uint, iupopkId int) (bool, error) {

	tx := r.db.Begin()
	var reportDmo reportdmo.ReportDmo

	findReportDmo := tx.Where("id = ? AND iupopk_id = ?", idReportDmo, iupopkId).First(&reportDmo).Error

	if findReportDmo != nil {
		tx.Rollback()
		return false, findReportDmo
	}

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", idReportDmo, iupopkId).Delete(&reportDmo).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Report Dmo with id number %s and id %v", *reportDmo.IdNumber, reportDmo.ID)
	history.UserId = userId
	history.IupopkId = &reportDmo.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

// COA Report

func (r *repository) CreateCoaReport(dateFrom string, dateTo string, iupopkId int, userId uint) (coareport.CoaReport, error) {
	var coaReport coareport.CoaReport
	var iup iupopk.Iupopk

	var counterTransaction counter.Counter

	tx := r.db.Begin()

	errFindIup := tx.Where("id = ?", iupopkId).First(&iup).Error

	if errFindIup != nil {
		tx.Rollback()
		return coaReport, errFindIup
	}

	errFindCounterTransaction := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if errFindCounterTransaction != nil {
		tx.Rollback()
		return coaReport, errFindCounterTransaction
	}

	errFind := tx.Where("date_from = ? AND date_to = ? AND iupopk_id = ?", dateFrom, dateTo, iupopkId).First(&coaReport).Error

	if errFind == nil {
		tx.Rollback()
		return coaReport, errors.New("Report already has been created")
	}

	idNumber := "RCO-" + iup.Code

	coaReport.IdNumber = createIdNumber(idNumber, uint(counterTransaction.CoaReport))
	coaReport.DateFrom = dateFrom
	coaReport.DateTo = dateTo
	coaReport.IupopkId = iup.ID

	errCreate := tx.Create(&coaReport).Error

	if errCreate != nil {
		tx.Rollback()
		return coaReport, errCreate
	}

	coaReportData, _ := json.Marshal(coaReport)

	var history History

	history.CoaReportId = &coaReport.ID
	history.Status = "Created Coa Report"
	history.UserId = userId
	history.IupopkId = &iup.ID
	history.BeforeData = coaReportData
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return coaReport, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("coa_report", counterTransaction.CoaReport+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return coaReport, updateCounterErr
	}

	tx.Commit()
	return coaReport, nil
}

func (r *repository) DeleteCoaReport(id int, iupopkId int, userId uint) (bool, error) {
	var coaReport coareport.CoaReport
	var iup iupopk.Iupopk

	tx := r.db.Begin()

	errFindIup := tx.Where("id = ?", iupopkId).First(&iup).Error

	if errFindIup != nil {
		tx.Rollback()
		return false, errFindIup
	}

	errFind := tx.Where("id = ?", id).First(&coaReport).Error

	if errFind != nil {
		tx.Rollback()
		return false, errFind
	}

	coaReportData, _ := json.Marshal(coaReport)

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", id, iupopkId).Delete(&coaReport).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Coa Report with id number %s and id %v", &coaReport.IdNumber, coaReport.ID)
	history.UserId = userId
	history.IupopkId = &iup.ID
	history.BeforeData = coaReportData
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentCoaReport(id int, documentLink coareport.InputUpdateDocumentCoaReport, userId uint, iupopkId int) (coareport.CoaReport, error) {
	tx := r.db.Begin()
	var coaReport coareport.CoaReport

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&coaReport).Error

	if errFind != nil {
		tx.Rollback()
		return coaReport, errFind
	}

	editData := make(map[string]interface{})

	for _, value := range documentLink.Data {
		if value["Location"] != nil {
			if strings.Contains(value["Location"].(string), "coa_report") {
				editData["coa_report_document_link"] = value["Location"]
			}
		}
	}

	errEdit := tx.Model(&coaReport).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return coaReport, errEdit
	}

	var history History

	history.CoaReportId = &coaReport.ID
	history.UserId = userId
	history.Status = fmt.Sprintf("Update upload document coa report with id = %v", coaReport.ID)
	history.IupopkId = &coaReport.IupopkId
	dataInput, _ := json.Marshal(documentLink)
	history.AfterData = dataInput

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return coaReport, createHistoryErr
	}

	tx.Commit()
	return coaReport, nil
}

// RKAB

func (r *repository) CreateRkab(input rkab.RkabInput, iupopkId int, userId uint) (rkab.Rkab, error) {
	var createdRkab rkab.Rkab
	var iup iupopk.Iupopk
	var findRkab []rkab.Rkab

	var counterTransaction counter.Counter

	tx := r.db.Begin()

	errFindIup := tx.Where("id = ?", iupopkId).First(&iup).Error

	if errFindIup != nil {
		tx.Rollback()
		return createdRkab, errFindIup
	}

	errFindRkab := tx.Where("year = ?", input.Year).Find(&findRkab).Error

	if errFindRkab != nil {
		tx.Rollback()
		return createdRkab, nil
	}

	errFindCounterTransaction := tx.Where("iupopk_id = ?", iupopkId).First(&counterTransaction).Error

	if errFindCounterTransaction != nil {
		tx.Rollback()
		return createdRkab, errFindCounterTransaction
	}

	idNumber := "RKB-" + iup.Code
	createdRkab.IdNumber = createIdNumber(idNumber, uint(counterTransaction.Rkab))
	createdRkab.LetterNumber = strings.ToUpper(input.LetterNumber)
	createdRkab.DateOfIssue = input.DateOfIssue
	createdRkab.Year = input.Year
	createdRkab.ProductionQuota = input.ProductionQuota
	createdRkab.IupopkId = uint(iupopkId)

	if len(findRkab) > 0 {
		createdRkab.IsRevision = true

		errUpd := tx.Model(&findRkab).Where("year = ? AND iupopk_id = ?", input.Year, iupopkId).Update("is_revision", true).Error

		if errUpd != nil {
			return createdRkab, errUpd
		}
	}

	errCreate := tx.Create(&createdRkab).Error

	if errCreate != nil {
		tx.Rollback()
		return createdRkab, errCreate
	}

	rkabData, _ := json.Marshal(createdRkab)

	var history History

	history.RkabId = &createdRkab.ID
	history.Status = "Created Rkab"
	history.UserId = userId
	history.IupopkId = &iup.ID
	history.BeforeData = rkabData
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdRkab, createHistoryErr
	}

	updateCounterErr := tx.Model(&counterTransaction).Where("iupopk_id = ?", iupopkId).Update("rkab", counterTransaction.Rkab+1).Error

	if updateCounterErr != nil {
		tx.Rollback()
		return createdRkab, updateCounterErr
	}

	tx.Commit()
	return createdRkab, nil
}

func (r *repository) DeleteRkab(id int, iupopkId int, userId uint) (bool, error) {
	var rkabDeleted rkab.Rkab
	var iup iupopk.Iupopk

	var rkabYear string

	tx := r.db.Begin()

	errFindIup := tx.Where("id = ?", iupopkId).First(&iup).Error

	if errFindIup != nil {
		tx.Rollback()
		return false, errFindIup
	}

	errFind := tx.Where("id = ?", id).First(&rkabDeleted).Error

	if errFind != nil {
		tx.Rollback()
		return false, errFind
	}

	rkabYear = rkabDeleted.Year
	rkabData, _ := json.Marshal(rkabDeleted)

	errDelete := tx.Unscoped().Where("id = ? AND iupopk_id = ?", id, iupopkId).Delete(&rkabDeleted).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var listRkab []rkab.Rkab

	errFindRkabWithYear := tx.Where("year = ?", rkabYear).Order("created_at desc").Find(&listRkab).Error

	if errFindRkabWithYear == nil {
		var newUpdRkab rkab.Rkab
		var isRevision = true
		if len(listRkab) > 0 {
			if len(listRkab) == 1 {
				isRevision = false
			}
			newUpdRkab = listRkab[0]

			errUpd := tx.Model(&newUpdRkab).Update("is_revision", isRevision).Error

			if errUpd != nil {
				tx.Rollback()
				return false, errUpd
			}
		}
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Rkab with id number %s and id %v", &rkabDeleted.IdNumber, rkabDeleted.ID)
	history.UserId = userId
	history.IupopkId = &iup.ID
	history.BeforeData = rkabData
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UploadDocumentRkab(id uint, urlS3 string, userId uint, iupopkId int) (rkab.Rkab, error) {
	var rkab rkab.Rkab

	tx := r.db.Begin()

	errFind := tx.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&rkab).Error

	if errFind != nil {
		return rkab, errFind
	}

	errEdit := tx.Model(&rkab).Update("rkab_document_link", urlS3).Error

	if errEdit != nil {
		tx.Rollback()
		return rkab, errEdit
	}

	var history History

	history.RkabId = &rkab.ID
	history.UserId = userId
	history.Status = "Uploaded document rkab"

	history.IupopkId = &rkab.IupopkId
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return rkab, createHistoryErr
	}

	tx.Commit()
	return rkab, nil
}
