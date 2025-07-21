package user

import (
	"mrpbackend/model/employee"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EmployeeId uint   `json:"employee_id" gorm:"unique"` // note: unique should be lowercase
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
	Email      string `json:"email" gorm:"unique"`
	IsActive   bool   `json:"is_active"`
	CodeEmp    uint   `json:"code_emp"`

	Employee employee.Employee `gorm:"foreignKey:EmployeeId;references:ID" json:"employee"`
}
