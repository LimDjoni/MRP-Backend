package royaltyrecon

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type RoyaltyRecon struct {
	gorm.Model
	IdNumber                 string        `json:"id_number" gorm:"UNIQUE"`
	DateFrom                 string        `json:"date_from" gorm:"type:DATE"`
	DateTo                   string        `json:"date_to" gorm:"type:DATE"`
	RoyaltyReconDocumentLink string        `json:"royalty_recon_document_link"`
	IupopkId                 uint          `json:"iupopk_id"`
	Iupopk                   iupopk.Iupopk `json:"iupopk"`
}

type SortFilterRoyaltyRecon struct {
	Field     string
	Sort      string
	DateStart string
	DateEnd   string
}
