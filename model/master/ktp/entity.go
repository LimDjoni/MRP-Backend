package ktp

import (
	"gorm.io/gorm"
)

type KTP struct {
	gorm.Model
	NamaSesuaiKTP string `json:"nama_sesuai_ktp"`
	NomorKTP      string `json:"nomor_ktp"`
	TempatLahir   string `json:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir"`
	Gender        string `json:"gender"`
	Alamat        string `json:"alamat"`
	RT            string `json:"rt"`
	RW            string `json:"rw"`
	Kel           string `json:"kelurahan"`
	Kec           string `json:"kecamatan"`
	Kota          string `json:"kota"`
	Prov          string `json:"provinsi"`
	KodePos       string `json:"kode_pos"`
	GolonganDarah string `json:"golongan_darah"`
	Agama         string `json:"agama"`
	RingKTP       string `json:"ring_ktp"`
}
