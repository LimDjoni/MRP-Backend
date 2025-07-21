package position

import (
	"gorm.io/gorm"
)

type Position struct {
	gorm.Model
	PositionName string `json:"position_name"`
}
