package dmovessel

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"

	"gorm.io/gorm"
)

type DmoVessel struct {
	gorm.Model
	DmoId uint    `json:"dmo_id"`
	Dmo   dmo.Dmo `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	// Start Will be deleted
	VesselName         string  `json:"vessel_name"`
	Quantity           float64 `json:"quantity"`
	Adjustment         float64 `json:"adjustment"`
	GrandTotalQuantity float64 `json:"grand_total_quantity"`
	BlDate             string  `json:"bl_date" gorm:"type:DATE"`
	// End Will be deleted
	GroupingVesselDnId *uint                              `json:"grouping_vessel_dn_id"`
	GroupingVesselDn   *groupingvesseldn.GroupingVesselDn `json:"grouping_vessel_dn"`
}
