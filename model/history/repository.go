package history

import (
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/minerba"
	"ajebackend/model/production"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/vessel"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint, transactionType string) (bool, error)
	UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string) (transaction.Transaction, error)
	CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	CreateMinerba(period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error)
	UpdateMinerba(id int, updateTransaction []int, userId uint) (minerba.Minerba, error)
	DeleteMinerba(idMinerba int, userId uint) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error)
	CreateDmo(dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error)
	DeleteDmo(idDmo int, userId uint) (bool, error)
	UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint) (dmo.Dmo, error)
	UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error)
	UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string) (dmo.Dmo, error)
	UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error)
	CreateProduction(input production.InputCreateProduction, userId uint) (production.Production, error)
	UpdateProduction(input production.InputCreateProduction, productionId int, userId uint) (production.Production, error)
	DeleteProduction(productionId int, userId uint) (bool, error)
	CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error)
	EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error)
	UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string) (groupingvesselln.GroupingVesselLn, error)
	DeleteGroupingVesselLn(id int, userId uint) (bool, error)
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

	idNumber := fmt.Sprintf("%s-%v-%v-%v", model, monthNumber, year, helper.CreateIdNumber(int(id)))

	return idNumber
}

// Transaction

func (r *repository) CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var createdTransaction transaction.Transaction

	tx := r.db.Begin()

	createdTransaction.DmoId = nil
	createdTransaction.TransactionType = "DN"
	if inputTransactionDN.Seller == "" {
		createdTransaction.Seller = "PT Angsana Jaya Energi"
	} else {
		createdTransaction.Seller = inputTransactionDN.Seller
	}
	createdTransaction.ShippingDate = inputTransactionDN.ShippingDate
	createdTransaction.Quantity = inputTransactionDN.Quantity
	createdTransaction.TugboatName = inputTransactionDN.TugboatName
	createdTransaction.BargeName = inputTransactionDN.BargeName
	createdTransaction.VesselName = inputTransactionDN.VesselName
	createdTransaction.CustomerName = inputTransactionDN.CustomerName
	createdTransaction.LoadingPortName = inputTransactionDN.LoadingPortName
	createdTransaction.LoadingPortLocation = inputTransactionDN.LoadingPortLocation
	createdTransaction.UnloadingPortName = inputTransactionDN.UnloadingPortName
	createdTransaction.UnloadingPortLocation = inputTransactionDN.UnloadingPortLocation
	createdTransaction.DmoDestinationPort = inputTransactionDN.DmoDestinationPort
	createdTransaction.SkbDate = inputTransactionDN.SkbDate
	createdTransaction.SkbNumber = strings.ToUpper(inputTransactionDN.SkbNumber)
	createdTransaction.SkabDate = inputTransactionDN.SkabDate
	createdTransaction.SkabNumber = strings.ToUpper(inputTransactionDN.SkabNumber)
	createdTransaction.BillOfLadingDate = inputTransactionDN.BillOfLadingDate
	createdTransaction.BillOfLadingNumber = strings.ToUpper(inputTransactionDN.BillOfLadingNumber)
	createdTransaction.RoyaltyRate = inputTransactionDN.RoyaltyRate
	createdTransaction.DpRoyaltyPrice = inputTransactionDN.DpRoyaltyPrice
	createdTransaction.DpRoyaltyCurrency = strings.ToUpper(inputTransactionDN.DpRoyaltyCurrency)
	if inputTransactionDN.DpRoyaltyCurrency == "" {
		createdTransaction.DpRoyaltyCurrency = "IDR"
	}
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
	createdTransaction.PaymentDpRoyaltyCurrency = strings.ToUpper(inputTransactionDN.PaymentDpRoyaltyCurrency)
	if inputTransactionDN.PaymentDpRoyaltyCurrency == "" {
		createdTransaction.PaymentDpRoyaltyCurrency = "IDR"
	}

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
	createdTransaction.SurveyorName = inputTransactionDN.SurveyorName
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
	createdTransaction.SalesSystem = inputTransactionDN.SalesSystem
	createdTransaction.InvoiceDate = inputTransactionDN.InvoiceDate
	createdTransaction.InvoiceNumber = strings.ToUpper(inputTransactionDN.InvoiceNumber)
	createdTransaction.InvoicePriceUnit = inputTransactionDN.InvoicePriceUnit
	createdTransaction.InvoicePriceTotal = inputTransactionDN.InvoicePriceTotal
	createdTransaction.DmoReconciliationLetter = inputTransactionDN.DmoReconciliationLetter
	createdTransaction.ContractDate = inputTransactionDN.ContractDate
	createdTransaction.ContractNumber = strings.ToUpper(inputTransactionDN.ContractNumber)
	createdTransaction.DmoBuyerName = inputTransactionDN.DmoBuyerName
	createdTransaction.DmoIndustryType = inputTransactionDN.DmoIndustryType
	createdTransaction.DmoCategory = strings.ToUpper(inputTransactionDN.DmoCategory)
	createdTransaction.IsCoaFinish = inputTransactionDN.IsCoaFinish
	createdTransaction.IsRoyaltyFinalFinish = inputTransactionDN.IsRoyaltyFinalFinish
	createTransactionErr := tx.Create(&createdTransaction).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, createTransactionErr
	}

	idNumber := createIdNumber("DN", createdTransaction.ID)

	updateTransactionsErr := tx.Model(&createdTransaction).Update("id_number", idNumber).Error

	if updateTransactionsErr != nil {
		tx.Rollback()
		return createdTransaction, updateTransactionsErr
	}

	var history History

	history.TransactionId = &createdTransaction.ID
	history.Status = "Created"
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdTransaction, createHistoryErr
	}

	tx.Commit()
	return createdTransaction, createHistoryErr
}

func (r *repository) DeleteTransaction(id int, userId uint, transactionType string) (bool, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ? AND transaction_type = ?", id, transactionType).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := tx.Unscoped().Where("id = ? AND transaction_type = ?", id, transactionType).Delete(&transaction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Transaction %v with id number %s and id %v", transactionType, *transaction.IdNumber, transaction.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, createHistoryErr
}

func (r *repository) UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ?", idTransaction).First(&transaction).Error

	if errFind != nil {
		tx.Rollback()
		return transaction, errFind
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorBeforeDataJsonMarshal
	}

	inputEditTransactionDN.Seller = inputEditTransactionDN.Seller
	inputEditTransactionDN.TugboatName = inputEditTransactionDN.TugboatName
	inputEditTransactionDN.BargeName = inputEditTransactionDN.BargeName
	inputEditTransactionDN.VesselName = inputEditTransactionDN.VesselName
	inputEditTransactionDN.CustomerName = inputEditTransactionDN.CustomerName
	inputEditTransactionDN.LoadingPortName = inputEditTransactionDN.LoadingPortName
	inputEditTransactionDN.LoadingPortLocation = inputEditTransactionDN.LoadingPortLocation
	inputEditTransactionDN.UnloadingPortName = inputEditTransactionDN.UnloadingPortName
	inputEditTransactionDN.UnloadingPortLocation = inputEditTransactionDN.UnloadingPortLocation
	inputEditTransactionDN.DmoDestinationPort = inputEditTransactionDN.DmoDestinationPort
	inputEditTransactionDN.SkbNumber = strings.ToUpper(inputEditTransactionDN.SkbNumber)
	inputEditTransactionDN.SkabNumber = strings.ToUpper(inputEditTransactionDN.SkabNumber)
	inputEditTransactionDN.BillOfLadingNumber = strings.ToUpper(inputEditTransactionDN.BillOfLadingNumber)
	inputEditTransactionDN.DpRoyaltyCurrency = strings.ToUpper(inputEditTransactionDN.DpRoyaltyCurrency)
	if inputEditTransactionDN.DpRoyaltyNtpn != nil {
		dpNtpn := strings.ToUpper(*inputEditTransactionDN.DpRoyaltyNtpn)
		inputEditTransactionDN.DpRoyaltyNtpn = &dpNtpn
	}

	if inputEditTransactionDN.DpRoyaltyBillingCode != nil {
		dpBillingCode := strings.ToUpper(*inputEditTransactionDN.DpRoyaltyBillingCode)
		inputEditTransactionDN.DpRoyaltyBillingCode = &dpBillingCode
	}

	inputEditTransactionDN.PaymentDpRoyaltyCurrency = strings.ToUpper(inputEditTransactionDN.PaymentDpRoyaltyCurrency)
	inputEditTransactionDN.PaymentDpRoyaltyCurrency = strings.ToUpper(inputEditTransactionDN.PaymentDpRoyaltyCurrency)

	if inputEditTransactionDN.PaymentDpRoyaltyNtpn != nil {
		paymentDpNtpn := strings.ToUpper(*inputEditTransactionDN.PaymentDpRoyaltyNtpn)
		inputEditTransactionDN.PaymentDpRoyaltyNtpn = &paymentDpNtpn
	}

	if inputEditTransactionDN.PaymentDpRoyaltyBillingCode != nil {
		paymentDpBillingCode := strings.ToUpper(*inputEditTransactionDN.PaymentDpRoyaltyBillingCode)
		inputEditTransactionDN.PaymentDpRoyaltyBillingCode = &paymentDpBillingCode
	}

	inputEditTransactionDN.LhvNumber = strings.ToUpper(inputEditTransactionDN.LhvNumber)
	inputEditTransactionDN.SurveyorName = inputEditTransactionDN.SurveyorName
	inputEditTransactionDN.CowNumber = strings.ToUpper(inputEditTransactionDN.CowNumber)
	inputEditTransactionDN.CoaNumber = strings.ToUpper(inputEditTransactionDN.CoaNumber)
	inputEditTransactionDN.SalesSystem = inputEditTransactionDN.SalesSystem
	inputEditTransactionDN.InvoiceNumber = strings.ToUpper(inputEditTransactionDN.InvoiceNumber)
	inputEditTransactionDN.ContractNumber = strings.ToUpper(inputEditTransactionDN.ContractNumber)
	inputEditTransactionDN.DmoBuyerName = inputEditTransactionDN.DmoBuyerName
	inputEditTransactionDN.DmoIndustryType = inputEditTransactionDN.DmoIndustryType
	inputEditTransactionDN.DmoCategory = strings.ToUpper(inputEditTransactionDN.DmoCategory)

	dataInput, errorMarshal := json.Marshal(inputEditTransactionDN)

	if errorMarshal != nil {
		tx.Rollback()
		return transaction, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return transaction, createHistoryErr
	}

	tx.Commit()
	return transaction, nil
}

func (r *repository) UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string) (transaction.Transaction, error) {
	var uploadedTransaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ? AND transaction_type = ?", idTransaction, transactionType).First(&uploadedTransaction).Error

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return uploadedTransaction, createHistoryErr
	}

	tx.Commit()
	return uploadedTransaction, nil
}

// Transaction LN

func (r *repository) CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var createdTransactionLn transaction.Transaction

	tx := r.db.Begin()

	createdTransactionLn.TransactionType = "LN"
	if inputTransactionLN.Seller == "" {
		createdTransactionLn.Seller = "PT Angsana Jaya Energi"
	} else {
		createdTransactionLn.Seller = inputTransactionLN.Seller
	}
	createdTransactionLn.ShippingDate = inputTransactionLN.ShippingDate
	createdTransactionLn.Quantity = inputTransactionLN.Quantity
	createdTransactionLn.TugboatName = inputTransactionLN.TugboatName
	createdTransactionLn.BargeName = inputTransactionLN.BargeName
	createdTransactionLn.VesselName = inputTransactionLN.VesselName
	createdTransactionLn.CustomerName = inputTransactionLN.CustomerName
	createdTransactionLn.LoadingPortName = inputTransactionLN.LoadingPortName
	createdTransactionLn.LoadingPortLocation = inputTransactionLN.LoadingPortLocation
	createdTransactionLn.UnloadingPortName = inputTransactionLN.UnloadingPortName
	createdTransactionLn.UnloadingPortLocation = inputTransactionLN.UnloadingPortLocation
	createdTransactionLn.DmoDestinationPort = inputTransactionLN.DmoDestinationPort
	createdTransactionLn.SkbDate = inputTransactionLN.SkbDate
	createdTransactionLn.SkbNumber = strings.ToUpper(inputTransactionLN.SkbNumber)
	createdTransactionLn.SkabDate = inputTransactionLN.SkabDate
	createdTransactionLn.SkabNumber = strings.ToUpper(inputTransactionLN.SkabNumber)
	createdTransactionLn.BillOfLadingDate = inputTransactionLN.BillOfLadingDate
	createdTransactionLn.BillOfLadingNumber = strings.ToUpper(inputTransactionLN.BillOfLadingNumber)
	createdTransactionLn.RoyaltyRate = inputTransactionLN.RoyaltyRate
	createdTransactionLn.DpRoyaltyPrice = inputTransactionLN.DpRoyaltyPrice
	createdTransactionLn.DpRoyaltyCurrency = strings.ToUpper(inputTransactionLN.DpRoyaltyCurrency)
	if inputTransactionLN.DpRoyaltyCurrency == "" {
		createdTransactionLn.DpRoyaltyCurrency = "USD"
	}
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
	createdTransactionLn.PaymentDpRoyaltyCurrency = strings.ToUpper(inputTransactionLN.PaymentDpRoyaltyCurrency)
	if inputTransactionLN.PaymentDpRoyaltyCurrency == "" {
		createdTransactionLn.PaymentDpRoyaltyCurrency = "USD"
	}

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
	createdTransactionLn.SurveyorName = inputTransactionLN.SurveyorName
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
	createdTransactionLn.SalesSystem = inputTransactionLN.SalesSystem
	createdTransactionLn.InvoiceDate = inputTransactionLN.InvoiceDate
	createdTransactionLn.InvoiceNumber = strings.ToUpper(inputTransactionLN.InvoiceNumber)
	createdTransactionLn.InvoicePriceUnit = inputTransactionLN.InvoicePriceUnit
	createdTransactionLn.InvoicePriceTotal = inputTransactionLN.InvoicePriceTotal
	createdTransactionLn.DmoReconciliationLetter = inputTransactionLN.DmoReconciliationLetter
	createdTransactionLn.ContractDate = inputTransactionLN.ContractDate
	createdTransactionLn.ContractNumber = strings.ToUpper(inputTransactionLN.ContractNumber)
	createdTransactionLn.DmoBuyerName = inputTransactionLN.DmoBuyerName
	createdTransactionLn.DmoIndustryType = inputTransactionLN.DmoIndustryType
	createdTransactionLn.DmoCategory = strings.ToUpper(inputTransactionLN.DmoCategory)

	createTransactionErr := tx.Create(&createdTransactionLn).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransactionLn, createTransactionErr
	}

	idNumber := createIdNumber("LN", createdTransactionLn.ID)

	updateTransactionsErr := tx.Model(&createdTransactionLn).Update("id_number", idNumber).Error

	if updateTransactionsErr != nil {
		tx.Rollback()
		return createdTransactionLn, updateTransactionsErr
	}

	if inputTransactionLN.VesselName != "" {
		var createVesselMaster vessel.Vessel

		createVesselMaster.Name = inputTransactionLN.VesselName

		errCreateVesselMaster := tx.FirstOrCreate(&createVesselMaster, createVesselMaster).Error

		if errCreateVesselMaster != nil {
			tx.Rollback()
			return createdTransactionLn, errCreateVesselMaster
		}
	}

	var history History

	history.TransactionId = &createdTransactionLn.ID
	history.Status = "Created LN"
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdTransactionLn, createHistoryErr
	}

	tx.Commit()
	return createdTransactionLn, createHistoryErr
}

func (r *repository) UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var updatedTransactionLn transaction.Transaction
	var updateTransaction transaction.DataTransactionInput

	tx := r.db.Begin()

	errFind := r.db.Where("id = ? AND transaction_type = ?", id, "LN").First(&updatedTransactionLn).Error

	if errFind != nil {
		tx.Rollback()
		return updatedTransactionLn, errFind
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(updatedTransactionLn)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedTransactionLn, errorBeforeDataJsonMarshal
	}

	if inputTransactionLN.Seller == "" {
		updateTransaction.Seller = "PT Angsana Jaya Energi"
	} else {
		updateTransaction.Seller = inputTransactionLN.Seller
	}
	updateTransaction.ShippingDate = inputTransactionLN.ShippingDate
	updateTransaction.Quantity = inputTransactionLN.Quantity
	updateTransaction.TugboatName = inputTransactionLN.TugboatName
	updateTransaction.BargeName = inputTransactionLN.BargeName
	updateTransaction.VesselName = inputTransactionLN.VesselName
	updateTransaction.CustomerName = inputTransactionLN.CustomerName
	updateTransaction.LoadingPortName = inputTransactionLN.LoadingPortName
	updateTransaction.LoadingPortLocation = inputTransactionLN.LoadingPortLocation
	updateTransaction.UnloadingPortName = inputTransactionLN.UnloadingPortName
	updateTransaction.UnloadingPortLocation = inputTransactionLN.UnloadingPortLocation
	updateTransaction.DmoDestinationPort = inputTransactionLN.DmoDestinationPort
	updateTransaction.SkbDate = inputTransactionLN.SkbDate
	updateTransaction.SkbNumber = strings.ToUpper(inputTransactionLN.SkbNumber)
	updateTransaction.SkabDate = inputTransactionLN.SkabDate
	updateTransaction.SkabNumber = strings.ToUpper(inputTransactionLN.SkabNumber)
	updateTransaction.BillOfLadingDate = inputTransactionLN.BillOfLadingDate
	updateTransaction.BillOfLadingNumber = strings.ToUpper(inputTransactionLN.BillOfLadingNumber)
	updateTransaction.RoyaltyRate = inputTransactionLN.RoyaltyRate
	updateTransaction.DpRoyaltyPrice = inputTransactionLN.DpRoyaltyPrice
	updateTransaction.DpRoyaltyCurrency = strings.ToUpper(inputTransactionLN.DpRoyaltyCurrency)
	if inputTransactionLN.DpRoyaltyCurrency == "" {
		updateTransaction.DpRoyaltyCurrency = "USD"
	}
	updateTransaction.DpRoyaltyDate = inputTransactionLN.DpRoyaltyDate
	if inputTransactionLN.DpRoyaltyNtpn != nil {
		dpNtpn := strings.ToUpper(*inputTransactionLN.DpRoyaltyNtpn)
		updateTransaction.DpRoyaltyNtpn = &dpNtpn
	}

	if inputTransactionLN.DpRoyaltyBillingCode != nil {
		dpBillingCode := strings.ToUpper(*inputTransactionLN.DpRoyaltyBillingCode)
		updateTransaction.DpRoyaltyBillingCode = &dpBillingCode
	}

	updateTransaction.DpRoyaltyTotal = inputTransactionLN.DpRoyaltyTotal
	updateTransaction.PaymentDpRoyaltyPrice = inputTransactionLN.PaymentDpRoyaltyPrice
	updateTransaction.PaymentDpRoyaltyCurrency = strings.ToUpper(inputTransactionLN.PaymentDpRoyaltyCurrency)
	if inputTransactionLN.PaymentDpRoyaltyCurrency == "" {
		updateTransaction.PaymentDpRoyaltyCurrency = "USD"
	}

	updateTransaction.PaymentDpRoyaltyDate = inputTransactionLN.PaymentDpRoyaltyDate
	if inputTransactionLN.PaymentDpRoyaltyNtpn != nil {
		paymentDpNtpn := strings.ToUpper(*inputTransactionLN.PaymentDpRoyaltyNtpn)
		updateTransaction.PaymentDpRoyaltyNtpn = &paymentDpNtpn
	}

	if inputTransactionLN.PaymentDpRoyaltyBillingCode != nil {
		paymentDpBillingCode := strings.ToUpper(*inputTransactionLN.PaymentDpRoyaltyBillingCode)
		updateTransaction.PaymentDpRoyaltyBillingCode = &paymentDpBillingCode
	}

	updateTransaction.PaymentDpRoyaltyTotal = inputTransactionLN.PaymentDpRoyaltyTotal
	updateTransaction.LhvDate = inputTransactionLN.LhvDate
	updateTransaction.LhvNumber = strings.ToUpper(inputTransactionLN.LhvNumber)
	updateTransaction.SurveyorName = inputTransactionLN.SurveyorName
	updateTransaction.CowDate = inputTransactionLN.CowDate
	updateTransaction.CowNumber = strings.ToUpper(inputTransactionLN.CowNumber)
	updateTransaction.CoaDate = inputTransactionLN.CoaDate
	updateTransaction.CoaNumber = strings.ToUpper(inputTransactionLN.CoaNumber)
	updateTransaction.QualityTmAr = inputTransactionLN.QualityTmAr
	updateTransaction.QualityImAdb = inputTransactionLN.QualityImAdb
	updateTransaction.QualityAshAr = inputTransactionLN.QualityAshAr
	updateTransaction.QualityAshAdb = inputTransactionLN.QualityAshAdb
	updateTransaction.QualityVmAdb = inputTransactionLN.QualityVmAdb
	updateTransaction.QualityFcAdb = inputTransactionLN.QualityFcAdb
	updateTransaction.QualityTsAr = inputTransactionLN.QualityTsAr
	updateTransaction.QualityTsAdb = inputTransactionLN.QualityTsAdb
	updateTransaction.QualityCaloriesAr = inputTransactionLN.QualityCaloriesAr
	updateTransaction.QualityCaloriesAdb = inputTransactionLN.QualityCaloriesAdb
	updateTransaction.BargingDistance = inputTransactionLN.BargingDistance
	updateTransaction.SalesSystem = inputTransactionLN.SalesSystem
	updateTransaction.InvoiceDate = inputTransactionLN.InvoiceDate
	updateTransaction.InvoiceNumber = strings.ToUpper(inputTransactionLN.InvoiceNumber)
	updateTransaction.InvoicePriceUnit = inputTransactionLN.InvoicePriceUnit
	updateTransaction.InvoicePriceTotal = inputTransactionLN.InvoicePriceTotal
	updateTransaction.DmoReconciliationLetter = inputTransactionLN.DmoReconciliationLetter
	updateTransaction.ContractDate = inputTransactionLN.ContractDate
	updateTransaction.ContractNumber = strings.ToUpper(inputTransactionLN.ContractNumber)
	updateTransaction.DmoBuyerName = inputTransactionLN.DmoBuyerName
	updateTransaction.DmoIndustryType = inputTransactionLN.DmoIndustryType
	updateTransaction.DmoCategory = strings.ToUpper(inputTransactionLN.DmoCategory)
	updateTransaction.IsFinanceCheck = inputTransactionLN.IsFinanceCheck
	updateTransaction.IsNotClaim = inputTransactionLN.IsNotClaim

	editTransaction, errorMarshal := json.Marshal(updateTransaction)

	if errorMarshal != nil {
		tx.Rollback()
		return updatedTransactionLn, errorMarshal
	}

	var editTransactionInput map[string]interface{}

	errorUnmarshalTransaction := json.Unmarshal(editTransaction, &editTransactionInput)

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

	if inputTransactionLN.VesselName != "" {
		var createVesselMaster vessel.Vessel

		createVesselMaster.Name = inputTransactionLN.VesselName

		errCreateVesselMaster := tx.FirstOrCreate(&createVesselMaster, createVesselMaster).Error

		if errCreateVesselMaster != nil {
			tx.Rollback()
			return updatedTransactionLn, errCreateVesselMaster
		}
	}

	var history History

	history.TransactionId = &updatedTransactionLn.ID
	history.Status = "Updated LN"
	history.UserId = userId
	history.BeforeData = beforeData
	history.AfterData = afterData

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedTransactionLn, createHistoryErr
	}

	tx.Commit()
	return updatedTransactionLn, createHistoryErr
}

// Minerba

func (r *repository) CreateMinerba(period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error) {
	var createdMinerba minerba.Minerba

	tx := r.db.Begin()
	createdMinerba.Period = period
	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ?", updateTransaction).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdMinerba, findTransactionsErr
	}

	if len(transactions) != len(updateTransaction) {
		tx.Rollback()
		return createdMinerba, errors.New("please check some of transactions not found")
	}

	var tempQuantity float64
	for _, v := range transactions {
		tempQuantity += v.Quantity
	}

	stringTempQuantity := fmt.Sprintf("%.3f", tempQuantity)
	parseTempQuantity, _ := strconv.ParseFloat(stringTempQuantity, 64)

	createdMinerba.Quantity = parseTempQuantity

	errCreateMinerba := tx.Create(&createdMinerba).Error

	if errCreateMinerba != nil {
		tx.Rollback()
		return createdMinerba, errCreateMinerba
	}

	idNumber := baseIdNumber + "-" + helper.CreateIdNumber(int(createdMinerba.ID))

	updateMinerbaErr := tx.Model(&createdMinerba).Update("id_number", idNumber).Error

	if updateMinerbaErr != nil {
		tx.Rollback()
		return createdMinerba, updateMinerbaErr
	}

	updateTransactionErr := tx.Table("transactions").Where("id IN ?", updateTransaction).Update("minerba_id", createdMinerba.ID).Error

	if updateTransactionErr != nil {
		tx.Rollback()
		return createdMinerba, updateTransactionErr
	}

	var history History

	history.MinerbaId = &createdMinerba.ID
	history.Status = "Created Minerba Report"
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdMinerba, createHistoryErr
	}

	tx.Commit()
	return createdMinerba, nil
}

func (r *repository) UpdateMinerba(id int, updateTransaction []int, userId uint) (minerba.Minerba, error) {
	var updatedMinerba minerba.Minerba
	var quantityMinerba float64

	historyBefore := make(map[string]interface{})
	historyAfter := make(map[string]interface{})
	tx := r.db.Begin()

	findMinerbaErr := tx.Where("id = ?", id).First(&updatedMinerba).Error

	if findMinerbaErr != nil {
		return updatedMinerba, findMinerbaErr
	}

	historyBefore["minerba"] = updatedMinerba

	var beforeTransaction []transaction.Transaction
	findTransactionBeforeErr := tx.Where("minerba_id = ?", id).Find(&beforeTransaction).Error

	if findTransactionBeforeErr != nil {
		return updatedMinerba, findTransactionBeforeErr
	}

	var transactionBefore []uint

	for _, v := range beforeTransaction {
		transactionBefore = append(transactionBefore, v.ID)
	}

	historyBefore["transactions"] = transactionBefore

	errUpdMinerbaNil := tx.Model(&beforeTransaction).Where("minerba_id = ?", id).Update("minerba_id", nil).Error

	if errUpdMinerbaNil != nil {
		tx.Rollback()
		return updatedMinerba, errUpdMinerbaNil
	}

	var transactions []transaction.Transaction
	findTransactionsErr := tx.Where("id IN ?", updateTransaction).Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return updatedMinerba, findTransactionsErr
	}

	if len(transactions) != len(updateTransaction) {
		tx.Rollback()
		return updatedMinerba, errors.New("please check some of transactions not found")
	}

	for _, v := range transactions {
		quantityMinerba += v.Quantity
	}

	errUpdateMinerba := tx.Model(&updatedMinerba).Update("quantity", quantityMinerba).Error

	if errUpdateMinerba != nil {
		tx.Rollback()
		return updatedMinerba, errUpdateMinerba
	}

	historyAfter["minerba"] = updatedMinerba
	historyAfter["transactions"] = updateTransaction

	updateTransactionErr := tx.Table("transactions").Where("id IN ?", updateTransaction).Update("minerba_id", id).Error

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

	history.MinerbaId = &updatedMinerba.ID
	history.Status = "Updated Minerba Report"
	history.UserId = userId
	history.AfterData = afterData
	history.BeforeData = beforeData
	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedMinerba, createHistoryErr
	}

	tx.Commit()
	return updatedMinerba, nil
}

func (r *repository) DeleteMinerba(idMinerba int, userId uint) (bool, error) {

	tx := r.db.Begin()
	var minerba minerba.Minerba

	findMinerbaErr := tx.Where("id = ?", idMinerba).First(&minerba).Error

	if findMinerbaErr != nil {
		tx.Rollback()
		return false, findMinerbaErr
	}

	errDelete := tx.Unscoped().Where("id = ?", idMinerba).Delete(&minerba).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Minerba Report with id number %s and id %v", *minerba.IdNumber, minerba.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error) {
	tx := r.db.Begin()
	var minerba minerba.Minerba

	errFind := tx.Where("id = ?", id).First(&minerba).Error

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
			if strings.Contains(value["Location"].(string), "recapdmo") {
				editData["recap_dmo_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "detaildmo") {
				editData["detail_dmo_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "sp3meln") {
				editData["sp3_meln_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "inswexport") {
				editData["insw_export_document_link"] = value["Location"]
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

func (r *repository) CreateDmo(dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error) {
	var createdDmo dmo.Dmo
	var transactionBarge []transaction.Transaction
	var transactionVessel []transaction.Transaction

	barge := false
	vessel := false
	tx := r.db.Begin()

	createdDmo.Period = dmoInput.Period
	createdDmo.IsDocumentCustom = dmoInput.IsDocumentCustom
	createdDmo.DocumentDate = dmoInput.DocumentDate
	if len(dmoInput.TransactionBarge) > 0 {
		var bargeQuantity float64
		barge = true
		findTransactionBargeErr := tx.Where("id IN ?", dmoInput.TransactionBarge).Find(&transactionBarge).Error

		if findTransactionBargeErr != nil {
			return createdDmo, findTransactionBargeErr
		}

		for _, v := range transactionBarge {
			bargeQuantity += v.Quantity
		}

		stringBargeQuantity := fmt.Sprintf("%.3f", bargeQuantity)
		parseBargeQuantity, _ := strconv.ParseFloat(stringBargeQuantity, 64)

		createdDmo.BargeTotalQuantity = parseBargeQuantity
		createdDmo.BargeGrandTotalQuantity = parseBargeQuantity
	}

	if len(dmoInput.TransactionVessel) > 0 {
		var vesselQuantity float64
		var vesselAdjustment float64
		vessel = true
		findTransactionVesselErr := tx.Where("id IN ?", dmoInput.TransactionVessel).Find(&transactionVessel).Error

		if findTransactionVesselErr != nil {
			return createdDmo, findTransactionVesselErr
		}

		for _, v := range transactionVessel {
			vesselQuantity += v.Quantity
		}

		for _, v := range dmoInput.VesselAdjustment {
			vesselAdjustment += v.Adjustment
		}

		stringVesselAdjustment := fmt.Sprintf("%.3f", vesselAdjustment)
		parseVesselAdjustment, _ := strconv.ParseFloat(stringVesselAdjustment, 64)

		stringVesselTotalQuantity := fmt.Sprintf("%.3f", vesselQuantity)
		parseVesselTotalQuantity, _ := strconv.ParseFloat(stringVesselTotalQuantity, 64)

		stringVesselGrandTotalQuantity := fmt.Sprintf("%.3f", vesselQuantity+vesselAdjustment)
		parseVesselGrandTotalQuantity, _ := strconv.ParseFloat(stringVesselGrandTotalQuantity, 64)

		createdDmo.VesselAdjustment = parseVesselAdjustment
		createdDmo.VesselTotalQuantity = parseVesselTotalQuantity
		createdDmo.VesselGrandTotalQuantity = parseVesselGrandTotalQuantity
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

	createdDmoErr := tx.Create(&createdDmo).Error

	if createdDmoErr != nil {
		tx.Rollback()
		return createdDmo, createdDmoErr
	}

	idNumber := baseIdNumber + "-" + helper.CreateIdNumber(int(createdDmo.ID))

	updateDmoErr := tx.Model(&createdDmo).Update("id_number", idNumber).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return createdDmo, updateDmoErr
	}

	if len(dmoInput.TransactionBarge) > 0 {
		findTransactionsBargeErr := tx.Where("id IN ?", dmoInput.TransactionBarge).Find(&transactionBarge).Error

		if findTransactionsBargeErr != nil {
			tx.Rollback()
			return createdDmo, findTransactionsBargeErr
		}

		updateTransactionBargeErr := tx.Table("transactions").Where("id IN ?", dmoInput.TransactionBarge).Update("dmo_id", createdDmo.ID).Error

		if updateTransactionBargeErr != nil {
			tx.Rollback()
			return createdDmo, updateTransactionBargeErr
		}
	}

	if len(dmoInput.TransactionVessel) > 0 {
		findTransactionsVesselErr := tx.Where("id IN ?", dmoInput.TransactionVessel).Find(&transactionVessel).Error

		if findTransactionsVesselErr != nil {
			tx.Rollback()
			return createdDmo, findTransactionsVesselErr
		}

		updateTransactionVesselErr := tx.Table("transactions").Where("id IN ?", dmoInput.TransactionVessel).Update("dmo_id", createdDmo.ID).Error

		if updateTransactionVesselErr != nil {
			tx.Rollback()
			return createdDmo, updateTransactionVesselErr
		}

		var dmoVessels []dmovessel.DmoVessel

		if len(dmoInput.VesselAdjustment) > 0 {
			for _, value := range dmoInput.VesselAdjustment {
				var vesselDummy dmovessel.DmoVessel

				stringAdjustment := fmt.Sprintf("%.3f", value.Adjustment)
				parseAdjustment, _ := strconv.ParseFloat(stringAdjustment, 64)

				stringQuantity := fmt.Sprintf("%.3f", value.Quantity)
				parseQuantity, _ := strconv.ParseFloat(stringQuantity, 64)

				stringGrandTotalQuantity := fmt.Sprintf("%.3f", value.Quantity+value.Adjustment)
				parseGrandTotalQuantity, _ := strconv.ParseFloat(stringGrandTotalQuantity, 64)

				vesselDummy.VesselName = value.VesselName
				vesselDummy.Adjustment = parseAdjustment
				vesselDummy.Quantity = parseQuantity
				vesselDummy.DmoId = createdDmo.ID
				vesselDummy.GrandTotalQuantity = parseGrandTotalQuantity
				vesselDummy.BlDate = value.BlDate
				dmoVessels = append(dmoVessels, vesselDummy)
			}

			createDmoVesselsErr := tx.Create(&dmoVessels).Error

			if createDmoVesselsErr != nil {
				tx.Rollback()
				return createdDmo, createDmoVesselsErr
			}
		}
	}

	if len(transactionBarge) != len(dmoInput.TransactionBarge) && len(transactionVessel) != len(dmoInput.TransactionVessel) {
		tx.Rollback()
		return createdDmo, errors.New("please check some of transactions not found")
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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdDmo, createHistoryErr
	}

	tx.Commit()
	return createdDmo, nil
}

func (r *repository) DeleteDmo(idDmo int, userId uint) (bool, error) {

	tx := r.db.Begin()
	var dmo dmo.Dmo

	findDmoErr := tx.Where("id = ?", idDmo).First(&dmo).Error

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

	history.Status = fmt.Sprintf("Deleted Dmo with id number %s and id %v", *dmo.IdNumber, dmo.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

func (r *repository) UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint) (dmo.Dmo, error) {
	tx := r.db.Begin()
	var dmoUpdate dmo.Dmo

	errFind := tx.Where("id = ?", id).First(&dmoUpdate).Error

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

func (r *repository) UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

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

	updateErr := tx.Model(&dmoUpdate).Where("id = ?", id).Update(field, true).Error

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

func (r *repository) UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return dmoUpdate, createHistoryErr
	}

	tx.Commit()
	return dmoUpdate, nil
}

// Production
func (r *repository) CreateProduction(input production.InputCreateProduction, userId uint) (production.Production, error) {
	var createdProduction production.Production

	createdProduction.ProductionDate = input.ProductionDate
	createdProduction.Quantity = input.Quantity

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdProduction, createHistoryErr
	}

	tx.Commit()
	return createdProduction, nil
}

func (r *repository) UpdateProduction(input production.InputCreateProduction, productionId int, userId uint) (production.Production, error) {
	var updatedProduction production.Production

	tx := r.db.Begin()

	errFindProduction := tx.Where("id = ?", productionId).First(&updatedProduction).Error

	if errFindProduction != nil {
		tx.Rollback()
		return updatedProduction, errFindProduction
	}

	beforeData, _ := json.Marshal(updatedProduction)

	editData := make(map[string]interface{})
	editData["production_date"] = input.ProductionDate
	editData["quantity"] = input.Quantity

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

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedProduction, createHistoryErr
	}

	tx.Commit()
	return updatedProduction, nil
}

func (r *repository) DeleteProduction(productionId int, userId uint) (bool, error) {
	tx := r.db.Begin()
	var deletedProduction production.Production

	findDeletedProductionErr := tx.Where("id = ?", productionId).First(&deletedProduction).Error

	if findDeletedProductionErr != nil {
		tx.Rollback()
		return false, findDeletedProductionErr
	}

	errDelete := tx.Unscoped().Where("id = ?", productionId).Delete(&deletedProduction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Production with id %v", deletedProduction.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, nil
}

// Grouping Vessel LN
func (r *repository) CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error) {
	var createdGroupingVesselLn groupingvesselln.GroupingVesselLn
	var transactions []transaction.Transaction
	tx := r.db.Begin()

	findTransactionsErr := tx.Where("id IN ? AND transaction_type = ? AND grouping_vessel_ln_id is NULL", inputGrouping.ListTransactions, "LN").Find(&transactions).Error

	if findTransactionsErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, findTransactionsErr
	}

	if len(transactions) != len(inputGrouping.ListTransactions) {
		tx.Rollback()
		return createdGroupingVesselLn, errors.New("please check some of transactions not found or already created in another group")
	}

	createdGroupingVesselLn.VesselName = inputGrouping.VesselName
	createdGroupingVesselLn.Quantity = inputGrouping.Quantity
	createdGroupingVesselLn.Adjustment = inputGrouping.Adjustment
	createdGroupingVesselLn.GrandTotalQuantity = inputGrouping.GrandTotalQuantity
	createdGroupingVesselLn.DescriptionOfDocumentType = inputGrouping.DescriptionOfDocumentType
	createdGroupingVesselLn.CodeOfDocumentType = inputGrouping.CodeOfDocumentType
	createdGroupingVesselLn.AjuNumber = strings.ToUpper(inputGrouping.AjuNumber)
	createdGroupingVesselLn.PebRegisterNumber = strings.ToUpper(inputGrouping.PebRegisterNumber)
	createdGroupingVesselLn.PebRegisterDate = inputGrouping.PebRegisterDate
	createdGroupingVesselLn.DescriptionOfPabeanOffice = inputGrouping.DescriptionOfPabeanOffice
	createdGroupingVesselLn.CodeOfPabeanOffice = inputGrouping.CodeOfPabeanOffice
	createdGroupingVesselLn.SeriesPebGoods = inputGrouping.SeriesPebGoods
	createdGroupingVesselLn.DescriptionOfGoods = inputGrouping.DescriptionOfGoods
	createdGroupingVesselLn.TarifPosHs = strings.ToUpper(inputGrouping.TarifPosHs)
	createdGroupingVesselLn.PebQuantity = inputGrouping.PebQuantity
	createdGroupingVesselLn.PebUnit = inputGrouping.PebUnit
	createdGroupingVesselLn.ExportValue = inputGrouping.ExportValue
	createdGroupingVesselLn.Currency = inputGrouping.Currency
	createdGroupingVesselLn.LoadingPort = inputGrouping.LoadingPort
	createdGroupingVesselLn.SkaCooNumber = strings.ToUpper(inputGrouping.SkaCooNumber)
	createdGroupingVesselLn.SkaCooDate = inputGrouping.SkaCooDate
	createdGroupingVesselLn.DestinationCountry = inputGrouping.DestinationCountry
	createdGroupingVesselLn.CodeOfDestinationCountry = inputGrouping.CodeOfDestinationCountry
	createdGroupingVesselLn.LsExportNumber = strings.ToUpper(inputGrouping.LsExportNumber)
	createdGroupingVesselLn.LsExportDate = inputGrouping.LsExportDate
	createdGroupingVesselLn.InsuranceCompanyName = inputGrouping.InsuranceCompanyName
	createdGroupingVesselLn.PolisNumber = strings.ToUpper(inputGrouping.PolisNumber)
	createdGroupingVesselLn.NavyCompanyName = inputGrouping.NavyCompanyName
	createdGroupingVesselLn.NavyShipName = inputGrouping.NavyShipName
	createdGroupingVesselLn.NavyImoNumber = strings.ToUpper(inputGrouping.NavyImoNumber)
	createdGroupingVesselLn.Deadweight = inputGrouping.Deadweight

	errCreatedGroupingVesselLn := tx.Create(&createdGroupingVesselLn).Error

	if errCreatedGroupingVesselLn != nil {
		tx.Rollback()
		return createdGroupingVesselLn, errCreatedGroupingVesselLn
	}

	idNumber := createIdNumber("GML", createdGroupingVesselLn.ID)

	updatedGroupingVesselLnErr := tx.Model(&createdGroupingVesselLn).Update("id_number", idNumber).Error

	if updatedGroupingVesselLnErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, updatedGroupingVesselLnErr
	}

	updateTransactions := make(map[string]interface{})

	updateTransactions["vessel_name"] = inputGrouping.VesselName
	updateTransactions["grouping_vessel_ln_id"] = createdGroupingVesselLn.ID

	errUpdateTransactions := tx.Table("transactions").Where("id IN ?", inputGrouping.ListTransactions).Updates(updateTransactions).Error

	if errUpdateTransactions != nil {
		tx.Rollback()
		return createdGroupingVesselLn, errUpdateTransactions
	}

	var history History

	history.GroupingVesselLnId = &createdGroupingVesselLn.ID
	history.Status = "Created Grouping Vessel LN"
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdGroupingVesselLn, createHistoryErr
	}

	tx.Commit()

	return createdGroupingVesselLn, nil
}

func (r *repository) EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error) {
	var updatedGroupingVesselLn groupingvesselln.GroupingVesselLn
	tx := r.db.Begin()

	errFind := r.db.Where("id = ?", id).First(&updatedGroupingVesselLn).Error

	if errFind != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errFind
	}

	if updatedGroupingVesselLn.VesselName != editGrouping.VesselName {
		var transactions []transaction.Transaction
		var listIdTransaction []uint
		errFindTransaction := tx.Where("grouping_vessel_ln_id = ?", id).Find(&transactions).Error

		if errFindTransaction != nil {
			return updatedGroupingVesselLn, errFindTransaction
		}

		for _, trans := range transactions {
			listIdTransaction = append(listIdTransaction, trans.ID)
		}

		errUpdateTransaction := tx.Table("transactions").Where("id IN ?", listIdTransaction).Update("vessel_name", editGrouping.VesselName).Error

		if errUpdateTransaction != nil {
			tx.Rollback()
			return updatedGroupingVesselLn, errUpdateTransaction
		}
	}

	beforeData, errorBeforeDataJsonMarshal := json.Marshal(updatedGroupingVesselLn)

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

	updateGroupingVesselErr := tx.Model(&updatedGroupingVesselLn).Updates(editGroupingVesselLnInput).Error

	if updateGroupingVesselErr != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, updateGroupingVesselErr
	}

	afterData, errorAfterDataJsonMarshal := json.Marshal(updatedGroupingVesselLn)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, errorAfterDataJsonMarshal
	}

	var history History

	history.GroupingVesselLnId = &updatedGroupingVesselLn.ID
	history.Status = "Updated Grouping Vessel LN"
	history.UserId = userId
	history.BeforeData = beforeData
	history.AfterData = afterData

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return updatedGroupingVesselLn, createHistoryErr
	}

	tx.Commit()
	return updatedGroupingVesselLn, nil
}

func (r *repository) UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string) (groupingvesselln.GroupingVesselLn, error) {
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

func (r *repository) DeleteGroupingVesselLn(id int, userId uint) (bool, error) {

	tx := r.db.Begin()
	var groupingVesselLn groupingvesselln.GroupingVesselLn

	findGroupingVesselLnErr := tx.Where("id = ?", id).First(&groupingVesselLn).Error

	if findGroupingVesselLnErr != nil {
		tx.Rollback()
		return false, findGroupingVesselLnErr
	}

	errDelete := tx.Unscoped().Where("id = ?", id).Delete(&groupingVesselLn).Error

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
