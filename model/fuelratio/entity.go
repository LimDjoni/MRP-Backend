package fuelratio

import (
	"mrpbackend/model/employee"
	"mrpbackend/model/unit"

	"gorm.io/gorm"
)

type FuelRatio struct {
	gorm.Model
	UnitId      uint   `json:"unit_id"`
	EmployeeId  uint   `json:"employee_id"`
	Shift       string `json:"shift"`
	FirstHM     string `json:"first_hm" gorm:"DATETIME"`
	LastHM      string `json:"last_hm" gorm:"DATETIME"`
	TotalRefill uint   `json:"total_refill"`
	Status      bool   `json:"status"`

	Unit     unit.Unit         `gorm:"foreignKey:UnitId;references:ID" json:"Unit"`
	Employee employee.Employee `gorm:"foreignKey:EmployeeId;references:ID" json:"Employee"`
}
