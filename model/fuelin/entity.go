package fuelin

import (
	"gorm.io/gorm"
)

type FuelIn struct {
	gorm.Model
	Date            string  `json:"date"`
	Vendor          string  `json:"vendor"`
	Code            string  `json:"code"`
	NomorSuratJalan string  `json:"nomor_surat_jalan"`
	NomorPlatMobil  string  `json:"nomor_plat_mobil"`
	Qty             float64 `json:"qty"`
	QtyNow          float64 `json:"qty_now"`
	Driver          string  `json:"driver"`
	TujuanAwal      string  `json:"tujuan_awal"`
}
