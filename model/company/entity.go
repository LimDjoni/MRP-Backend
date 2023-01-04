package company

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyName  string `json:"company_name"`
	IndustryType string `json:"industry_type"`
	Address      string `json:"address"`
	Province     string `json:"province"`
	PhoneNumber  string `json:"phone_number"`
	FaxNumber    string `json:"fax_number"`
}
