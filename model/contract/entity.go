package contract

import (
	"ajebackend/model/master/company"
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model
	ContractDate   string          `json:"contract_date" gorm:"DATETIME"`
	ContractNumber string          `json:"contract_number" gorm:"UNIQUE"`
	CustomerId     uint            `json:"customer_id"`
	Customer       company.Company `json:"customer"`
	Quantity       float64         `json:"quantity"`
	Validity       string          `json:"validity" gorm:"DATETIME"`
	File           *string         `json:"file"`
	IupopkId       uint            `json:"iupopk_id"`
	Iupopk         iupopk.Iupopk   `json:"iupopk"`
}
