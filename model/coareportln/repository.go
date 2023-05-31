package coareportln

import (
	"ajebackend/model/transaction"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetTransactionCoaReportLn(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
	GetDetailTransactionCoaReportLn(id int, iupopkId int) (CoaReportLnDetail, error)
	ListCoaReportLn(page int, sortFilter SortFilterCoaReportLn, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionCoaReportLn(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	var listTransaction []transaction.Transaction

	errFind := r.db.Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ? AND transaction_type = ?", dateFrom, dateTo, iupopkId, "LN").Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return listTransaction, errFind
	}

	return listTransaction, nil
}

func (r *repository) GetDetailTransactionCoaReportLn(id int, iupopkId int) (CoaReportLnDetail, error) {
	var detailCoaReportLn CoaReportLnDetail
	var coaReportLn CoaReportLn
	var listTransaction []transaction.Transaction

	errFindCoaReportLn := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&coaReportLn).Error

	if errFindCoaReportLn != nil {
		return detailCoaReportLn, errFindCoaReportLn
	}

	detailCoaReportLn.Detail = coaReportLn

	errFind := r.db.Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ? AND transaction_type = ?", coaReportLn.DateFrom, coaReportLn.DateTo, iupopkId, "LN").Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return detailCoaReportLn, errFind
	}

	detailCoaReportLn.ListTransaction = listTransaction

	return detailCoaReportLn, nil
}

func (r *repository) ListCoaReportLn(page int, sortFilter SortFilterCoaReportLn, iupopkId int) (Pagination, error) {
	var coaReportLns []CoaReportLn
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "id desc"

	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)

	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.DateStart != "" {
		queryFilter += " AND date_from = '" + sortFilter.DateStart + "'"
	}

	if sortFilter.DateEnd != "" {
		queryFilter += " AND date_to = '" + sortFilter.DateEnd + "'"
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(coaReportLns, &pagination, r.db, queryFilter)).Find(&coaReportLns).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = coaReportLns

	return pagination, nil
}
