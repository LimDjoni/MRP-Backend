package rkab

type RkabInput struct {
	LetterNumber    string  `form:"letter_number" json:"letter_number" validate:"required"`
	DateOfIssue     string  `form:"date_of_issue" json:"date_of_issue" gorm:"type:DATE" validate:"DateValidation"`
	Year            string  `form:"year" json:"year" validate:"required"`
	ProductionQuota float64 `form:"production_quota" json:"production_quota" validate:"required"`
	DmoObligation   float64 `form:"dmo_obligation" json:"dmo_obligation"`
}

type SortFilterRkab struct {
	Field           string
	Sort            string
	DateOfIssue     string
	Year            string
	ProductionQuota string
	Status          string
}
