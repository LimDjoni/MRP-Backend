package history

import (
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/minerba"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransactionDN(id int, userId uint) (bool, error)
	UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocumentTransactionDN (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error)
	CreateMinerba (period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error)
	UpdateMinerba (id int, updateTransaction []int, userId uint) (minerba.Minerba, error)
	DeleteMinerba (idMinerba int, userId uint) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error)
	CreateDmo (dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error)
	DeleteDmo (idDmo int, userId uint) (bool, error)
	UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint) (dmo.Dmo, error)
	UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint) (dmo.Dmo, error)
	UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint, location string) (dmo.Dmo, error)
	UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint) (dmo.Dmo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func createIdNumber(model string, id uint) string{
	year, month, _ := time.Now().Date()

	monthNumber := strconv.Itoa(int(month))

	if len([]rune(monthNumber)) < 2 {
		monthNumber = "0" + monthNumber
	}

	idNumber := fmt.Sprintf("%s-%v-%v-%v", model, monthNumber, year, helper.CreateIdNumber(int(id)))

	return idNumber
}

// Transaction

func (r *repository) CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var createdTransaction transaction.Transaction

	tx := r.db.Begin()

	createdTransaction.DmoId = nil
	createdTransaction.TransactionType = "DN"
	if inputTransactionDN.Seller == "" {
		createdTransaction.Seller = "PT ANGSANA JAYA ENERGI"
	} else {
		createdTransaction.Seller = strings.ToUpper(inputTransactionDN.Seller)
	}
	createdTransaction.ShippingDate = inputTransactionDN.ShippingDate
	createdTransaction.Quantity = inputTransactionDN.Quantity
	createdTransaction.TugboatName = strings.ToUpper(inputTransactionDN.TugboatName)
	createdTransaction.BargeName = strings.ToUpper(inputTransactionDN.BargeName)
	createdTransaction.VesselName = strings.ToUpper(inputTransactionDN.VesselName)
	createdTransaction.CustomerName = strings.ToUpper(inputTransactionDN.CustomerName)
	createdTransaction.LoadingPortName = strings.ToUpper(inputTransactionDN.LoadingPortName)
	createdTransaction.LoadingPortLocation = strings.ToUpper(inputTransactionDN.LoadingPortLocation)
	createdTransaction.UnloadingPortName = strings.ToUpper(inputTransactionDN.UnloadingPortName)
	createdTransaction.UnloadingPortLocation = strings.ToUpper(inputTransactionDN.UnloadingPortLocation)
	createdTransaction.DmoDestinationPort = strings.ToUpper(inputTransactionDN.DmoDestinationPort)
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
	createdTransaction.SurveyorName = strings.ToUpper(inputTransactionDN.SurveyorName)
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
	createdTransaction.SalesSystem = strings.ToUpper(inputTransactionDN.SalesSystem)
	createdTransaction.InvoiceDate = inputTransactionDN.InvoiceDate
	createdTransaction.InvoiceNumber = strings.ToUpper(inputTransactionDN.InvoiceNumber)
	createdTransaction.InvoicePriceUnit = inputTransactionDN.InvoicePriceUnit
	createdTransaction.InvoicePriceTotal = inputTransactionDN.InvoicePriceTotal
	createdTransaction.DmoReconciliationLetter = inputTransactionDN.DmoReconciliationLetter
	createdTransaction.ContractDate = inputTransactionDN.ContractDate
	createdTransaction.ContractNumber = strings.ToUpper(inputTransactionDN.ContractNumber)
	createdTransaction.DmoBuyerName = strings.ToUpper(inputTransactionDN.DmoBuyerName)
	createdTransaction.DmoIndustryType = strings.ToUpper(inputTransactionDN.DmoIndustryType)
	createdTransaction.DmoCategory = strings.ToUpper(inputTransactionDN.DmoCategory)

	createTransactionErr := tx.Create(&createdTransaction).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, createTransactionErr
	}

	idNumber := createIdNumber("DN", createdTransaction.ID)

	updateTransactionsErr := tx.Model(&createdTransaction).Update("id_number", idNumber).Error

	if updateTransactionsErr != nil {
		tx.Rollback()
		return  createdTransaction, updateTransactionsErr
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

func (r *repository) DeleteTransactionDN(id int, userId uint) (bool, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ?", id).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := tx.Unscoped().Where("id = ?", id).Delete(&transaction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Minerba Dmo with id number %s and id %v", *transaction.IdNumber, transaction.ID)
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, createHistoryErr
}

func (r *repository) UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ?", idTransaction).First(&transaction).Error

	if errFind != nil {
		tx.Rollback()
		return transaction, errFind
	}

	beforeData , errorBeforeDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorBeforeDataJsonMarshal
	}

	inputEditTransactionDN.Seller = strings.ToUpper(inputEditTransactionDN.Seller)
	inputEditTransactionDN.TugboatName = strings.ToUpper(inputEditTransactionDN.TugboatName)
	inputEditTransactionDN.BargeName = strings.ToUpper(inputEditTransactionDN.BargeName)
	inputEditTransactionDN.VesselName = strings.ToUpper(inputEditTransactionDN.VesselName)
	inputEditTransactionDN.CustomerName = strings.ToUpper(inputEditTransactionDN.CustomerName)
	inputEditTransactionDN.LoadingPortName = strings.ToUpper(inputEditTransactionDN.LoadingPortName)
	inputEditTransactionDN.LoadingPortLocation = strings.ToUpper(inputEditTransactionDN.LoadingPortLocation)
	inputEditTransactionDN.UnloadingPortName = strings.ToUpper(inputEditTransactionDN.UnloadingPortName)
	inputEditTransactionDN.UnloadingPortLocation = strings.ToUpper(inputEditTransactionDN.UnloadingPortLocation)
	inputEditTransactionDN.DmoDestinationPort = strings.ToUpper(inputEditTransactionDN.DmoDestinationPort)
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
	inputEditTransactionDN.SurveyorName = strings.ToUpper(inputEditTransactionDN.SurveyorName)
	inputEditTransactionDN.CowNumber = strings.ToUpper(inputEditTransactionDN.CowNumber)
	inputEditTransactionDN.CoaNumber = strings.ToUpper(inputEditTransactionDN.CoaNumber)
	inputEditTransactionDN.SalesSystem = strings.ToUpper(inputEditTransactionDN.SalesSystem)
	inputEditTransactionDN.InvoiceNumber = strings.ToUpper(inputEditTransactionDN.InvoiceNumber)
	inputEditTransactionDN.ContractNumber = strings.ToUpper(inputEditTransactionDN.ContractNumber)
	inputEditTransactionDN.DmoBuyerName = strings.ToUpper(inputEditTransactionDN.DmoBuyerName)
	inputEditTransactionDN.DmoIndustryType = strings.ToUpper(inputEditTransactionDN.DmoIndustryType)
	inputEditTransactionDN.DmoCategory = strings.ToUpper(inputEditTransactionDN.DmoCategory)

	dataInput, errorMarshal := json.Marshal(inputEditTransactionDN)

	if errorMarshal != nil {
		tx.Rollback()
		return  transaction, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		tx.Rollback()
		return  transaction, errorUnmarshal
	}

	updateErr := tx.Model(&transaction).Updates(dataInputMapString).Error

	if updateErr != nil {
		tx.Rollback()
		return  transaction, updateErr
	}

	afterData , errorAfterDataJsonMarshal := json.Marshal(transaction)

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

func (r *repository) UploadDocumentTransactionDN (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error) {
	var uploadedTransaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ?", idTransaction).First(&uploadedTransaction).Error

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

// Minerba

func (r *repository) CreateMinerba (period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error) {
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

	for _, v := range transactions {
		createdMinerba.Quantity += v.Quantity
	}

	errCreateMinerba := tx.Create(&createdMinerba).Error

	if errCreateMinerba != nil {
		tx.Rollback()
		return createdMinerba, errCreateMinerba
	}

	idNumber := baseIdNumber + "-" + helper.CreateIdNumber(int(createdMinerba.ID))

	updateMinerbaErr := tx.Model(&createdMinerba).Update("id_number", idNumber).Error

	if updateMinerbaErr != nil {
		tx.Rollback()
		return  createdMinerba, updateMinerbaErr
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

func (r *repository) UpdateMinerba (id int, updateTransaction []int, userId uint) (minerba.Minerba, error) {
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
	beforeData , errorBeforeDataJsonMarshal := json.Marshal(historyBefore)
	afterData , errorAfterDataJsonMarshal := json.Marshal(historyAfter)

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

func (r *repository) DeleteMinerba (idMinerba int, userId uint) (bool, error) {

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
	return  minerba, nil
}

// Dmo

func (r *repository) CreateDmo (dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error) {
	var createdDmo dmo.Dmo
	var transactionBarge []transaction.Transaction
	var transactionVessel []transaction.Transaction

	barge := false
	vessel := false
	tx := r.db.Begin()

	createdDmo.Period = dmoInput.Period
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

		createdDmo.BargeTotalQuantity = bargeQuantity
		createdDmo.BargeGrandTotalQuantity = bargeQuantity
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

		createdDmo.VesselAdjustment = vesselAdjustment
		createdDmo.VesselTotalQuantity = vesselQuantity
		createdDmo.VesselGrandTotalQuantity = vesselQuantity + vesselAdjustment
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

	createdDmoErr := tx.Create(&createdDmo).Error

	if createdDmoErr != nil {
		tx.Rollback()
		return  createdDmo, createdDmoErr
	}

	idNumber := baseIdNumber + "-" + helper.CreateIdNumber(int(createdDmo.ID))

	updateDmoErr := tx.Model(&createdDmo).Update("id_number", idNumber).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return  createdDmo, updateDmoErr
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
				vesselDummy.VesselName = value.VesselName
				vesselDummy.Adjustment = value.Adjustment
				vesselDummy.Quantity = value.Quantity
				vesselDummy.DmoId = createdDmo.ID
				vesselDummy.GrandTotalQuantity = value.Quantity + value.Adjustment
				dmoVessels = append(dmoVessels, vesselDummy)
			}

			createDmoVesselsErr := tx.Create(&dmoVessels).Error

			if createDmoVesselsErr != nil {
				tx.Rollback()
				return createdDmo, createDmoVesselsErr
			}
		}
	}

	if len(transactionBarge) != len(dmoInput.TransactionBarge) && len(transactionVessel) != len(dmoInput.TransactionVessel){
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

func (r *repository) DeleteDmo (idDmo int, userId uint) (bool, error) {

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
			if strings.Contains(value["Location"].(string), "berita_acara") {
				editData["reconciliation_letter_document_link"] = value["Location"]
			}
			if strings.Contains(value["Location"].(string), "surat_pernyataan") {
				editData["statement_letter_document_link"] = value["Location"]
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
	return  dmoUpdate, nil
}

func(r *repository) UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData , errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

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

func(r *repository) UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint, location string) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData , errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return dmoUpdate, errorBeforeDataJsonMarshal
	}

	var field string
	updatesDmo := make(map[string]interface{})

	if isBast {
		field = "bast"
		updatesDmo["is_bast_document_signed"] = true
		updatesDmo["signed_bast_letter_document_link"] = location
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

func(r *repository) UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, id int, userId uint) (dmo.Dmo, error) {
	var dmoUpdate dmo.Dmo

	tx := r.db.Begin()

	findErr := tx.Where("id = ?", id).First(&dmoUpdate).Error

	if findErr != nil {
		tx.Rollback()
		return dmoUpdate, findErr
	}

	beforeData , errorBeforeDataJsonMarshal := json.Marshal(dmoUpdate)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return dmoUpdate, errorBeforeDataJsonMarshal
	}

	var field string
	updatesDmo := make(map[string]interface{})

	if isBast {
		field = "signed_bast false"
		updatesDmo["is_bast_document_signed"] = false
		updatesDmo["signed_bast_letter_document_link"] = nil
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
