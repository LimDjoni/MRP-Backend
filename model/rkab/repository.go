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
		queryFilter = queryFilter + " AND cast(production_quota AS TEXT) LIKE '%" + sortFilter.ProductionQuota + "%'"
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

	errFind := r.db.Table("rkabs").Preload(clause.Associations).Select("DISTINCT ON (rkabs.year) rkabs.year, rkabs.id, rkabs.created_at, rkabs.id_number, rkabs.letter_number, rkabs.date_of_issue, rkabs.production_quota, rkabs.rkab_document_link, rkabs.iupopk_id, rkabs.is_revision").Where(queryFilter).Order(sortString).Joins("left join iupopks on rkabs.iupopk_id = iupopks.id").Scopes(paginateData(listRkab, &pagination, r.db, queryFilter)).Find(&listRkab).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listRkab

	return pagination, nil
}

func (r *repository) DetailRkabWithYear(year int, iupopkId int) (DetailRkab, error) {
	var listRkabDetail []Rkab

	var detailRkab DetailRkab

	filter := fmt.Sprintf("year = '%v' AND iupopk_id = %v", year, iupopkId)

	errFind := r.db.Where(filter).Order("created_at desc").Find(&listRkabDetail).Error

	if errFind != nil {
		return detailRkab, errFind
	}

	detailRkab.ListRkab = listRkabDetail

	var rkabProductionQuantity RkabProductionQuantity

	filterProduction := fmt.Sprintf("production_date >= '%v-01-01' AND production_date <= '%v-12-31' AND iupopk_id = %v", year, year, iupopkId)

	errProd := r.db.Table("productions").Select("SUM(quantity) as total_production").Where(filterProduction).Scan(&rkabProductionQuantity).Error

	if errProd != nil {
		return detailRkab, errProd
	}

	detailRkab.TotalProduction = rkabProductionQuantity.TotalProduction

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
