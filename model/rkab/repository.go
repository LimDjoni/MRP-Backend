package rkab

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error) {
	var listRkab []Rkab
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "id desc"

	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)

	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.DateFrom != "" {
		queryFilter += " AND date_of_issue >= '" + sortFilter.DateFrom + "'"
	}

	if sortFilter.DateTo != "" {
		queryFilter += " AND date_of_issue <= '" + sortFilter.DateTo + "'"
	}

	if sortFilter.Year != "" {
		queryFilter += " AND year = " + sortFilter.Year
	}

	if sortFilter.ProductionQuota != "" {
		queryFilter = queryFilter + " AND cast(production_quota AS TEXT) LIKE '%" + sortFilter.ProductionQuota + "%'"
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(listRkab, &pagination, r.db, queryFilter)).Find(&listRkab).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listRkab

	return pagination, nil
}
