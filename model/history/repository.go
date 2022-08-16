package history

import (
	"ajebackend/model/transaction"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository interface {
	CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var createdTransaction transaction.Transaction
	var totalCount int64
	year, month, _ := time.Now().Date()
	startDate := fmt.Sprintf("%v-%v-01  00:00:00", year, int(month))
	endDate := fmt.Sprintf("%v-%v-31  00:00:00", year, int(month))

	tx := r.db.Begin()

	findErr := tx.Unscoped().Where("created_at BETWEEN ? AND ?", startDate, endDate).Model(transaction.Transaction{}).Count(&totalCount).Error

	if findErr != nil {
		return createdTransaction, findErr
	}

	createdTransaction.DmoId = nil
	createdTransaction.Number = int(totalCount) + 1
	createdTransaction.IdNumber += fmt.Sprintf("DN-%v-%v-%v", year, int(month), strconv.Itoa(int(totalCount) + 1))
	createdTransaction.TransactionType = "DN"
	if inputTransactionDN.ShippingDate != "" {
		createdTransaction.ShippingDate = &inputTransactionDN.ShippingDate
	} else {
		createdTransaction.ShippingDate = nil
	}
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

	if inputTransactionDN.SkbDate != "" {
		createdTransaction.SkbDate = &inputTransactionDN.SkbDate
	} else {
		createdTransaction.SkbDate = nil
	}
	createdTransaction.SkbNumber = inputTransactionDN.SkbNumber

	if inputTransactionDN.SkabDate != "" {
		createdTransaction.SkabDate = &inputTransactionDN.SkabDate
	} else {
		createdTransaction.SkabDate = nil
	}
	createdTransaction.SkabNumber = inputTransactionDN.SkabNumber
	if inputTransactionDN.BillOfLadingDate != "" {
		createdTransaction.BillOfLadingDate = &inputTransactionDN.BillOfLadingDate
	} else {
		createdTransaction.BillOfLadingDate = nil
	}
	createdTransaction.BillOfLadingNumber = inputTransactionDN.BillOfLadingNumber
	createdTransaction.RoyaltyRate = inputTransactionDN.RoyaltyRate
	createdTransaction.DpRoyaltyPrice = inputTransactionDN.DpRoyaltyPrice

	if inputTransactionDN.DpRoyaltyDate != "" {
		createdTransaction.DpRoyaltyDate = &inputTransactionDN.DpRoyaltyDate
	} else {
		createdTransaction.DpRoyaltyDate = nil
	}
	createdTransaction.DpRoyaltyNtpn = inputTransactionDN.DpRoyaltyNtpn
	createdTransaction.DpRoyaltyBillingCode = inputTransactionDN.DpRoyaltyBillingCode
	createdTransaction.DpRoyaltyTotal = inputTransactionDN.DpRoyaltyTotal
	createdTransaction.PaymentDpRoyaltyPrice = inputTransactionDN.PaymentDpRoyaltyPrice
	if inputTransactionDN.PaymentDpRoyaltyDate != "" {
		createdTransaction.PaymentDpRoyaltyDate = &inputTransactionDN.PaymentDpRoyaltyDate
	} else {
		createdTransaction.PaymentDpRoyaltyDate = nil
	}
	createdTransaction.PaymentDpRoyaltyNtpn = inputTransactionDN.PaymentDpRoyaltyNtpn
	createdTransaction.PaymentDpRoyaltyBillingCode = inputTransactionDN.PaymentDpRoyaltyBillingCode
	createdTransaction.PaymentDpRoyaltyTotal = inputTransactionDN.PaymentDpRoyaltyTotal
	if inputTransactionDN.LhvDate != "" {
		createdTransaction.LhvDate = &inputTransactionDN.LhvDate
	} else {
		createdTransaction.LhvDate = nil
	}
	createdTransaction.LhvNumber = inputTransactionDN.LhvNumber
	createdTransaction.SurveyorName = inputTransactionDN.SurveyorName
	if inputTransactionDN.CowDate != "" {
		createdTransaction.CowDate = &inputTransactionDN.CowDate
	} else {
		createdTransaction.CowDate = nil
	}
	createdTransaction.CowNumber = inputTransactionDN.CowNumber
	if inputTransactionDN.CoaDate != "" {
		createdTransaction.CoaDate = &inputTransactionDN.CoaDate
	} else {
		createdTransaction.CoaDate = nil
	}
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
	if inputTransactionDN.InvoiceDate != "" {
		createdTransaction.InvoiceDate = &inputTransactionDN.InvoiceDate
	} else {
		createdTransaction.InvoiceDate = nil
	}
	createdTransaction.InvoiceNumber = inputTransactionDN.InvoiceNumber
	createdTransaction.InvoicePriceUnit = inputTransactionDN.InvoicePriceUnit
	createdTransaction.InvoicePriceTotal = inputTransactionDN.InvoicePriceTotal
	createdTransaction.DmoReconciliationLetter = inputTransactionDN.DmoReconciliationLetter
	if inputTransactionDN.ContractDate != "" {
		createdTransaction.ContractDate = &inputTransactionDN.ContractDate
	} else {
		createdTransaction.ContractDate = nil
	}
	createdTransaction.ContractNumber = inputTransactionDN.ContractNumber
	createdTransaction.DmoBuyerName = inputTransactionDN.DmoBuyerName
	createdTransaction.DmoIndustryType = inputTransactionDN.DmoIndustryType
	createdTransaction.DmoStatusReconciliationLetter = inputTransactionDN.DmoStatusReconciliationLetter

	createTransactionErr := tx.Create(&createdTransaction).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, createTransactionErr
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

func (r *repository) DeleteTransaction(id int, userId uint) (bool, error) {
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
