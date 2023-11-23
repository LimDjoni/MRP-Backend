package pitloss

import (
	"ajebackend/model/jettybalance"

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

	detailOutput.JettyBalance = jettyBalance
	detailOutput.PitLoss = pitLoss

	return detailOutput, nil
}
