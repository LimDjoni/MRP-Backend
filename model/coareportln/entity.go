package coareportln

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type CoaReportLn struct {
	gorm.Model
	IdNumber                string        `json:"id_number" gorm:"UNIQUE"`
	DateFrom                string        `json:"date_from" gorm:"type:DATE"`
	DateTo                  string        `json:"date_to" gorm:"type:DATE"`
	CoaReportLnDocumentLink string        `json:"coa_report_ln_document_link"`
	IupopkId                uint          `json:"iupopk_id"`
	Iupopk                  iupopk.Iupopk `json:"iupopk"`
}
