package company

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyName  string `json:"company_name"`
	Address     string	`json:"address"`
	Email       string	`json:"email"`
	Province	string	`json:"province"`
	PhoneNumber	string	`json:"phone_number"`
	FaxNumber	string	`json:"fax_number"`
}
