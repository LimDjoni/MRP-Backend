package insw

import (
	"gorm.io/gorm"
)

type Insw struct {
	gorm.Model
	IdNumber         *string `json:"id_number" gorm:"UNIQUE"`
	Month            string  `json:"month"`
	Year             int     `json:"year"`
	InswDocumentLink string  `json:"insw_document_link"`
}
