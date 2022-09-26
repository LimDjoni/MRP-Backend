package trader

import "gorm.io/gorm"

type Repository interface {
	ListTrader() ([]Trader, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListTrader() ([]Trader, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Find(&listTrader).Error

	return listTrader, getListTraderErr
}
