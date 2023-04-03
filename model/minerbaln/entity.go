package minerbaln

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type MinerbaLn struct {
	gorm.Model
	IdNumber            *string       `json:"id_number" gorm:"UNIQUE"`
	Period              string        `json:"period"`
	Quantity            float64       `json:"quantity"`
	SP3MELNDocumentLink *string       `json:"sp3meln_document_link"`
	IupopkId            uint          `json:"iupopk_id"`
	Iupopk              iupopk.Iupopk `json:"iupopk"`
}
