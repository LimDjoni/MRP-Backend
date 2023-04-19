package electricassignment

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListElectricAssignment(page int, sortFilter SortFilterElectricAssignment, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListElectricAssignment(page int, sortFilter SortFilterElectricAssignment, iupopkId int) (Pagination, error) {
	var listElectricAssignment []ElectricAssignment
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	querySort := "id desc"
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.Quantity != "" {
		queryFilter = queryFilter + " AND cast(grand_total_quantity AS TEXT) LIKE '%" + sortFilter.Quantity + "%'"
	}

	if sortFilter.Year != "" {
		queryFilter += " AND year = '" + sortFilter.Year + "'"
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listElectricAssignment, &pagination, r.db, queryFilter)).Find(&listElectricAssignment).Error

	if errFind != nil {

		return pagination, errFind

	}

	pagination.Data = listElectricAssignment

	return pagination, nil
}
