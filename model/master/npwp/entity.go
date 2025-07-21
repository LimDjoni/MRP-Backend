package npwp

import (
	"gorm.io/gorm"
)

type NPWP struct {
	gorm.Model
	NomorNPWP   *string `json:"nomor_npwp"`
	StatusPajak string  `json:"status_pajak"`
}
