package user

import (
	"errors"
	"mrpbackend/helper"
	"mrpbackend/model/master/userrole"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	RegisterUser(user RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (TokenUser, error)
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

	newUser.EmployeeId = user.EmployeeId
	newUser.Username = strings.ToLower(user.Username)
	newUser.Password = user.Password
	newUser.Email = strings.ToLower(user.Email)
	newUser.IsActive = true
	newUser.CodeEmp = user.CodeEmp

	err := r.db.Create(&newUser).Error

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (r *repository) LoginUser(input LoginUserInput) (TokenUser, error) {
	var tokenUser TokenUser
	var user User
	var listUserRole []string

	dataUser := map[string]interface{}{
		"id":          0,
		"employee_id": 0,
		"username":    "",
		"email":       "",
	}

	// Check email vs username
	isEmail := strings.Contains(input.Data, "@")

	var err error
	if isEmail {
		err = r.db.Preload("Employee").Preload("Employee.Department").
			Where("email = ?", strings.ToLower(input.Data)).
			First(&user).Error
	} else {
		err = r.db.Preload("Employee").Preload("Employee.Department").
			Where("username = ?", input.Data).
			First(&user).Error
	}

	if err != nil {
		return tokenUser, errors.New("User tidak ditemukan")
	}

	dataUser["username"] = user.Username
	dataUser["email"] = user.Email
	dataUser["id"] = user.ID
	dataUser["employee_id"] = user.EmployeeId

	// Get roles
	var userRole []userrole.UserRole
	r.db.Preload(clause.Associations).Where("user_id = ?", user.ID).Find(&userRole)

	for _, v := range userRole {
		listUserRole = append(listUserRole, v.Role.Name)
	}

	if !user.IsActive {
		return tokenUser, errors.New("User Tidak Aktif")
	}

	isValidPassword := helper.CheckPassword(input.Password, user.Password)
	if !isValidPassword {
		return tokenUser, errors.New("Password Salah")
	}

	token, tokenErr := helper.GenerateToken(dataUser["id"].(uint), dataUser["username"].(string), dataUser["email"].(string))
	if tokenErr != nil {
		return tokenUser, errors.New("invalid Email / Username / Password")
	}

	tokenUser = TokenUser{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Token:      token,
		EmployeeId: user.EmployeeId,
		Role:       listUserRole,
		Employee:   user.Employee, // ✅ Assign preloaded employee
		CodeEmp:    user.CodeEmp,
	}

	return tokenUser, nil
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
