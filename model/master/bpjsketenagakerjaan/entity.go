package bpjsketenagakerjaan

import (
	"gorm.io/gorm"
)

type BPJSKetenagakerjaan struct {
	gorm.Model
	NomorKetenagakerjaan string `json:"nomor_ketenagakerjaan"`
}
