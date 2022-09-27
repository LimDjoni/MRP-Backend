package dmo

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetReportDmoWithPeriod(period string) (Dmo, error)
	GetListReportDmoAll(page int) (Pagination, error)
	GetDataDmo(id int) (Dmo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func(r *repository) GetReportDmoWithPeriod(period string) (Dmo, error) {
	var reportDmo Dmo

	errFind := r.db.Where("period = ?", period).First(&reportDmo).Error

	return reportDmo, errFind
}

func(r *repository) GetListReportDmoAll(page int) (Pagination, error) {
	var listReportDmo []Dmo

	var pagination Pagination
	pagination.Limit = 10
	pagination.Page = page
	errFind := r.db.Scopes(paginateDmo(listReportDmo, &pagination, r.db)).Find(&listReportDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportDmo

	return pagination, nil
}

func(r *repository) GetDataDmo(id int) (Dmo, error) {
	var dmo Dmo

	errFind := r.db.Where("id = ?", id).First(&dmo).Error

	return dmo, errFind
}

