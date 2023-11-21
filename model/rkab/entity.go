package rkab

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

// Year -> 1st year created rkab (3 years period)
// Year2 -> 2nd year
// Year3 -> 3rd year (last period)

// data without number -> data 1st year
// data with 2 -> data 2nd year
// data with 3 3 -> data 3rd year

type Rkab struct {
	gorm.Model
	IdNumber         string        `json:"id_number" gorm:"UNIQUE"`
	LetterNumber     string        `json:"letter_number"`
	DateOfIssue      string        `json:"date_of_issue" gorm:"type:DATE"`
	Year             string        `json:"year"`
	ProductionQuota  float64       `json:"production_quota"`
	SalesQuota       float64       `json:"sales_quota"`
	DmoObligation    float64       `json:"dmo_obligation"`
	Year2            string        `json:"year_2"`
	SalesQuota2      float64       `json:"sales_quota_2"`
	ProductionQuota2 float64       `json:"production_quota_2"`
	DmoObligation2   float64       `json:"dmo_obligation_2"`
	Year3            string        `json:"year_3"`
	ProductionQuota3 float64       `json:"production_quota_3"`
	SalesQuota3      float64       `json:"sales_quota_3"`
	DmoObligation3   float64       `json:"dmo_obligation_3"`
	RkabDocumentLink string        `json:"rkab_document_link"`
	IsRevision       bool          `json:"is_revision"`
	IupopkId         uint          `json:"iupopk_id"`
	Iupopk           iupopk.Iupopk `json:"iupopk"`
}
