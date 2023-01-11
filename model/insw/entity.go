package insw

import (
	"gorm.io/gorm"
)

type Insw struct {
	gorm.Model
	IdNumber         *string `json:"id_number" gorm:"UNIQUE"`
	Periode          string  `json:"periode" gorm:"UNIQUE"`
	InswDocumentLink string  `json:"insw_document_link"`
}
