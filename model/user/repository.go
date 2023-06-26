package user

import (
	"ajebackend/helper"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user RegisterUserInput) (User, error)
	FindUser(id uint) (User, error)
	ChangePassword(newPassword string, id uint) (User, error)
	ResetPassword(email string, newPassword string) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) RegisterUser(user RegisterUserInput) (User, error) {
	var newUser User

	newUser.Username = strings.ToLower(user.Username)
	newUser.Password = user.Password
	newUser.Email = strings.ToLower(user.Email)
	newUser.IsActive = true

	err := r.db.Create(&newUser).Error

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (r *repository) FindUser(id uint) (User, error) {
	var user User

	errFind := r.db.Where("id = ?", id).First(&user).Error

	return user, errFind
}

func (r *repository) ChangePassword(newPassword string, id uint) (User, error) {
	var user User

	errFind := r.db.Where("id = ?", id).Find(&user).Error

	if errFind != nil {
		return user, errFind
	}

	hashPassword, err := helper.GeneratePasswordHash(newPassword)

	if err != nil {
		return user, err
	}

	errUpd := r.db.Model(&user).Update("password", hashPassword).Error

	if errUpd != nil {
		return user, errUpd
	}

	return user, nil
}

func (r *repository) ResetPassword(email string, newPassword string) (User, error) {
	var user User

	errFind := r.db.Where("email = ?", email).Find(&user).Error

	if errFind != nil {
		return user, errFind
	}

	hashPassword, err := helper.GeneratePasswordHash(newPassword)

	if err != nil {
		return user, err
	}

	errUpd := r.db.Model(&user).Update("password", hashPassword).Error

	if errUpd != nil {
		return user, errUpd
	}

	return user, nil
}
