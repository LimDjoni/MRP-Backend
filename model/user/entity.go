package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"UNIQUE"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"UNIQUE"`
	IsActive bool   `json:"is_active"`
}
