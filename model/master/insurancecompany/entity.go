package insurancecompany

import (
	"gorm.io/gorm"
)

type InsuranceCompany struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
