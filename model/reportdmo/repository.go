package reportdmo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetReportDmoWithPeriod(period string, iupopkId int) (ReportDmo, error)
	GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportDmoWithPeriod(period string, iupopkId int) (ReportDmo, error) {
	var reportDmo ReportDmo

	errFind := r.db.Where("period = ? AND iupopk_id = ?", period, iupopkId).First(&reportDmo).Error

	return reportDmo, errFind
}

func (r *repository) GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo, iupopkId int) (Pagination, error) {
	var listReportDmo []ReportDmo

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
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
		queryFilter = queryFilter + " AND updated_at >= '" + filterReportDmo.UpdatedStart + "'"
	}

	if filterReportDmo.UpdatedEnd != "" {
		queryFilter = queryFilter + " AND updated_at <= '" + filterReportDmo.UpdatedEnd + "T23:59:59'"
	}

	if filterReportDmo.Month != "" && filterReportDmo.Year != "" {
		queryFilter = queryFilter + " AND period = '" + filterReportDmo.Month + " " + filterReportDmo.Year + "'"
	}

	if filterReportDmo.Month != "" && filterReportDmo.Year == "" {
		queryFilter = queryFilter + " AND period LIKE '" + filterReportDmo.Month + "%'"
	}

	if filterReportDmo.Month == "" && filterReportDmo.Year != "" {
		queryFilter = queryFilter + " AND period LIKE '%" + filterReportDmo.Year + "'"
	}

	if filterReportDmo.Quantity != "" {
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + filterReportDmo.Quantity + "%'"
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateReportDmo(listReportDmo, &pagination, r.db, queryFilter)).Find(&listReportDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportDmo

	return pagination, nil
}
