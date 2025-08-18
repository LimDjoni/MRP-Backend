package backlog

import (
	"mrpbackend/model/unit"

	"gorm.io/gorm"
)

type BackLog struct {
	gorm.Model
	UnitId             uint     `json:"unit_id"`
	HMBreakdown        float64  `json:"hm_breakdown"`
	Problem            string   `json:"problem"`
	Component          string   `json:"component"`
	PartNumber         string   `json:"part_number"`
	PartDescription    string   `json:"part_description"`
	QtyOrder           uint     `json:"qty_order"`
	DateOfInspection   string   `json:"date_of_inspection" gorm:"DATE"`
	PlanReplaceRepair  *string  `json:"plan_replace_repair" gorm:"DATE"`
	HMReady            *float64 `json:"hm_ready"`
	PPNumber           *string  `json:"pp_number"`
	PONumber           *string  `json:"po_number"`
	Status             string   `json:"status"`
	AgingBacklogByDate int      `gorm:"->;column:aging_backlog_by_date"`

	Unit unit.Unit `gorm:"foreignKey:UnitId;references:ID" json:"Unit"`
}
