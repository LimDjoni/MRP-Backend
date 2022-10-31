package minerba

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	GetReportMinerbaWithPeriod(period string) (Minerba, error)
	GetListReportMinerbaAll(page int, filterMinerba FilterMinerba) (Pagination, error)
	GetDataMinerba(id int) (Minerba, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func(r *repository) GetReportMinerbaWithPeriod(period string) (Minerba, error) {
	var reportMinerba Minerba

	errFind := r.db.Where("period = ?", period).First(&reportMinerba).Error

	return reportMinerba, errFind
}

func(r *repository) GetListReportMinerbaAll(page int, filterMinerba FilterMinerba) (Pagination, error) {
	var listReportMinerba []Minerba

	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	queryFilter := ""

	if filterMinerba.CreatedStart != "" {
		queryFilter = queryFilter + "created_at >= '" + filterMinerba.CreatedStart + "'"
	}

	if filterMinerba.CreatedEnd != "" {
		if queryFilter != "" {
			queryFilter = queryFilter + " AND created_at <= '" + filterMinerba.CreatedEnd + "T23:59:59'"
		} else {
			queryFilter = "created_at <= '" + filterMinerba.CreatedEnd + "T23:59:59'"
		}
	}

	if filterMinerba.Quantity != 0 {
		quantity := fmt.Sprintf("%v", filterMinerba.Quantity)
		if queryFilter != "" {
			queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" +  quantity + "%'"
		} else {
			queryFilter = "cast(quantity AS TEXT) LIKE '%" +  quantity + "%'"
		}
	}

	errFind := r.db.Where(queryFilter).Scopes(paginateMinerba(listReportMinerba, &pagination, r.db, queryFilter)).Find(&listReportMinerba).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportMinerba

	return pagination, nil
}

func(r *repository) GetDataMinerba(id int) (Minerba, error) {
	var minerba Minerba

	errFind := r.db.Where("id = ?", id).First(&minerba).Error

	return minerba, errFind
}

