package currency

import (
	"gorm.io/gorm"
)

type Currency struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
	Code string `json:"code" gorm:"UNIQUE"`
}
