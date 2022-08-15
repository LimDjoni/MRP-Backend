package transaction

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository interface {
	CreateTransactionDN (inputTransactionDN DataTransactionInput) (Transaction, error)
	ListDataDN (page int) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	DeleteTransaction(id int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateTransactionDN (inputTransactionDN DataTransactionInput) (Transaction, error) {
	var createdTransaction Transaction
	var transactions []Transaction

	year, month, _ := time.Now().Date()
	startDate := fmt.Sprintf("%v-%v-01  00:00:00", year, int(month))
	endDate := fmt.Sprintf("%v-%v-31  00:00:00", year, int(month))

	findErr := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error

	if findErr != nil {
		return createdTransaction, findErr
	}

	createdTransaction.DmoId = nil
	createdTransaction.Number = len(transactions) + 1
	createdTransaction.IdNumber += fmt.Sprintf("DN-%v-%v-%v", year, int(month), strconv.Itoa(len(transactions) + 1))
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

	createErr := r.db.Create(&createdTransaction).Error

	return createdTransaction, createErr
}

func (r *repository) ListDataDN (page int) (Pagination, error) {
	var transactions []Transaction
	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	errFind := r.db.Where("transaction_type = ?", "DN").Scopes(paginateDataDN(transactions, &pagination, r.db)).Find(&transactions).Error
	pagination.Data = transactions

	return pagination, errFind
}

func (r *repository) DetailTransactionDN(id int) (Transaction, error) {
	var transaction Transaction

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	return transaction, errFind
}

func (r *repository) DeleteTransaction(id int) (bool, error) {
	var transaction Transaction

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Where("id = ?", id).Delete(&transaction).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, errDelete
}
