package reportdmo

import (
	"gorm.io/gorm"
)

type ReportDmo struct {
	gorm.Model
	IdNumber              *string `json:"id_number" gorm:"UNIQUE"`
	Period                string  `json:"period" gorm:"UNIQUE"`
	Quantity              float64 `json:"quantity"`
	RecapDmoDocumentLink  *string `json:"recap_dmo_document_link"`
	DetailDmoDocumentLink *string `json:"detail_dmo_document_link"`
}
