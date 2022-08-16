package user

type Service interface {
	RegisterUser(user RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (TokenUser, error)
	FindUser(id uint) (User, error)
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
	token, err := s.repository.LoginUser(input)

	return token, err
}

func (s *service) FindUser(id uint) (User, error) {
	user, userErr := s.repository.FindUser(id)

	return user, userErr
}
