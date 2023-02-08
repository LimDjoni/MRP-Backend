package salessystem

import (
	"gorm.io/gorm"
)

type SalesSystem struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
