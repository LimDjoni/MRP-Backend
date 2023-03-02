package reportdmo

import (
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

		if strings.ToLower(filterReportDmo.Field) == "quantity" {
			sortFilter = "quantity " + filterReportDmo.Sort
		}

		if strings.ToLower(filterReportDmo.Field) == "updatedat" {
			sortFilter = "updated_at " + filterReportDmo.Sort
		}
	}

	if filterReportDmo.UpdatedStart != "" {
		queryFilter = queryFilter + "updated_at >= '" + filterReportDmo.UpdatedStart + "'"
	}

	if filterReportDmo.UpdatedEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND updated_at <= '" + filterReportDmo.UpdatedEnd + "T23:59:59'"
		} else {
			queryFilter = "updated_at <= '" + filterReportDmo.UpdatedEnd + "T23:59:59'"
		}
	}

	if filterReportDmo.Month != "" && filterReportDmo.Year != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period = '" + filterReportDmo.Month + " " + filterReportDmo.Year + "'"
		} else {
			queryFilter = "period = '" + filterReportDmo.Month + " " + filterReportDmo.Year + "'"
		}
	}

	if filterReportDmo.Month != "" && filterReportDmo.Year == "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period LIKE '" + filterReportDmo.Month + "%'"
		} else {
			queryFilter = "period LIKE '" + filterReportDmo.Month + "%'"
		}
	}

	if filterReportDmo.Month == "" && filterReportDmo.Year != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period LIKE '%" + filterReportDmo.Year + "'"
		} else {
			queryFilter = "period LIKE '%" + filterReportDmo.Year + "'"
		}
	}

	if filterReportDmo.Quantity != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + filterReportDmo.Quantity + "%'"
		} else {
			queryFilter = "cast(quantity AS TEXT) LIKE '%" + filterReportDmo.Quantity + "%'"
		}
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateReportDmo(listReportDmo, &pagination, r.db, queryFilter)).Find(&listReportDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportDmo

	return pagination, nil
}
