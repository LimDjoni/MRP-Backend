package pabeanoffice

import (
	"gorm.io/gorm"
)

type PabeanOffice struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
	Code string `json:"code" gorm:"UNIQUE"`
}
