package transactionrequestreport

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type TransactionRequestReport struct {
	gorm.Model
	IdNumber       *string       `json:"id_number" gorm:"UNIQUE"`
	DateFrom       string        `json:"date_from" gorm:"type:DATE"`
	DateTo         string        `json:"date_to" gorm:"type:DATE"`
	DocumentDnLink string        `json:"document_dn_link"`
	DocumentLnLink string        `json:"document_ln_link"`
	IupopkId       uint          `json:"iupopk_id"`
	Iupopk         iupopk.Iupopk `json:"iupopk"`
}
