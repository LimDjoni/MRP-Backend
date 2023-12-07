package userrole

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	Name string `json:"name"`
}
