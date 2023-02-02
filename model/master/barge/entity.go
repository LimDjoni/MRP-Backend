package barge

import (
	"gorm.io/gorm"
)

type Barge struct {
	gorm.Model
	Name            string  `json:"name" gorm:"UNIQUE"`
	Height          float64 `json:"height"`
	Deadweight      float64 `json:"deadweight"`
	MinimumQuantity float64 `json:"minimum_quantity"`
	MaximumQuantity float64 `json:"maximum_quantity"`
}
