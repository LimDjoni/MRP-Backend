package bank

import (
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	NamaBank        string `json:"nama_bank"`
	NomorRekening   string `json:"nomor_rekening"`
	NamaPemilikBank string `json:"nama_pemilik_bank"`
}
