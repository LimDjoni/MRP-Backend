package ici

import (
	"gorm.io/gorm"
)

type Ici struct {
	gorm.Model
	Date      string  `json:"date"`
	Level     string  `json:"level"`
	Avarage   float64 `json:"average"`
	UnitPrice float64 `json:"unit_price"`
	Currency  string  `json:"currency"`
}
