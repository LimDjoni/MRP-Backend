package transaction

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error) {
	var transactions []Transaction
	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	defaultSort := "id desc"
	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	queryFilter := fmt.Sprintf("transaction_type = '%s' ", "DN")

	if sortFilter.ShipName != "" {
		queryFilter = queryFilter + " AND ship_name ILIKE '%" + sortFilter.ShipName + "%'"
	}

	if sortFilter.BargeName != "" {
		queryFilter = queryFilter + " AND barge_name ILIKE '%" + sortFilter.BargeName + "%'"
	}

	if sortFilter.ShippingFrom != "" {
		queryFilter = queryFilter + " AND shipping_date >= '" + sortFilter.ShippingFrom + "'"
	}

	if sortFilter.ShippingTo != "" {
		queryFilter = queryFilter + " AND shipping_date <= '" + sortFilter.ShippingTo + "'"
	}

	errFind := r.db.Where(queryFilter).Order(sortString).Scopes(paginateDataDN(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

	if errFind != nil {
		errWithoutOrder := r.db.Where(queryFilter).Order(defaultSort).Scopes(paginateDataDN(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

		if errWithoutOrder != nil {
			pagination.Data = transactions
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = transactions

	return pagination, nil
}

func (r *repository) DetailTransactionDN(id int) (Transaction, error) {
	var transaction Transaction

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	return transaction, errFind
}
