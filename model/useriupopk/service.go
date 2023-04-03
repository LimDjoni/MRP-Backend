package useriupopk

import (
	"ajebackend/model/user"
)

type Service interface {
	CreateUserIupopk(userId int, iupopkId int) (UserIupopk, error)
	LoginUser(input user.LoginUserInput) (user.TokenUser, error)
	DeleteUserIupopk(userId int, iupopkId int) error
	FindUser(id uint, iupopkId int) (user.User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateUserIupopk(userId int, iupopkId int) (UserIupopk, error) {
	createdUserIupopk, createdUserIupopkErr := s.repository.CreateUserIupopk(userId, iupopkId)

	return createdUserIupopk, createdUserIupopkErr
}

func (s *service) LoginUser(input user.LoginUserInput) (user.TokenUser, error) {
	loginUser, loginUserErr := s.repository.LoginUser(input)

	return loginUser, loginUserErr
}

func (s *service) DeleteUserIupopk(userId int, iupopkId int) error {
	deleteUserIupopkErr := s.repository.DeleteUserIupopk(userId, iupopkId)

	return deleteUserIupopkErr
}

func (s *service) FindUser(id uint, iupopkId int) (user.User, error) {
	user, userErr := s.repository.FindUser(id, iupopkId)

	return user, userErr
}
