package user

import (
	"ajebackend/helper"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	RegisterUser(user RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (TokenUser, error)
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

func (r *repository) LoginUser(input LoginUserInput) (TokenUser, error) {
	var tokenUser TokenUser
	var isValidPassword bool
	var username User
	var email User

	dataUser := map[string]interface{}{
		"id":       0,
		"username": "",
		"email":    "",
		"role":     "",
	}
	usernameErr := r.db.Where("username = ?", input.Data).First(&username).Error

	emailErr := r.db.Where("email = ?", strings.ToLower(input.Data)).First(&email).Error

	if emailErr != nil && usernameErr != nil {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	if usernameErr != nil {
		dataUser["username"] = email.Username
		dataUser["email"] = email.Email
		dataUser["id"] = email.ID
		dataUser["role"] = email.Role
		isValidPassword = helper.CheckPassword(input.Password, email.Password)
	}

	if emailErr != nil {
		dataUser["username"] = username.Username
		dataUser["email"] = username.Email
		dataUser["id"] = username.ID
		dataUser["role"] = username.Role
		isValidPassword = helper.CheckPassword(input.Password, username.Password)
	}

	if username.IsActive == false && email.IsActive == false {

		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	if !isValidPassword {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	token, tokenErr := helper.GenerateToken(dataUser["id"].(uint), dataUser["username"].(string), dataUser["email"].(string))
	tokenUser.Token = token
	tokenUser.Username = dataUser["username"].(string)
	tokenUser.Email = dataUser["email"].(string)
	tokenUser.Role = dataUser["role"].(string)
	if tokenErr != nil {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	return tokenUser, nil
}

func (r *repository) FindUser(id uint) (User, error) {
	var user User

	errFind := r.db.Where("id = ? ", id).First(&user).Error

	return user, errFind
}
