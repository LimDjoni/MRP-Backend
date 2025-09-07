package backlog

import (
	"mrpbackend/model/unit"

	"gorm.io/gorm"
)

type BackLog struct {
	gorm.Model
	UnitId             uint          `json:"unit_id"`
	HMBreakdown        float64       `json:"hm_breakdown"`
	Problem            string        `json:"problem"`
	Component          string        `json:"component"`
	DateOfInspection   string        `json:"date_of_inspection" gorm:"DATE"`
	PlanReplaceRepair  *string       `json:"plan_replace_repair" gorm:"DATE"`
	HMReady            *float64      `json:"hm_ready"`
	PPNumber           *string       `json:"pp_number"`
	PONumber           *string       `json:"po_number"`
	Status             string        `json:"status"`
	Parts              []BackLogPart `gorm:"foreignKey:BackLogID" json:"parts"`
	AgingBacklogByDate int           `gorm:"->;column:aging_backlog_by_date"`

	Unit unit.Unit `gorm:"foreignKey:UnitId;references:ID" json:"Unit"`
}

// BackLogPart represents individual parts linked to a BackLog
type BackLogPart struct {
	gorm.Model
	BackLogID       uint   `json:"back_log_id"`
	PartNumber      string `json:"part_number"`
	PartDescription string `json:"part_description"`
	QtyOrder        uint   `json:"qty_order"`
}
