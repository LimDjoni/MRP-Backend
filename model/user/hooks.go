package user

import (
	"mrpbackend/helper"

	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	newPassword, err := helper.GeneratePasswordHash(u.Password)

	if err != nil {
		return err
	}

	u.Password = newPassword

	return
}
