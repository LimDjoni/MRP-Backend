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
	DeleteMinerba (idMinerba int, userId uint) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error)
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
		createdTransaction.Seller = inputTransactionDN.Seller
	}
	createdTransaction.ShippingDate = inputTransactionDN.ShippingDate
	createdTransaction.Quantity = inputTransactionDN.Quantity
	createdTransaction.ShipName = inputTransactionDN.ShipName
	createdTransaction.BargeName = inputTransactionDN.BargeName
	createdTransaction.VesselName = inputTransactionDN.VesselName
	createdTransaction.CustomerName = inputTransactionDN.CustomerName
	createdTransaction.LoadingPortName = inputTransactionDN.LoadingPortName
	createdTransaction.LoadingPortLocation = inputTransactionDN.LoadingPortLocation
	createdTransaction.UnloadingPortName = inputTransactionDN.UnloadingPortName
	createdTransaction.UnloadingPortLocation = inputTransactionDN.UnloadingPortLocation
	createdTransaction.DmoDestinationPort = inputTransactionDN.DmoDestinationPort
	createdTransaction.SkbDate = inputTransactionDN.SkbDate
	createdTransaction.SkbNumber = inputTransactionDN.SkbNumber
	createdTransaction.SkabDate = inputTransactionDN.SkabDate
	createdTransaction.SkabNumber = inputTransactionDN.SkabNumber
	createdTransaction.BillOfLadingDate = inputTransactionDN.BillOfLadingDate
	createdTransaction.BillOfLadingNumber = inputTransactionDN.BillOfLadingNumber
	createdTransaction.RoyaltyRate = inputTransactionDN.RoyaltyRate
	createdTransaction.DpRoyaltyPrice = inputTransactionDN.DpRoyaltyPrice
	createdTransaction.DpRoyaltyCurrency = inputTransactionDN.DpRoyaltyCurrency
	if inputTransactionDN.DpRoyaltyCurrency == "" {
		createdTransaction.DpRoyaltyCurrency = "IDR"
	}
	createdTransaction.DpRoyaltyDate = inputTransactionDN.DpRoyaltyDate
	createdTransaction.DpRoyaltyNtpn = inputTransactionDN.DpRoyaltyNtpn
	createdTransaction.DpRoyaltyBillingCode = inputTransactionDN.DpRoyaltyBillingCode
	createdTransaction.DpRoyaltyTotal = inputTransactionDN.DpRoyaltyTotal
	createdTransaction.PaymentDpRoyaltyPrice = inputTransactionDN.PaymentDpRoyaltyPrice
	createdTransaction.PaymentDpRoyaltyCurrency = inputTransactionDN.PaymentDpRoyaltyCurrency
	if inputTransactionDN.PaymentDpRoyaltyCurrency == "" {
		createdTransaction.PaymentDpRoyaltyCurrency = "IDR"
	}
	createdTransaction.PaymentDpRoyaltyDate = inputTransactionDN.PaymentDpRoyaltyDate
	createdTransaction.PaymentDpRoyaltyNtpn = inputTransactionDN.PaymentDpRoyaltyNtpn
	createdTransaction.PaymentDpRoyaltyBillingCode = inputTransactionDN.PaymentDpRoyaltyBillingCode
	createdTransaction.PaymentDpRoyaltyTotal = inputTransactionDN.PaymentDpRoyaltyTotal
	createdTransaction.LhvDate = inputTransactionDN.LhvDate
	createdTransaction.LhvNumber = inputTransactionDN.LhvNumber
	createdTransaction.SurveyorName = inputTransactionDN.SurveyorName
	createdTransaction.CowDate = inputTransactionDN.CowDate
	createdTransaction.CowNumber = inputTransactionDN.CowNumber
	createdTransaction.CoaDate = inputTransactionDN.CoaDate
	createdTransaction.CoaNumber = inputTransactionDN.CoaNumber
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
	createdTransaction.InvoiceNumber = inputTransactionDN.InvoiceNumber
	createdTransaction.InvoicePriceUnit = inputTransactionDN.InvoicePriceUnit
	createdTransaction.InvoicePriceTotal = inputTransactionDN.InvoicePriceTotal
	createdTransaction.DmoReconciliationLetter = inputTransactionDN.DmoReconciliationLetter
	createdTransaction.ContractDate = inputTransactionDN.ContractDate
	createdTransaction.ContractNumber = inputTransactionDN.ContractNumber
	createdTransaction.DmoBuyerName = inputTransactionDN.DmoBuyerName
	createdTransaction.DmoIndustryType = inputTransactionDN.DmoIndustryType
	createdTransaction.DmoCategory = inputTransactionDN.DmoCategory

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

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Where("id = ?", id).Delete(&transaction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	uId := uint(id)
	history.TransactionId = &uId
	history.UserId = userId
	history.Status = "Deleted"

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

func (r *repository) DeleteMinerba (idMinerba int, userId uint) (bool, error) {

	tx := r.db.Begin()
	var minerba minerba.Minerba

	findErr := tx.Where("id = ?", idMinerba).First(&minerba).Error

	if findErr != nil {
		tx.Rollback()
		return false, findErr
	}

	updateTransactionErr := tx.Table("transactions").Where("minerba_id = ?", idMinerba).Update("minerba_id", nil).Error

	if updateTransactionErr != nil {
		tx.Rollback()
		return false, updateTransactionErr
	}

	errDelete := tx.Unscoped().Where("id = ?", idMinerba).Delete(&minerba).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	history.Status = fmt.Sprintf("Deleted Minerba Report with id number %s and id %v", minerba.IdNumber, minerba.ID)
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

func (r *repository) CreateDmo (dmoInput dmo.CreateDmoInput, userId uint) (dmo.Dmo, error) {
	var createdDmo dmo.Dmo
	var transactionBarge []transaction.Transaction
	var transactionVessel []transaction.Transaction

	barge := false
	vessel := false
	tx := r.db.Begin()

	if len(dmoInput.TransactionBarge) > 0 {
		var bargeQuantity float64
		barge = true
		findTransactionBargeErr := tx.Where("id IN ?", dmoInput.TransactionBarge).Find(&transactionBarge).Error

		if findTransactionBargeErr != nil {
			return createdDmo, findTransactionBargeErr
		}

		for _, v := range transactionBarge {
			bargeQuantity += bargeQuantity + v.Quantity
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

		for _, v := range transactionBarge {
			vesselQuantity += vesselQuantity + v.Quantity
		}

		createdDmo.VesselTotalQuantity = vesselQuantity

		for _, v := range dmoInput.VesselAdjustment {
			vesselAdjustment += vesselAdjustment + v.Adjustment
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

	createdDmo.EndUser = dmoInput.EndUser

	createdDmoErr := tx.Create(&createdDmo).Error

	if createdDmoErr != nil {
		tx.Rollback()
		return  createdDmo, createdDmoErr
	}

	idNumber := createIdNumber("DD", createdDmo.ID)

	updateDmoErr := tx.Model(&createdDmo).Update("id_number", idNumber).Error

	if updateDmoErr != nil {
		tx.Rollback()
		return  createdDmo, updateDmoErr
	}

	findTransactionsBargeErr := tx.Where("id IN ?", dmoInput.TransactionBarge).Find(&transactionBarge).Error

	if findTransactionsBargeErr != nil {
		tx.Rollback()
		return createdDmo, findTransactionsBargeErr
	}

	findTransactionsVesselErr := tx.Where("id IN ?", dmoInput.TransactionVessel).Find(&transactionVessel).Error

	if findTransactionsVesselErr != nil {
		tx.Rollback()
		return createdDmo, findTransactionsVesselErr
	}

	if len(transactionBarge) != len(dmoInput.TransactionBarge) && len(transactionVessel) != len(dmoInput.TransactionVessel){
		tx.Rollback()
		return createdDmo, errors.New("please check some of transactions not found")
	}

	updateTransactionBargeErr := tx.Table("transactions").Where("id IN ?", dmoInput.TransactionBarge).Update("dmo_id", createdDmo.ID).Error

	if updateTransactionBargeErr != nil {
		tx.Rollback()
		return createdDmo, updateTransactionBargeErr
	}

	updateTransactionVesselErr := tx.Table("transactions").Where("id IN ?", dmoInput.TransactionVessel).Update("dmo_id", createdDmo.ID).Error

	if updateTransactionVesselErr != nil {
		tx.Rollback()
		return createdDmo, updateTransactionVesselErr
	}

	var dmoVessels []dmovessel.DmoVessel

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

	var traderDmo []traderdmo.TraderDmo

	for idx, value := range dmoInput.Trader {
		var traderDummy traderdmo.TraderDmo

		traderDummy.DmoId = createdDmo.ID
		traderDummy.TraderId = value.ID
		traderDummy.Order = idx + 1

		traderDmo = append(traderDmo, traderDummy)
	}

	createTraderDmoErr := tx.Create(&dmoVessels).Error

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
