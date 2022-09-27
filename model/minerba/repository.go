package minerba

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetReportMinerbaWithPeriod(period string) (Minerba, error)
	GetListReportMinerbaAll(page int) (Pagination, error)
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

func(r *repository) GetListReportMinerbaAll(page int) (Pagination, error) {
	var listReportMinerba []Minerba

	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	errFind := r.db.Scopes(paginateMinerba(listReportMinerba, &pagination, r.db)).Find(&listReportMinerba).Error

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

