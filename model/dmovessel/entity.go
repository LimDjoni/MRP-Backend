package dmovessel

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"

	"gorm.io/gorm"
)

type DmoVessel struct {
	gorm.Model
	DmoId              uint                              `json:"dmo_id"`
	Dmo                dmo.Dmo                           `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	GroupingVesselDnId uint                              `json:"grouping_vessel_dn_id"`
	GroupingVesselDn   groupingvesseldn.GroupingVesselDn `json:"grouping_vessel_dn"`
}
