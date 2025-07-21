package jabatan

import (
	"mrpbackend/model/master/position"

	"gorm.io/gorm"
)

type Jabatan struct {
	gorm.Model
	EmployeeId uint   `json:"employee_id"`
	DateMove   string `json:"date_move"`
	PositionId uint   `json:"position_id"`

	Position position.Position `gorm:"foreignKey:PositionId;references:ID" json:"Position"`
}
