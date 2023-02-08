package insw

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListInsw(page int, sortFilter SortFilterInsw) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListInsw(page int, sortFilter SortFilterInsw) (Pagination, error) {
	var insw []Insw
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	sortString := "id desc"
	queryFilter := ""

	if sortFilter.SortMonth != "" {
		sortString = "to_date(month,'Month') " + sortFilter.SortMonth
	}

	if sortFilter.SortYear != "" {
		if sortString == "id desc" {
			sortString = "year " + sortFilter.SortYear
		} else {
			sortString += ", year " + sortFilter.SortYear
		}
	}

	if sortFilter.Month != "" {
		queryFilter = "month = '" + sortFilter.Month + "'"
	}

	if sortFilter.Year != "" {
		if queryFilter != "" {
			queryFilter += "AND year = '" + sortFilter.Year + "'"
		} else {
			queryFilter = "year = '" + sortFilter.Year + "'"
		}
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(insw, &pagination, r.db, queryFilter)).Find(&insw).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = insw

	return pagination, nil
}
