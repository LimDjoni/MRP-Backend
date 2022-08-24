package minerba

import (
	"gorm.io/gorm"
	"time"
)

type Minerba struct {
	gorm.Model
	IdNumber string `json:"id_number"`
	Date time.Time `json:"date" gorm:"DATE"`
	Periode string `json:"periode"`
	SP3MEDNDocumentLink string `json:"sp3medn_document_link"`
	RekapDmoDocumentLink string `json:"rekap_dmo_document_link"`
	RincianDmoDocumentLink string `json:"rincian_dmo_document_link"`
	SP3MELNDocumentLink string `json:"sp3meln_document_link"`
	INSWEksporDocumentLink string `json:"insw_ekspor_document_link"`
}
