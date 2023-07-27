package royaltyrecon

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetTransactionRoyaltyRecon(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReconData, error)
	GetDetailTransactionRoyaltyRecon(id int, iupopkId int) (RoyaltyReconDetail, error)
	ListRoyaltyRecon(page int, sortFilter SortFilterRoyaltyRecon, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionRoyaltyRecon(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReconData, error) {
	var listTransaction []RoyaltyReconData

	errFind := r.db.Table("transactions").Preload(clause.Associations).Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ?", dateFrom, dateTo, iupopkId).Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return listTransaction, errFind
	}

	return listTransaction, nil
}

func (r *repository) GetDetailTransactionRoyaltyRecon(id int, iupopkId int) (RoyaltyReconDetail, error) {
	var detailRoyaltyRecon RoyaltyReconDetail
	var royaltyRecon RoyaltyRecon
	var listTransaction []RoyaltyReconData

	errFindRoyaltyRecon := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&royaltyRecon).Error

	if errFindRoyaltyRecon != nil {
		return detailRoyaltyRecon, errFindRoyaltyRecon
	}

	detailRoyaltyRecon.Detail = royaltyRecon

	errFind := r.db.Table("transactions").Preload(clause.Associations).Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("shipping_date >= ? AND shipping_date <= ? AND seller_id = ?", royaltyRecon.DateFrom, royaltyRecon.DateTo, iupopkId).Order("shipping_date desc").Find(&listTransaction).Error

	if errFind != nil {
		return detailRoyaltyRecon, errFind
	}

	detailRoyaltyRecon.ListTransaction = listTransaction

	return detailRoyaltyRecon, nil
}

func (r *repository) ListRoyaltyRecon(page int, sortFilter SortFilterRoyaltyRecon, iupopkId int) (Pagination, error) {
	var royaltyRecon []RoyaltyRecon
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

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(royaltyRecon, &pagination, r.db, queryFilter)).Find(&royaltyRecon).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = royaltyRecon

	return pagination, nil
}
