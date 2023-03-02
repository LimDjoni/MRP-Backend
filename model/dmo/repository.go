package dmo

import (
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
			sortFilter = "to_date(a.period,'Mon Year') " + filterDmo.Sort
		}

		if strings.ToLower(filterDmo.Field) == "grand_total_quantity" {
			sortFilter = "a.vessel_grand_total_quantity + a.barge_grand_total_quantity " + filterDmo.Sort
		}

		if strings.ToLower(filterDmo.Field) == "buyer_id" {
			sortFilter = "d.company_name " + filterDmo.Sort
		}
	}

	if filterDmo.Quantity != "" {
		queryFilter += " AND cast(a.vessel_grand_total_quantity + a.barge_grand_total_quantity AS TEXT) LIKE '%" + filterDmo.Quantity + "%'"
	}

	if filterDmo.Month != "" && filterDmo.Year != "" {
		queryFilter = queryFilter + " AND a.period = '" + filterDmo.Month + " " + filterDmo.Year + "'"
	}

	if filterDmo.Month != "" && filterDmo.Year == "" {
		queryFilter = queryFilter + " AND a.period LIKE '" + filterDmo.Month + "%'"
	}

	if filterDmo.Month == "" && filterDmo.Year != "" {
		queryFilter = queryFilter + " AND a.period LIKE '%" + filterDmo.Year + "'"
	}

	if filterDmo.BuyerId != "" {
		queryFilter = queryFilter + " AND d.id = " + filterDmo.BuyerId
	}

	var listDmo []map[string]interface{}

	var rawQuery = `select a.*, d.company_name from  dmos a
LEFT JOIN trader_dmos b on a.id = b.dmo_id
LEFT JOIN traders c on b.trader_id = c.id
LEFT JOIN companies d on c.company_id = d.id
Where b.is_end_user = TRUE`

	if queryFilter != "" {
		rawQuery += queryFilter
	}
	rawQuery += ` order by ` + sortFilter

	errFind := r.db.Scopes(paginateDmo(listReportDmo, &pagination, r.db, queryFilter)).Raw(rawQuery).Scan(&listDmo).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listDmo

	return pagination, nil
}

func (r *repository) GetDataDmo(id int) (Dmo, error) {
	var dmo Dmo

	errFind := r.db.Where("id = ?", id).First(&dmo).Error

	return dmo, errFind
}
