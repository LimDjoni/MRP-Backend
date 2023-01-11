package minerbaln

import (
	"gorm.io/gorm"
)

type MinerbaLn struct {
	gorm.Model
	IdNumber            *string `json:"id_number" gorm:"UNIQUE"`
	Period              string  `json:"period" gorm:"UNIQUE"`
	Quantity            float64 `json:"quantity"`
	SP3MELNDocumentLink *string `json:"sp3meln_document_link"`
}
