package user

type Service interface {
	RegisterUser(user RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (TokenUser, error)
	FindUser(id uint) (User, error)
	ChangePassword(newPassword string, id uint) (User, error)
	ResetPassword(email string, newPassword string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(user RegisterUserInput) (User, error) {
	newUser, err := s.repository.RegisterUser(user)

	return newUser, err
}

func (s *service) LoginUser(input LoginUserInput) (TokenUser, error) {
	loginUser, loginUserErr := s.repository.LoginUser(input)

	return loginUser, loginUserErr
}

func (s *service) FindUser(id uint) (User, error) {
	user, userErr := s.repository.FindUser(id)

	return user, userErr
}

func (s *service) ChangePassword(newPassword string, id uint) (User, error) {
	user, userErr := s.repository.ChangePassword(newPassword, id)

	return user, userErr
}

func (s *service) ResetPassword(email string, newPassword string) (User, error) {
	user, userErr := s.repository.ResetPassword(email, newPassword)

	return user, userErr
}
