package coareport

import (
	"ajebackend/model/transaction"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
	GetDetailTransactionCoaReport(id int, iupopkId int) (CoaReportDetail, error)
	ListCoaReport(page int, sortFilter SortFilterCoaReport, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	var listTransaction []transaction.Transaction

	errFind := r.db.Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ? AND transaction_type = ?", dateFrom, dateTo, iupopkId, "DN").Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return listTransaction, errFind
	}

	return listTransaction, nil
}

func (r *repository) GetDetailTransactionCoaReport(id int, iupopkId int) (CoaReportDetail, error) {
	var detailCoaReport CoaReportDetail
	var coaReport CoaReport
	var listTransaction []transaction.Transaction

	errFindCoaReport := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&coaReport).Error

	if errFindCoaReport != nil {
		return detailCoaReport, errFindCoaReport
	}

	detailCoaReport.Detail = coaReport

	errFind := r.db.Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ? AND transaction_type = ?", coaReport.DateFrom, coaReport.DateTo, iupopkId, "DN").Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return detailCoaReport, errFind
	}

	detailCoaReport.ListTransaction = listTransaction

	return detailCoaReport, nil
}

func (r *repository) ListCoaReport(page int, sortFilter SortFilterCoaReport, iupopkId int) (Pagination, error) {
	var coaReports []CoaReport
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

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(coaReports, &pagination, r.db, queryFilter)).Find(&coaReports).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = coaReports

	return pagination, nil
}
