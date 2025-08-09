package stockfuel

import (
	"gorm.io/gorm"
)

type StockFuel struct {
	gorm.Model
	Date           string  `json:"date"`
	FirstStock     float64 `json:"first_stock"`
	Day            float64 `json:"day"`
	Night          float64 `json:"night"`
	Total          float64 `json:"total"`
	GrandTotal     float64 `json:"grand_total"`
	FuelIn         float64 `json:"fuel_in"`
	EndStock       float64 `json:"end_stock"`
	MtdConsump     float64 `json:"mtd_consump"`
	PlanPermintaan float64 `json:"plan_permintaan"`
	MJSU           float64 `json:"mjsu"`
	PPP            float64 `json:"ppp"`
	SADP           float64 `json:"sadp"`
	BTP            float64 `json:"btp"`
	SURPLUS        float64 `json:"surplus"`
}
