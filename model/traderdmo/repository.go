package traderdmo

import "gorm.io/gorm"

type Repository interface {
	DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error)
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
