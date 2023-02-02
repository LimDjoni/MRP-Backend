package surveyor

import (
	"gorm.io/gorm"
)

type Surveyor struct {
	gorm.Model
	Name string `json:"name" gorm:"UNIQUE"`
}
