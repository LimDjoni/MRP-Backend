package dmo

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Dmo struct {
	gorm.Model
	IdNumber string `json:"id_number"`
	Type string `json:"type"`
	Date string `json:"date" gorm:"type:DATE"`
	TongkangTotalQuantity float64 `json:"tongkang_total_quantity"`
	TongkangAdjustment float64 `json:"tongkang_adjustment"`
	TongkangGrandTotalQuantity float64 `json:"tongkang_grand_total_quantity"`
	VesselTotalQuantity float64 `json:"vessel_total_quantity"`
	VesselAdjustment float64 `json:"vessel_adjustment"`
	VesselGrandTotalQuantity float64 `json:"vessel_grand_total_quantity"`
	EndUser pgtype.JSONB `json:"end_user"`
	ReconciliationLetterDocument string `json:"berita_acara_document"`
	BASTDocument string `json:"bast_document"`
	StatementLetterDocument string `json:"statement_letter_document"`
}
