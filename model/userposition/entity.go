package userposition

import (
	"mrpbackend/model/employee"
	"mrpbackend/model/master/position"

	"gorm.io/gorm"
)

type UserPosition struct {
	gorm.Model
	EmployeeId uint   `json:"employee_id"`
	PositionId uint   `json:"position_id"`
	DateMove   string `json:"date_move" gorm:"DATETIME"`

	Employee employee.Employee `gorm:"foreignKey:EmployeeId;references:ID" json:"Employee"`
	Position position.Position `gorm:"foreignKey:PositionId;references:ID" json:"Position"`
}
