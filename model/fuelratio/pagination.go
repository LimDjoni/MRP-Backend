package fuelratio

import (
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 7
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func paginateData(value interface{}, pagination *Pagination, db *gorm.DB, queryFilter string) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.
		Joins("JOIN units u ON fuel_ratios.unit_id = u.id").
		Joins("JOIN employees o ON fuel_ratios.employee_id = o.id").Preload(clause.Associations).Where(queryFilter).Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func paginateDataPage(value interface{}, pagination *Pagination, baseQuery *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	// Count using baseQuery (already includes JOINs and filters)
	baseQuery.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
