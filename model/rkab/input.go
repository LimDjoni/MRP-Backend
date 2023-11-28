package rkab

type RkabInput struct {
	LetterNumber     string  `form:"letter_number" json:"letter_number" validate:"required"`
	DateOfIssue      string  `form:"date_of_issue" json:"date_of_issue" gorm:"type:DATE" validate:"DateValidation"`
	Year             string  `form:"year" json:"year" validate:"required"`
	ProductionQuota  float64 `form:"production_quota" json:"production_quota" validate:"required"`
	SalesQuota       float64 `form:"sales_quota" json:"sales_quota"`
	DmoObligation    float64 `form:"dmo_obligation" json:"dmo_obligation"`
	Year2            string  `form:"year_2" json:"year_2" validate:"required"`
	ProductionQuota2 float64 `form:"production_quota_2" json:"production_quota_2" validate:"required"`
	SalesQuota2      float64 `form:"sales_quota_2" json:"sales_quota_2"`
	DmoObligation2   float64 `form:"dmo_obligation_2" json:"dmo_obligation_2"`
	Year3            string  `form:"year_3" json:"year_3" validate:"required"`
	ProductionQuota3 float64 `form:"production_quota_3" json:"production_quota_3" validate:"required"`
	SalesQuota3      float64 `form:"sales_quota_3" json:"sales_quota_3"`
	DmoObligation3   float64 `form:"dmo_obligation_3" json:"dmo_obligation_3"`
}

type SortFilterRkab struct {
	Field           string
	Sort            string
	DateOfIssue     string
	Year            string
	ProductionQuota string
	SalesQuota      string
	Status          string
}
