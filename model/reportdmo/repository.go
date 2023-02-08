package reportdmo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetReportDmoWithPeriod(period string) (ReportDmo, error)
	GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportDmoWithPeriod(period string) (ReportDmo, error) {
	var reportDmo ReportDmo

	errFind := r.db.Where("period = ?", period).First(&reportDmo).Error

	return reportDmo, errFind
}

func (r *repository) GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo) (Pagination, error) {
	var listReportDmo []ReportDmo

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := ""
	sortFilter := "id desc"

	if filterReportDmo.Field != "" && filterReportDmo.Sort != "" {
		sortFilter = filterReportDmo.Field + " " + filterReportDmo.Sort

		if strings.ToLower(filterReportDmo.Field) == "period" {
			sortFilter = "to_date(period,'Mon Year') " + filterReportDmo.Sort
		}
	}

	if filterReportDmo.CreatedStart != "" {
		queryFilter = queryFilter + "created_at >= '" + filterReportDmo.CreatedStart + "'"
	}

	if filterReportDmo.CreatedEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND created_at <= '" + filterReportDmo.CreatedEnd + "T23:59:59'"
		} else {
			queryFilter = "created_at <= '" + filterReportDmo.CreatedEnd + "T23:59:59'"
		}
	}

	if filterReportDmo.Quantity != 0 {
		quantity := fmt.Sprintf("%v", filterReportDmo.Quantity)
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
		} else {
			queryFilter = "cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
		}
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateReportDmo(listReportDmo, &pagination, r.db, queryFilter)).Find(&listReportDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportDmo

	return pagination, nil
}
