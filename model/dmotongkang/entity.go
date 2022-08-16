package dmotongkang

import (
	"ajebackend/model/dmo"
	"gorm.io/gorm"
)

type DmoTongkang struct {
	gorm.Model
	DmoId uint `json:"dmo_id"`
	Dmo dmo.Dmo `json:"dmo"`
	Quantity float64 `json:"quantity"`
	Adjustment float64 `json:"adjustment"`
	GrandTotalQuantity float64 `json:"grand_total_quantity"`
}
