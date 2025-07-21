package departmentform

import (
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/role"

	"gorm.io/gorm"
)

type DepartmentForm struct {
	gorm.Model
	DepartmentId uint `json:"department_id"`
	RoleId       uint `json:"role_id"`

	Department department.Department `gorm:"foreignKey:DepartmentId" json:"Department"` // Add this line
	Role       role.Role             `gorm:"foreignKey:RoleId" json:"role"`             // Add this line
}
