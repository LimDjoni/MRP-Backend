package rkab

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Rkab struct {
	gorm.Model
	IdNumber         string        `json:"id_number" gorm:"UNIQUE"`
	LetterNumber     string        `json:"letter_number"`
	DateOfIssue      string        `json:"date_of_issue" gorm:"type:DATE"`
	Year             string        `json:"year"`
	ProductionQuota  float64       `json:"production_quota"`
	SalesQuota       float64       `json:"sales_quota"`
	DmoObligation    float64       `json:"dmo_obligation"`
	RkabDocumentLink string        `json:"rkab_document_link"`
	IsRevision       bool          `json:"is_revision"`
	IupopkId         uint          `json:"iupopk_id"`
	Iupopk           iupopk.Iupopk `json:"iupopk"`
}
