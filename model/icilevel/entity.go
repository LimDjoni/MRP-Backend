package icilevel

import (
	"gorm.io/gorm"
)

type IciLevel struct {
	gorm.Model
	Name        string `json:"level"`
	Description string `json:"description"`
}
