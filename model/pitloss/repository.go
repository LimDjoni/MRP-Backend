package pitloss

import (
	"ajebackend/model/jettybalance"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	DetailJettyBalance(id int, iupopkId int) (OutputJettyBalancePitLossDetail, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) DetailJettyBalance(id int, iupopkId int) (OutputJettyBalancePitLossDetail, error) {
	var detailOutput OutputJettyBalancePitLossDetail

	var jettyBalance jettybalance.JettyBalance
	var pitLoss []PitLoss

	errFindJettyBalance := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&jettyBalance).Error

	if errFindJettyBalance != nil {
		return detailOutput, errFindJettyBalance
	}

	errFindPitLoss := r.db.Preload(clause.Associations).Where("jetty_balance_id = ?", jettyBalance.ID).Order("id asc").Find(&pitLoss).Error

	if errFindPitLoss != nil {
		return detailOutput, errFindPitLoss
	}

	var production float64
	var sales float64
	var loss float64
	var quantityProduction *float64
	var quantitySales *float64
	var quantityLoss *float64

	year, err := strconv.Atoi(jettyBalance.Year)
	if err != nil {
		return detailOutput, err
	}

	errProd := r.db.Table("productions").Select("SUM(quantity)").Where("iupopk_id = ? and jetty_id = ? and production_date <= ?", iupopkId, jettyBalance.JettyId, fmt.Sprintf("%v-12-31", year-1)).Scan(&quantityProduction).Error
	if errProd != nil {
		return detailOutput, errProd
	}

	errSales := r.db.Table("transactions").Select("SUM(quantity)").Where("seller_id = ? and loading_port_id = ? and shipping_date <= ?", iupopkId, jettyBalance.JettyId, fmt.Sprintf("%v-12-31", year-1)).Scan(&quantitySales).Error

	if errSales != nil {
		return detailOutput, errSales
	}

	errLoss := r.db.Table("jetty_balances").Select("SUM(total_loss)").Where("iupopk_id = ? and jetty_id = ? and cast(year AS INTEGER) < ?", iupopkId, jettyBalance.JettyId, year).Scan(&quantityLoss).Error

	if errLoss != nil {
		return detailOutput, errLoss
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

	jettyBalance.StartBalance = production - sales - loss

	detailOutput.JettyBalance = jettyBalance
	detailOutput.PitLoss = pitLoss

	return detailOutput, nil
}
