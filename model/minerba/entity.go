package minerba

import (
	"gorm.io/gorm"
)

type Minerba struct {
	gorm.Model
	IdNumber string `json:"id_number"`
	Period string `json:"period" gorm:"UNIQUE"`
	SP3MEDNDocumentLink *string `json:"sp3medn_document_link"`
	RecapDmoDocumentLink *string `json:"recap_dmo_document_link"`
	DetailDmoDocumentLink *string `json:"detail_dmo_document_link"`
	SP3MELNDocumentLink *string `json:"sp3meln_document_link"`
	INSWExportDocumentLink *string `json:"insw_export_document_link"`
}
