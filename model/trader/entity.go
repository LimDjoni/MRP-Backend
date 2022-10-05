package trader

import (
	"ajebackend/model/company"
	"gorm.io/gorm"
)

type Trader struct {
	gorm.Model
	TraderName  string `json:"trader_name"`
	Position    string `json:"position"`
	CompanyId	uint	`json:"company_id"`
	Company		company.Company `json:"company"`
}
