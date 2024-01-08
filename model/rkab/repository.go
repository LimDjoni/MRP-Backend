package rkab

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error)
	DetailRkabWithYear(year int, iupopkId int) (DetailRkab, error)
	DetailRkabWithId(id int, iupopkId int) (Rkab, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error) {
	var listRkab []Rkab
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "year desc, created_at desc"

	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)

	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	if sortFilter.Field == "year" {
		sortString += ", created_at desc"
	}

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.DateOfIssue != "" {
		queryFilter += " AND date_of_issue = '" + sortFilter.DateOfIssue + "'"
	}

	if sortFilter.Year != "" {
		queryFilter += " AND year = '" + sortFilter.Year + "'"
	} else {

	}

	if sortFilter.ProductionQuota != "" {
		queryFilter = queryFilter + " AND cast(production_quota + production_quota2 + production_quota3 AS TEXT) LIKE '%" + sortFilter.ProductionQuota + "%'"
	}

	if sortFilter.SalesQuota != "" {
		queryFilter = queryFilter + " AND cast(sales_quota + sales_quota2 + sales_quota3  AS TEXT) LIKE '%" + sortFilter.SalesQuota + "%'"
	}

	if sortFilter.Status != "" {
		var status bool

		if sortFilter.Status == "Revisi" {
			status = true
		}

		if sortFilter.Status == "Non-Revisi" {
			status = false
		}

		queryFilter = fmt.Sprintf("%s AND is_revision = %v", queryFilter, status)
	}

	errFind := r.db.Table("rkabs").Preload(clause.Associations).Select("DISTINCT ON (rkabs.year) rkabs.year, rkabs.id, rkabs.dmo_obligation, rkabs.created_at, rkabs.id_number, rkabs.letter_number, rkabs.date_of_issue, rkabs.production_quota, rkabs.sales_quota, rkabs.year2, rkabs.year3 , rkabs.production_quota2, rkabs.sales_quota2, rkabs.production_quota3, rkabs.sales_quota3, rkabs.dmo_obligation2, rkabs.dmo_obligation3,  rkabs.rkab_document_link, rkabs.iupopk_id, rkabs.is_revision").Where(queryFilter).Order(sortString).Joins("left join iupopks on rkabs.iupopk_id = iupopks.id").Scopes(paginateData(listRkab, &pagination, r.db, queryFilter)).Find(&listRkab).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listRkab

	return pagination, nil
}

func (r *repository) DetailRkabWithYear(year int, iupopkId int) (DetailRkab, error) {
	var detail []Rkab

	var detailRkab DetailRkab

	filter := fmt.Sprintf("(year = '%v' OR year2 = '%v' OR year3 = '%v' ) AND iupopk_id = %v", year, year, year, iupopkId)

	errFind := r.db.Where(filter).Order("created_at desc").Find(&detail).Error

	if errFind != nil {
		return detailRkab, errFind
	}

	if len(detail) == 0 {
		return detailRkab, nil
	}

	detailRkab.ListRkab = detail

	var rkabProductionQuantity RkabProductionQuantity
	var rkabProductionQuantity2 RkabProductionQuantity
	var rkabProductionQuantity3 RkabProductionQuantity
	var rkabSalesQuantity RkabSalesQuantity
	var rkabSalesQuantity2 RkabSalesQuantity
	var rkabSalesQuantity3 RkabSalesQuantity

	filterProduction := fmt.Sprintf("production_date >= '%v-01-01' AND production_date <= '%v-12-31' AND iupopk_id = %v", detail[0].Year, detail[0].Year, iupopkId)

	errProd := r.db.Table("productions").Select("SUM(quantity) as total_production").Where(filterProduction).Scan(&rkabProductionQuantity).Error

	if errProd != nil {
		return detailRkab, errProd
	}

	filterTransaction := fmt.Sprintf("shipping_date >= '%v-01-01' AND shipping_date <= '%v-12-31' AND seller_id = %v and is_not_claim = FALSE", detail[0].Year, detail[0].Year, iupopkId)

	errTrans := r.db.Table("transactions").Select("SUM(quantity) as total_sales").Where(filterTransaction).Scan(&rkabSalesQuantity).Error

	if errTrans != nil {
		return detailRkab, errTrans
	}

	if detail[0].Year2 != nil {
		filterProduction2 := fmt.Sprintf("production_date >= '%v-01-01' AND production_date <= '%v-12-31' AND iupopk_id = %v", *detail[0].Year2, *detail[0].Year2, iupopkId)

		errProd2 := r.db.Table("productions").Select("SUM(quantity) as total_production").Where(filterProduction2).Scan(&rkabProductionQuantity2).Error

		if errProd2 != nil {
			return detailRkab, errProd2
		}

		filterTransaction2 := fmt.Sprintf("shipping_date >= '%v-01-01' AND shipping_date <= '%v-12-31' AND seller_id = %v and is_not_claim = FALSE", *detail[0].Year2, *detail[0].Year2, iupopkId)

		errTrans2 := r.db.Table("transactions").Select("SUM(quantity) as total_sales").Where(filterTransaction2).Scan(&rkabSalesQuantity2).Error

		if errTrans2 != nil {
			return detailRkab, errTrans2
		}
	}

	if detail[0].Year3 != nil {
		filterProduction3 := fmt.Sprintf("production_date >= '%v-01-01' AND production_date <= '%v-12-31' AND iupopk_id = %v", *detail[0].Year3, *detail[0].Year3, iupopkId)

		errProd3 := r.db.Table("productions").Select("SUM(quantity) as total_production").Where(filterProduction3).Scan(&rkabProductionQuantity3).Error

		if errProd3 != nil {
			return detailRkab, errProd3
		}

		filterTransaction3 := fmt.Sprintf("shipping_date >= '%v-01-01' AND shipping_date <= '%v-12-31' AND seller_id = %v and is_not_claim = FALSE", *detail[0].Year3, *detail[0].Year3, iupopkId)

		errTrans3 := r.db.Table("transactions").Select("SUM(quantity) as total_sales").Where(filterTransaction3).Scan(&rkabSalesQuantity3).Error

		if errTrans3 != nil {
			return detailRkab, errTrans3
		}
	}

	detailRkab.TotalProduction = rkabProductionQuantity.TotalProduction
	detailRkab.TotalProduction2 = rkabProductionQuantity2.TotalProduction
	detailRkab.TotalProduction3 = rkabProductionQuantity3.TotalProduction

	detailRkab.TotalSales = rkabSalesQuantity.TotalSales
	detailRkab.TotalSales2 = rkabSalesQuantity2.TotalSales
	detailRkab.TotalSales3 = rkabSalesQuantity3.TotalSales

	return detailRkab, nil
}

func (r *repository) DetailRkabWithId(id int, iupopkId int) (Rkab, error) {
	var detailRkab Rkab

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&detailRkab).Error

	if errFind != nil {
		return detailRkab, errFind
	}

	return detailRkab, nil
}
