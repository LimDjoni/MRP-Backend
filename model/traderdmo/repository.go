package traderdmo

import (
	"ajebackend/model/master/trader"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error)
	TraderListWithDmoId(idDmo int) ([]trader.Trader, trader.Trader, error)
	GetTraderEndUserDmo(idDmo int) (trader.Trader, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error) {
	var traderDmoList []TraderDmo

	findErr := r.db.Where("trader_id = ? ", idTrader).Find(&traderDmoList).Error

	if findErr != nil {
		return traderDmoList, findErr
	}

	return traderDmoList, nil
}

func (r *repository) TraderListWithDmoId(idDmo int) ([]trader.Trader, trader.Trader, error) {
	var traderDmoList []TraderDmo

	findErr := r.db.Where("dmo_id = ?", idDmo).Order("id asc").Find(&traderDmoList).Error

	var traderList []trader.Trader
	var traderEndUser trader.Trader

	if findErr != nil {
		return traderList, traderEndUser, findErr
	}

	for _, v := range traderDmoList {
		var traderTemp trader.Trader

		findTempErr := r.db.Preload(clause.Associations).Preload("Company.IndustryType").Where("id = ?", v.TraderId).First(&traderTemp).Error

		fmt.Println(traderTemp)
		if findTempErr != nil {
			return traderList, traderEndUser, findTempErr
		}

		if v.IsEndUser {
			traderEndUser = traderTemp
		} else {
			traderList = append(traderList, traderTemp)
		}
	}

	return traderList, traderEndUser, nil
}

func (r *repository) GetTraderEndUserDmo(idDmo int) (trader.Trader, error) {
	var endUserDmo TraderDmo

	errFind := r.db.Preload(clause.Associations).Preload("Trader.Company.IndustryType").Where("dmo_id = ? AND is_end_user = ?", idDmo, true).First(&endUserDmo).Error

	if errFind != nil {
		return endUserDmo.Trader, errFind
	}

	return endUserDmo.Trader, nil
}
