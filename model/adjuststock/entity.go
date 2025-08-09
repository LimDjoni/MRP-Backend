package adjuststock

import (
	"gorm.io/gorm"
)

type AdjustStock struct {
	gorm.Model
	Date  string  `json:"date"`
	Stock float64 `json:"stock"`
}
