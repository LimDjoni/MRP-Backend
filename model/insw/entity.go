package insw

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Insw struct {
	gorm.Model
	IdNumber         *string       `json:"id_number" gorm:"UNIQUE"`
	Month            string        `json:"month"`
	Year             int           `json:"year"`
	InswDocumentLink string        `json:"insw_document_link"`
	IupopkId         uint          `json:"iupopk_id"`
	Iupopk           iupopk.Iupopk `json:"iupopk"`
}
