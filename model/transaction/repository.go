package transaction

import (
	"gorm.io/gorm"
)

type Repository interface {
	ListDataDN (page int) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
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
