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

		if sortFilter.Field == "vessel_id" {
			querySort = "vessels.name " + sortFilter.Sort
		}

		if querySort == "id desc" {
			querySort = fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
		}
	}

	if sortFilter.Quantity != 0 {
		quantity := fmt.Sprintf("%v", sortFilter.Quantity)
		queryFilter = queryFilter + "cast(grand_total_quantity AS TEXT) LIKE '%" + quantity + "%'"
	}

	if sortFilter.VesselId != "" {
		if queryFilter != "" {
			queryFilter += "AND vessel_id = " + sortFilter.VesselId
		} else {
			queryFilter = "vessel_id = " + sortFilter.VesselId
		}
	}

	if sortFilter.BlDateStart != "" {
		if queryFilter != "" {
			queryFilter += "AND bl_date >= '" + sortFilter.BlDateStart + "'"
		} else {
			queryFilter = "bl_date >= '" + sortFilter.BlDateStart + "'"
		}
	}

	if sortFilter.BlDateEnd != "" {
		if queryFilter != "" {
			queryFilter += "AND bl_date <= '" + sortFilter.BlDateEnd + "T23:59:59'"
		} else {
			queryFilter = "bl_date <= '" + sortFilter.BlDateEnd + "T23:59:59'"
		}
	}

	errFind := r.db.Preload(clause.Associations).Select("grouping_vessel_dns.*").Joins("LEFT JOIN vessels vessels on grouping_vessel_dns.vessel_id = vessels.id").Order(querySort).Where(queryFilter).Scopes(paginateData(listGroupingVesselDn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselDn).Error

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
