package groupingvesseldn

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn) (Pagination, error) {
	var listGroupingVesselDn []GroupingVesselDn
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	querySort := "id desc"
	queryFilter := ""

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	}

	if sortFilter.Quantity != 0 {
		quantity := fmt.Sprintf("%v", sortFilter.Quantity)
		queryFilter = queryFilter + "cast(grand_total_quantity AS TEXT) LIKE '%" + quantity + "%'"
	}

	if sortFilter.VesselName != "" {
		if queryFilter != "" {
			queryFilter += "AND vessel_name LIKE '%" + sortFilter.VesselName + "%'"
		} else {
			queryFilter = "vessel_name LIKE '%" + sortFilter.VesselName + "%'"
		}
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselDn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselDn).Error

	if errFind != nil {
		errWithoutOrder := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselDn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselDn).Error

		if errWithoutOrder != nil {
			pagination.Data = listGroupingVesselDn
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = listGroupingVesselDn

	return pagination, nil
}
