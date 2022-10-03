package traderdmo

import (
	"ajebackend/model/dmo"
	"ajebackend/model/trader"
	"gorm.io/gorm"
)

type TraderDmo struct {
	gorm.Model
	DmoId uint `json:"dmo_id"`
	Dmo dmo.Dmo `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	TraderId uint `json:"trader_id"`
	Trader trader.Trader `json:"trader"`
	Order int `json:"order"`
	IsEndUser bool `json:"is_end_user"`
}
