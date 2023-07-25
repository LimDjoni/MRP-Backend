package royaltyreport

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetTransactionRoyaltyReport(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReportData, error)
	GetDetailTransactionRoyaltyReport(id int, iupopkId int) (RoyaltyReportDetail, error)
	ListRoyaltyReport(page int, sortFilter SortFilterRoyaltyReport, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionRoyaltyReport(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReportData, error) {
	var listTransaction []RoyaltyReportData

	errFind := r.db.Table("transactions").Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ?", dateFrom, dateTo, iupopkId).Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return listTransaction, errFind
	}

	return listTransaction, nil
}

func (r *repository) GetDetailTransactionRoyaltyReport(id int, iupopkId int) (RoyaltyReportDetail, error) {
	var detailRoyaltyReport RoyaltyReportDetail
	var royaltyReport RoyaltyReport
	var listTransaction []RoyaltyReportData

	errFindRoyaltyReport := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&royaltyReport).Error

	if errFindRoyaltyReport != nil {
		return detailRoyaltyReport, errFindRoyaltyReport
	}

	detailRoyaltyReport.Detail = royaltyReport

	errFind := r.db.Preload(clause.Associations).Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ?", royaltyReport.DateFrom, royaltyReport.DateTo, iupopkId).Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return detailRoyaltyReport, errFind
	}

	detailRoyaltyReport.ListTransaction = listTransaction

	return detailRoyaltyReport, nil
}

func (r *repository) ListRoyaltyReport(page int, sortFilter SortFilterRoyaltyReport, iupopkId int) (Pagination, error) {
	var royaltyReport []RoyaltyReport
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

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(royaltyReport, &pagination, r.db, queryFilter)).Find(&royaltyReport).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = royaltyReport

	return pagination, nil
}
