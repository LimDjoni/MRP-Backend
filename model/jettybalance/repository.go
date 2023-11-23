package jettybalance

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListJettyBalance(page int, sortFilter SortFilterJettyBalance, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListJettyBalance(page int, sortFilter SortFilterJettyBalance, iupopkId int) (Pagination, error) {
	var listJettyBalance []JettyBalance
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	querySort := "id desc"

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	}

	if sortFilter.JettyId != "" {
		queryFilter += " AND jetty_id = " + sortFilter.JettyId
	}

	if sortFilter.Year != "" {
		queryFilter += " AND year = '" + sortFilter.Year + "'"
	}

	errFind := r.db.Preload(clause.Associations).Order(querySort).Where(queryFilter).Scopes(paginateData(listJettyBalance, &pagination, r.db, queryFilter)).Find(&listJettyBalance).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listJettyBalance

	return pagination, nil
}
