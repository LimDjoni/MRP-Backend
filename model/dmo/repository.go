package dmo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetReportDmoWithPeriod(period string, iupopkId int) (Dmo, error)
	GetListReportDmoAll(page int, filterDmo FilterAndSortDmo, iupopkId int) (Pagination, error)
	GetDataDmo(id int, iupopkId int) (Dmo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetReportDmoWithPeriod(period string, iupopkId int) (Dmo, error) {
	var reportDmo Dmo

	errFind := r.db.Preload(clause.Associations).Where("period = ? AND iupopk_id = ?", period, iupopkId).First(&reportDmo).Error

	return reportDmo, errFind
}

func (r *repository) GetListReportDmoAll(page int, filterDmo FilterAndSortDmo, iupopkId int) (Pagination, error) {
	var listReportDmo []Dmo

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("AND iupopk_id = %v", iupopkId)
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

func (r *repository) GetDataDmo(id int, iupopkId int) (Dmo, error) {
	var dmo Dmo

	errFind := r.db.Preload(clause.Associations).Where("id = ? AND iupopkId = ?", id, iupopkId).First(&dmo).Error

	return dmo, errFind
}
