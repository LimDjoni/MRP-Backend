package minerba

import (
	"gorm.io/gorm"
	"time"
)

type Minerba struct {
	gorm.Model
	IdNumber string `json:"id_number"`
	Date time.Time `json:"date" gorm:"DATE"`
	Period string `json:"period"`
	SP3MEDNDocumentLink string `json:"sp3medn_document_link"`
	RecapDmoDocumentLink string `json:"recap_dmo_document_link"`
	DetailDmoDocumentLink string `json:"detail_dmo_document_link"`
	SP3MELNDocumentLink string `json:"sp3meln_document_link"`
	INSWExportDocumentLink string `json:"insw_export_document_link"`
}
