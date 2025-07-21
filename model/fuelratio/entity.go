package fuelratio

import (
	"mrpbackend/model/employee"
	"mrpbackend/model/unit"

	"gorm.io/gorm"
)

type FuelRatio struct {
	gorm.Model
	UnitId       uint    `json:"unit_id"`
	EmployeeId   uint    `json:"employee_id"`
	Shift        string  `json:"shift"`
	Tanggal      string  `json:"tanggal"`
	FirstHM      float64 `json:"first_hm"`
	LastHM       float64 `json:"last_hm"`
	TanggalAwal  string  `json:"tanggal_awal" gorm:"DATETIME"`
	TanggalAkhir string  `json:"tanggal_akhir" gorm:"DATETIME"`
	TotalRefill  uint    `json:"total_refill"`
	Status       bool    `json:"status"`

	Unit     unit.Unit         `gorm:"foreignKey:UnitId;references:ID" json:"Unit"`
	Employee employee.Employee `gorm:"foreignKey:EmployeeId;references:ID" json:"Employee"`
}
