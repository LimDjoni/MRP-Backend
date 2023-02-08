package tugboat

import (
	"ajebackend/model/master/barge"

	"gorm.io/gorm"
)

type Tugboat struct {
	gorm.Model
	Name            string       `json:"name" gorm:"UNIQUE"`
	Height          float64      `json:"height"`
	Deadweight      float64      `json:"deadweight"`
	MinimumQuantity float64      `json:"minimum_quantity"`
	MaximumQuantity float64      `json:"maximum_quantity"`
	BargeId         *uint        `json:"barge_id"`
	Barge           *barge.Barge `json:"barge" gorm:"constraint:OnDelete:SET NULL;"`
}
