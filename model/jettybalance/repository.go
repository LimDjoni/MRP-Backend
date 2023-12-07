package jettybalance

import (
	"fmt"
	"strconv"

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

	var list []JettyBalance

	for _, v := range listJettyBalance {
		var production float64
		var sales float64
		var loss float64
		var quantityProduction *float64
		var quantitySales *float64
		var quantityLoss *float64

		year, err := strconv.Atoi(v.Year)
		if err != nil {
			return pagination, err
		}

		errProd := r.db.Table("productions").Select("SUM(quantity)").Where("iupopk_id = ? and jetty_id = ? and production_date <= ?", iupopkId, v.JettyId, fmt.Sprintf("%v-12-31", year-1)).Scan(&quantityProduction).Error
		if errProd != nil {
			return pagination, errProd
		}

		errSales := r.db.Table("transactions").Select("SUM(quantity)").Where("seller_id = ? and loading_port_id = ? and shipping_date <= ?", iupopkId, v.JettyId, fmt.Sprintf("%v-12-31", year-1)).Scan(&quantitySales).Error

		if errSales != nil {
			return pagination, errSales
		}

		errLoss := r.db.Table("jetty_balances").Select("SUM(total_loss)").Where("iupopk_id = ? and jetty_id = ? and cast(year AS INTEGER) < ?", iupopkId, v.JettyId, year).Scan(&quantityLoss).Error

		if errLoss != nil {
			return pagination, errLoss
		}

		if quantityProduction != nil {
			production = *quantityProduction
		}

		if quantitySales != nil {
			sales = *quantitySales
		}

		if quantityLoss != nil {
			loss = *quantityLoss
		}

		v.StartBalance = production - sales - loss

		list = append(list, v)
	}

	pagination.Data = list

	return pagination, nil
}
