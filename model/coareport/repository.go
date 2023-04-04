package coareport

import (
	"ajebackend/model/transaction"

	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	var listTransaction []transaction.Transaction

	errFind := r.db.Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ?", dateFrom, dateTo, iupopkId).Find(&listTransaction).Error

	if errFind != nil {
		return listTransaction, errFind
	}

	return listTransaction, nil
}
