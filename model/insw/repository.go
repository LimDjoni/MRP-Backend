package insw

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListInsw(page int, sortFilter SortFilterInsw, iupopkId int) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListInsw(page int, sortFilter SortFilterInsw, iupopkId int) (Pagination, error) {
	var insw []Insw
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	sortString := "id desc"
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		sortString = "to_date(month,'Mon') " + sortFilter.Sort + " , year " + sortFilter.Sort
	}

	if sortFilter.Month != "" {
		queryFilter = " AND month = '" + sortFilter.Month + "'"
	}

	if sortFilter.Year != "" {
		queryFilter += " AND year = '" + sortFilter.Year + "'"
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateData(insw, &pagination, r.db, queryFilter)).Find(&insw).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = insw

	return pagination, nil
}
