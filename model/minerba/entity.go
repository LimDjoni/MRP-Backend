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
	SP3MEDNDocument string `json:"sp3medn_document"`
	RekapDmoDocument string `json:"rekap_dmo_document"`
	RincianDmoDocument string `json:"rincian_dmo_document"`
	SP3MELNDocument string `json:"sp3meln_document"`
	INSWEksporDocument string `json:"insw_ekspor_document"`
}
