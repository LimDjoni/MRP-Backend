package bpjskesehatan

import (
	"gorm.io/gorm"
)

type BPJSKesehatan struct {
	gorm.Model
	NomorKesehatan string `json:"nomor_kesehatan"`
}
