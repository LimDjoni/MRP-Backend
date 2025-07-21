package userrole

import (
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/role"

	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	UserId       uint                  `json:"user_id"`
	RoleId       uint                  `json:"role_id"`
	Role         role.Role             `json:"Role"`
	DepartmentId uint                  `json:"department_id"`
	Department   department.Department `json:"Department"`
}
