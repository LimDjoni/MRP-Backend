package laporan

import (
	"gorm.io/gorm"
)

type Laporan struct {
	gorm.Model
	RingSerapan             string `json:"ring_serapan"`
	RingRIPPM               string `json:"ring_rippm"`
	KategoriLaporanTwiwulan string `json:"kategori_laporan_twiwulan"`
	KategoriLokalNonLokal   string `json:"kategori_lokal_non_lokal"`
	Rekomendasi             string `json:"rekomendasi"`
}
