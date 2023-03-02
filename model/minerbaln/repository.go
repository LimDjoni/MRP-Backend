package minerbaln

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetReportMinerbaLnWithPeriod(period string) (MinerbaLn, error)
	GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn) (Pagination, error)
	GetDataMinerbaLn(id int) (MinerbaLn, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportMinerbaLnWithPeriod(period string) (MinerbaLn, error) {
	var reportMinerbaLn MinerbaLn

	errFind := r.db.Where("period = ?", period).First(&reportMinerbaLn).Error

	return reportMinerbaLn, errFind
}

func (r *repository) GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn) (Pagination, error) {
	var listReportMinerbaLn []MinerbaLn

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := ""
	sortFilter := "id desc"

	if filterMinerbaLn.Field != "" && filterMinerbaLn.Sort != "" {
		sortFilter = filterMinerbaLn.Field + " " + filterMinerbaLn.Sort

		if strings.ToLower(filterMinerbaLn.Field) == "period" {
			sortFilter = "to_date(period,'Mon Year') " + filterMinerbaLn.Sort
		}
	}

	if filterMinerbaLn.UpdatedStart != "" {
		queryFilter = queryFilter + "updated_at >= '" + filterMinerbaLn.UpdatedStart + "'"
	}

	if filterMinerbaLn.UpdatedEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND updated_at <= '" + filterMinerbaLn.UpdatedEnd + "T23:59:59'"
		} else {
			queryFilter = "updated_at <= '" + filterMinerbaLn.UpdatedEnd + "T23:59:59'"
		}
	}

	if filterMinerbaLn.Quantity != "" {
		quantity := fmt.Sprintf("%v", filterMinerbaLn.Quantity)
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
		} else {
			queryFilter = "cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
		}
	}

	if filterMinerbaLn.Month != "" && filterMinerbaLn.Year != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period = '" + filterMinerbaLn.Month + " " + filterMinerbaLn.Year + "'"
		} else {
			queryFilter = "period = '" + filterMinerbaLn.Month + " " + filterMinerbaLn.Year + "'"
		}
	}

	if filterMinerbaLn.Month != "" && filterMinerbaLn.Year == "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period LIKE '" + filterMinerbaLn.Month + "%'"
		} else {
			queryFilter = "period LIKE '" + filterMinerbaLn.Month + "%'"
		}
	}

	if filterMinerbaLn.Month == "" && filterMinerbaLn.Year != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND period LIKE '%" + filterMinerbaLn.Year + "'"
		} else {
			queryFilter = "period LIKE '%" + filterMinerbaLn.Year + "'"
		}
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateMinerbaLn(listReportMinerbaLn, &pagination, r.db, queryFilter)).Find(&listReportMinerbaLn).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportMinerbaLn

	return pagination, nil
}

func (r *repository) GetDataMinerbaLn(id int) (MinerbaLn, error) {
	var minerbaLn MinerbaLn

	errFind := r.db.Where("id = ?", id).First(&minerbaLn).Error

	return minerbaLn, errFind
}
