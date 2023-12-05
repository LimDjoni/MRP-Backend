package useriupopk

import (
	"ajebackend/helper"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"
	"ajebackend/model/userrole"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateUserIupopk(userId int, iupopkId int) (UserIupopk, error)
	LoginUser(input user.LoginUserInput) (user.TokenUser, error)
	DeleteUserIupopk(userId int, iupopkId int) error
	FindUser(id uint, iupopkId int) (user.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUserIupopk(userId int, iupopkId int) (UserIupopk, error) {
	var findUser user.User
	var findIupopk iupopk.Iupopk
	var createdUserIupopk UserIupopk
	findUserErr := r.db.Where("id = ?", userId).First(&findUser).Error

	if findUserErr != nil {
		return createdUserIupopk, findUserErr
	}

	findIupopkErr := r.db.Where("id = ?", iupopkId).First(&findIupopk).Error

	if findIupopkErr != nil {
		return createdUserIupopk, findIupopkErr
	}

	var findUserIupopk UserIupopk

	findUserIupopkErr := r.db.Where("user_id = ? AND iupopk_id = ?", userId, iupopkId).First(&findUserIupopk).Error

	if findUserIupopkErr == nil {
		return createdUserIupopk, errors.New("User dengan Iupopk sudah terbuat")
	}

	createdUserIupopk.UserId = uint(userId)
	createdUserIupopk.IupopkId = uint(iupopkId)
	createUserIupopkErr := r.db.Create(&createdUserIupopk).Error

	return createdUserIupopk, createUserIupopkErr
}

func (r *repository) LoginUser(input user.LoginUserInput) (user.TokenUser, error) {
	var tokenUser user.TokenUser
	var isValidPassword bool
	var username user.User
	var email user.User

	var listUserRole []string

	dataUser := map[string]interface{}{
		"id":       0,
		"username": "",
		"email":    "",
	}
	usernameErr := r.db.Preload(clause.Associations).Where("username = ?", input.Data).First(&username).Error

	emailErr := r.db.Preload(clause.Associations).Where("email = ?", strings.ToLower(input.Data)).First(&email).Error

	if emailErr != nil && usernameErr != nil {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	if usernameErr != nil {
		dataUser["username"] = email.Username
		dataUser["email"] = email.Email
		dataUser["id"] = email.ID

		var userRole []userrole.UserRole

		r.db.Preload(clause.Associations).Where("user_id = ?", email.ID).Find(&userRole)

		for _, v := range userRole {
			listUserRole = append(listUserRole, v.Role.Name)
		}

		isValidPassword = helper.CheckPassword(input.Password, email.Password)
	}

	if emailErr != nil {
		dataUser["username"] = username.Username
		dataUser["email"] = username.Email
		dataUser["id"] = username.ID
		var userRole []userrole.UserRole

		r.db.Preload(clause.Associations).Where("user_id = ?", username.ID).Find(&userRole)

		for _, v := range userRole {
			listUserRole = append(listUserRole, v.Role.Name)
		}

		isValidPassword = helper.CheckPassword(input.Password, username.Password)
	}

	if username.IsActive == false && email.IsActive == false {

		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	if !isValidPassword {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	var iupopks []UserIupopk

	findUserIupopkErr := r.db.Preload(clause.Associations).Order("id asc").Where("user_id = ?", dataUser["id"]).Find(&iupopks).Error

	if findUserIupopkErr != nil {
		return tokenUser, errors.New("user not have iupopk")
	}

	var iupopk []iupopk.Iupopk

	for _, v := range iupopks {
		iupopk = append(iupopk, v.Iupopk)
	}

	token, tokenErr := helper.GenerateToken(dataUser["id"].(uint), dataUser["username"].(string), dataUser["email"].(string))
	tokenUser.Token = token
	tokenUser.Username = dataUser["username"].(string)
	tokenUser.Email = dataUser["email"].(string)

	tokenUser.Role = listUserRole
	tokenUser.Iupopk = iupopk

	if tokenErr != nil {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	return tokenUser, nil
}

func (r *repository) DeleteUserIupopk(userId int, iupopkId int) error {

	var findUserIupopk UserIupopk

	findUserIupopkErr := r.db.Where("user_id = ? AND iupopk_id = ?", userId, iupopkId).First(&findUserIupopk).Error

	if findUserIupopkErr != nil {
		return findUserIupopkErr
	}

	deleteUserIupopkErr := r.db.Unscoped().Where("id = ? ", findUserIupopk.ID).Delete(&findUserIupopk).Error

	if deleteUserIupopkErr != nil {
		return deleteUserIupopkErr
	}

	return nil
}

func (r *repository) FindUser(id uint, iupopkId int) (user.User, error) {
	var userIupopk UserIupopk

	errFind := r.db.Preload(clause.Associations).Where("user_id = ? AND iupopk_id = ?", id, iupopkId).First(&userIupopk).Error

	var user user.User

	user = userIupopk.User

	return user, errFind
}
