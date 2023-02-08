package unit

import (
	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
	Code string `json:"code" gorm:"UNIQUE"`
}
