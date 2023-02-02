package navyship

import (
	"gorm.io/gorm"
)

type NavyShip struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
