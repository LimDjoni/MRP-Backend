package trader

import (
	"ajebackend/model/master/company"

	"gorm.io/gorm"
)

type Trader struct {
	gorm.Model
	TraderName string          `json:"trader_name"`
	Position   string          `json:"position"`
	Email      *string         `json:"email"`
	CompanyId  uint            `json:"company_id"`
	Company    company.Company `json:"company"`
}
