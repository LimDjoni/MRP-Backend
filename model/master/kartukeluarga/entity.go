package kartukeluarga

import (
	"gorm.io/gorm"
)

type KartuKeluarga struct {
	gorm.Model
	NomorKartuKeluarga    string `json:"nomor_kartu_keluarga"`
	NamaIbuKandung        string `json:"nama_ibu_kandung"`
	KontakDarurat         string `json:"kontak_darurat"`
	NamaKontakDarurat     string `json:"nama_kontak_darurat"`
	HubunganKontakDarurat string `json:"hubungan_kontak_darurat"`
}
