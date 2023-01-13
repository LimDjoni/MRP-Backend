package destination

import "gorm.io/gorm"

type Destination struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
