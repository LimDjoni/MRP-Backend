package production

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	GetListProduction(page int, filter FilterListProduction) (Pagination, error)
	DetailProduction(id int) (Production, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func(r *repository) GetListProduction(page int, filter FilterListProduction) (Pagination, error) {
	var listProduction []Production

	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	queryFilter := ""

	if filter.ProductionDateStart != "" {
		queryFilter = queryFilter + "production_date >= '" + filter.ProductionDateStart + "'"
	}

	if filter.ProductionDateEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND production_date <= '" + filter.ProductionDateEnd + "T23:59:59'"
		} else {
			queryFilter = "production_date <= '" + filter.ProductionDateEnd + "T23:59:59'"
		}
	}

	if filter.Quantity != 0 {
		quantity := fmt.Sprintf("%v", filter.Quantity)
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" +  quantity + "%'"
		} else {
			queryFilter = "cast(quantity AS TEXT) LIKE '%" +  quantity + "%'"
		}
	}

	errFind := r.db.Where(queryFilter).Scopes(paginateProduction(listProduction, &pagination, r.db, queryFilter)).Find(&listProduction).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listProduction

	return pagination, nil
}

func (r *repository) DetailProduction(id int) (Production, error) {
	var detailProduction Production

	errFind := r.db.Where("id = ?", id).First(&detailProduction).Error

	return detailProduction, errFind
}
