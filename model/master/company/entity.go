package company

import (
	"ajebackend/model/master/industrytype"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	CompanyName    string                     `json:"company_name"`
	IndustryTypeId *uint                      `json:"industry_type_id"`
	IndustryType   *industrytype.IndustryType `json:"industry_type"`
	Address        string                     `json:"address"`
	Province       string                     `json:"province"`
	PhoneNumber    string                     `json:"phone_number"`
	FaxNumber      string                     `json:"fax_number"`
	IsTrader       bool                       `json:"is_trader"`
	IsEndUser      bool                       `json:"is_end_user"`
}
