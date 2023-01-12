package groupingvesselln

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error) {
	var listGroupingVesselLn []GroupingVesselLn
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

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselLn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselLn).Error

	if errFind != nil {
		errWithoutOrder := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselLn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselLn).Error

		if errWithoutOrder != nil {
			pagination.Data = listGroupingVesselLn
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = listGroupingVesselLn

	return pagination, nil
}
