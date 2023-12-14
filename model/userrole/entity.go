package userrole

import (
	"ajebackend/model/master/role"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	UserId uint      `json:"user_id"`
	User   user.User `json:"user"`
	RoleId uint      `json:"role_id"`
	Role   role.Role `json:"role"`
}
