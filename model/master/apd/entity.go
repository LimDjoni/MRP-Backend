package apd

import (
	"gorm.io/gorm"
)

type APD struct {
	gorm.Model
	UkuranBaju   string `json:"ukuran_baju"`
	UkuranCelana string `json:"ukuran_celana"`
	UkuranSepatu string `json:"ukuran_sepatu"`
}
