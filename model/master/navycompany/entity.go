package navycompany

import (
	"gorm.io/gorm"
)

type NavyCompany struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
