package portlocation

import (
	"gorm.io/gorm"
)

type PortLocation struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
