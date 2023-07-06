package industrytype

import (
	"gorm.io/gorm"
)

type IndustryType struct {
	gorm.Model
	Name           string `json:"name" gorm:"UNIQUE"`
	Category       string `json:"category"`
	SystemCategory string `json:"system_category"`
}
