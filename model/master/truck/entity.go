package truck

import (
	"ajebackend/model/master/contractor"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type Truck struct {
	gorm.Model
	Rfid          *string               `json:"rfid" gorm:"UNIQUE"`
	Code          string                `json:"code" gorm:"UNIQUE;primaryKey"`
	NumberLambung string                `json:"number_lambung"`
	TruckModel    string                `json:"truck_model"`
	Tara          float64               `json:"tara"`
	Capacity      float64               `json:"capacity"`
	ContractorId  uint                  `json:"contractor_id"`
	Contractor    contractor.Contractor `json:"contractor"`
	CreatedById   uint                  `json:"created_by_id"`
	CreatedBy     user.User             `json:"created_by"`
	UpdatedById   uint                  `json:"updated_by_id"`
	UpdatedBy     user.User             `json:"updated_by"`
	IsActive      bool                  `json:"is_active"`
}
