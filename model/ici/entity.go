package ici

import (
	"ajebackend/model/icilevel"

	"gorm.io/gorm"
)

type Ici struct {
	gorm.Model
	Date       string            `json:"date"`
	Avarage    float64           `json:"average"`
	UnitPrice  float64           `json:"unit_price"`
	Currency   string            `json:"currency"`
	IciLevelId uint              `json:"ici_level_id`
	IciLevel   icilevel.IciLevel `json:"icilevel"`
}
