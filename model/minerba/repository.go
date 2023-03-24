package minerba

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetReportMinerbaWithPeriod(period string, iupopkId int) (Minerba, error)
	GetListReportMinerbaAll(page int, filterMinerba FilterAndSortMinerba, iupopkId int) (Pagination, error)
	GetDataMinerba(id int, iupopkId int) (Minerba, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportMinerbaWithPeriod(period string, iupopkId int) (Minerba, error) {
	var reportMinerba Minerba

	errFind := r.db.Where("period = ? AND iupopk_id = ?", period, iupopkId).First(&reportMinerba).Error

	return reportMinerba, errFind
}

func (r *repository) GetListReportMinerbaAll(page int, filterMinerba FilterAndSortMinerba, iupopkId int) (Pagination, error) {
	var listReportMinerba []Minerba

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
	sortFilter := "id desc"

	if filterMinerba.Field != "" && filterMinerba.Sort != "" {
		sortFilter = filterMinerba.Field + " " + filterMinerba.Sort

		if strings.ToLower(filterMinerba.Field) == "period" {
			sortFilter = "to_date(period,'Mon Year') " + filterMinerba.Sort
		}
	}

	if filterMinerba.UpdatedStart != "" {
		queryFilter = queryFilter + "AND updated_at >= '" + filterMinerba.UpdatedStart + "'"
	}

	if filterMinerba.UpdatedEnd != "" {
		queryFilter = queryFilter + " AND updated_at <= '" + filterMinerba.UpdatedEnd + "T23:59:59'"
	}

	if filterMinerba.Quantity != "" {
		quantity := fmt.Sprintf("%v", filterMinerba.Quantity)
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
	}

	if filterMinerba.Month != "" && filterMinerba.Year != "" {
		queryFilter = queryFilter + " AND period = '" + filterMinerba.Month + " " + filterMinerba.Year + "'"
	}

	if filterMinerba.Month != "" && filterMinerba.Year == "" {
		queryFilter = queryFilter + " AND period LIKE '" + filterMinerba.Month + "%'"
	}

	if filterMinerba.Month == "" && filterMinerba.Year != "" {
		queryFilter = queryFilter + " AND period LIKE '%" + filterMinerba.Year + "'"
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateMinerba(listReportMinerba, &pagination, r.db, queryFilter)).Find(&listReportMinerba).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportMinerba

	return pagination, nil
}

func (r *repository) GetDataMinerba(id int, iupopkId int) (Minerba, error) {
	var minerba Minerba

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&minerba).Error

	return minerba, errFind
}
