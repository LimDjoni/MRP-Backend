package trader

import (
	"errors"
	"gorm.io/gorm"
)

type Repository interface {
	ListTrader() ([]Trader, error)
	CheckListTrader(list []uint) (bool, error)
	CheckEndUser(id uint) (bool, error)
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

func (r *repository) CheckListTrader(list []uint) (bool, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Where("id IN ?", list).Find(&listTrader).Error

	if getListTraderErr != nil {
		return false, getListTraderErr
	}

	if len(listTrader) != len(list) {
		return false, errors.New("record not found")
	}

	return true, nil
}

func (r *repository) CheckEndUser(id uint) (bool, error) {
	var endUser Trader

	getEndUserErr := r.db.Where("id = ?", id).First(&endUser).Error

	if getEndUserErr != nil {
		return false, getEndUserErr
	}

	return true, nil
}
