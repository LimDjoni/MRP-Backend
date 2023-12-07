package user

import (
	"ajebackend/model/master/userrole"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string             `json:"username" gorm:"UNIQUE"`
	Password   string             `json:"password"`
	Email      string             `json:"email" gorm:"UNIQUE"`
	UserRoleId *uint              `json:"user_role_id"`
	UserRole   *userrole.UserRole `json:"user_role"`
	IsActive   bool               `json:"is_active"`
}
