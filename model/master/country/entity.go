package country

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name string `json:"name"`
	Code string `json:"code" gorm:"UNIQUE"`
}
