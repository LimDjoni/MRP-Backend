package vessel

import (
	"gorm.io/gorm"
)

type Vessel struct {
	gorm.Model
	Name            string   `json:"vessel_name" gorm:"UNIQUE"`
	Deadweight      *float64 `json:"deadweight"`
	MinimumQuantity *float64 `json:"minimum_quantity"`
	MaximumQuantity *float64 `json:"maximum_quantity"`
}
