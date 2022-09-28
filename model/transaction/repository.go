package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDNAndMinerba(listData []int)(bool, error)
	GetDetailMinerba(id int)(DetailMinerba, error)
	ListDataDNWithoutDmo() ([]Transaction, error)
	CheckDataDNAndDmo(listData []int)(bool, error)
	GetDetailDmo(id int)(DetailDmo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Transaction

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

	if sortFilter.Quantity != 0 {
		queryFilter = fmt.Sprintf("%s AND quantity = %v", queryFilter, sortFilter.Quantity)
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

func (r *repository) ListDataDNWithoutMinerba() ([]Transaction, error) {
	var listDataDnWithoutMinerba []Transaction

	errFind := r.db.Where("minerba_id is NULL AND transaction_type = ?", "DN").Find(&listDataDnWithoutMinerba).Error

	return listDataDnWithoutMinerba, errFind
}

// Minerba

func (r *repository) CheckDataDNAndMinerba(listData []int)(bool, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ?", listData).Find(&listDnValid).Error

	if errFindValid != nil {
		return false, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return false, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("minerba_id = ? AND id IN ?", nil, listData).Find(&listDn).Error

	if errFind != nil {
		return false, errFind
	}

	if len(listDn) == 0 {
		return false, errors.New("please check there is transaction already in report")
	}

	return true, nil
}

func(r *repository) GetDetailMinerba(id int)(DetailMinerba, error) {

	var detailMinerba DetailMinerba

	var minerba minerba.Minerba
	var transactions []Transaction

	minerbaFindErr := r.db.Where("id = ?", id).First(&minerba).Error

	if minerbaFindErr != nil {
		return detailMinerba, minerbaFindErr
	}

	detailMinerba.Detail = minerba

	transactionFindErr := r.db.Where("minerba_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailMinerba, transactionFindErr
	}

	detailMinerba.List = transactions
	return detailMinerba, nil
}

// DMO

func (r *repository) ListDataDNWithoutDmo() ([]Transaction, error) {
	var listDataDnWithoutDmo []Transaction

	errFind := r.db.Where("dmo_id is NULL AND transaction_type = ?", "DN").Find(&listDataDnWithoutDmo).Error

	return listDataDnWithoutDmo, errFind
}

func (r *repository) CheckDataDNAndDmo(listData []int)(bool, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ?", listData).Find(&listDnValid).Error

	if errFindValid != nil {
		return false, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return false, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("dmo_id = ? AND id IN ?", nil, listData).Find(&listDn).Error

	if errFind != nil {
		return false, errFind
	}

	if len(listDn) == 0 {
		return false, errors.New("please check there is transaction already in report")
	}

	return true, nil
}

func(r *repository) GetDetailDmo(id int)(DetailDmo, error) {

	var detailDmo DetailDmo

	var dmoData dmo.Dmo
	var transactions []Transaction

	dmoFindErr := r.db.Where("id = ?", id).First(&dmoData).Error

	if dmoFindErr != nil {
		return detailDmo, dmoFindErr
	}

	detailDmo.Detail = dmoData

	transactionFindErr := r.db.Where("dmo_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailDmo, transactionFindErr
	}

	detailDmo.List = transactions
	return detailDmo, nil
}
