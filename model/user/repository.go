package user

import (
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user RegisterUserInput) (User, error)
	FindUser(id uint) (User, error)
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
