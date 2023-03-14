package production

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error)
	DetailProduction(id int, iupopkId int) (Production, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error) {
	var listProduction []Production

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
	querySort := "id desc"

	if filter.Field != "" && filter.Sort != "" {
		querySort = filter.Field + " " + filter.Sort
	}

	if filter.ProductionDateStart != "" {
		queryFilter = queryFilter + " AND production_date >= '" + filter.ProductionDateStart + "'"
	}

	if filter.ProductionDateEnd != "" {
		queryFilter = queryFilter + " AND production_date <= '" + filter.ProductionDateEnd + "T23:59:59'"
	}

	if filter.Quantity != "" {
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + filter.Quantity + "%'"
	}

	errFind := r.db.Where(queryFilter).Order(querySort).Scopes(paginateProduction(listProduction, &pagination, r.db, queryFilter)).Find(&listProduction).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listProduction

	return pagination, nil
}

func (r *repository) DetailProduction(id int, iupopkId int) (Production, error) {
	var detailProduction Production

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&detailProduction).Error

	return detailProduction, errFind
}
