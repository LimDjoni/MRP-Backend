package trader

import "gorm.io/gorm"

type Trader struct {
	gorm.Model
	TraderName  string `json:"trader_name"`
	Position    string `json:"position"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber   string `json:"fax_number"`
	CompanyName string `json:"company_name"`
}
