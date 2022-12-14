package dmo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetReportDmoWithPeriod(period string) (Dmo, error)
	GetListReportDmoAll(page int, filterDmo FilterAndSortDmo) (Pagination, error)
	GetDataDmo(id int) (Dmo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportDmoWithPeriod(period string) (Dmo, error) {
	var reportDmo Dmo

	errFind := r.db.Where("period = ?", period).First(&reportDmo).Error

	return reportDmo, errFind
}

func (r *repository) GetListReportDmoAll(page int, filterDmo FilterAndSortDmo) (Pagination, error) {
	var listReportDmo []Dmo

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := ""
	sortFilter := "id desc"

	if filterDmo.Field != "" && filterDmo.Sort != "" {
		sortFilter = filterDmo.Field + " " + filterDmo.Sort

		if strings.ToLower(filterDmo.Field) == "period" {
			sortFilter = "to_date(period,'Mon Year') " + filterDmo.Sort
		}

		if strings.ToLower(filterDmo.Field) == "quantity" {
			sortFilter = "vessel_total_quantity + barge_total_quantity " + filterDmo.Sort
		}

		if strings.ToLower(filterDmo.Field) == "grand_total_quantity" {
			sortFilter = "vessel_grand_total_quantity + barge_grand_total_quantity " + filterDmo.Sort
		}
	}

	if filterDmo.CreatedStart != "" {
		queryFilter = queryFilter + "created_at >= '" + filterDmo.CreatedStart + "'"
	}

	if filterDmo.CreatedEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND created_at <= '" + filterDmo.CreatedEnd + "T23:59:59'"
		} else {
			queryFilter = "created_at <= '" + filterDmo.CreatedEnd + "T23:59:59'"
		}
	}

	if filterDmo.Quantity != 0 {
		quantity := fmt.Sprintf("%v", filterDmo.Quantity)
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(vessel_grand_total_quantity + barge_grand_total_quantity AS TEXT) LIKE '%" + quantity + "%'"
		} else {
			queryFilter = "cast(vessel_grand_total_quantity + barge_grand_total_quantity AS TEXT) LIKE '%" + quantity + "%'"
		}
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateDmo(listReportDmo, &pagination, r.db, queryFilter)).Find(&listReportDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportDmo

	return pagination, nil
}

func (r *repository) GetDataDmo(id int) (Dmo, error) {
	var dmo Dmo

	errFind := r.db.Where("id = ?", id).First(&dmo).Error

	return dmo, errFind
}
