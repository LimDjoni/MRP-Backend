package fuelratio

import (
	"mrpbackend/model/unit"

	"gorm.io/gorm"
)

type FuelRatio struct {
	gorm.Model
	UnitId       uint    `json:"unit_id"`
	OperatorName string  `json:"operator_name"`
	Shift        string  `json:"shift"`
	Tanggal      string  `json:"tanggal"`
	FirstHM      float64 `json:"first_hm"`
	LastHM       float64 `json:"last_hm"`
	TanggalAwal  string  `json:"tanggal_awal" gorm:"DATETIME"`
	TanggalAkhir string  `json:"tanggal_akhir" gorm:"DATETIME"`
	TotalRefill  uint    `json:"total_refill"`
	Status       bool    `json:"status"`

	Unit unit.Unit `gorm:"foreignKey:UnitId;references:ID" json:"Unit"`
}
