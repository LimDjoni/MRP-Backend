package ici

type Service interface {
	GetAllIci() ([]Ici, error)
	CreateIci(inputIci InputCreateUpdateIci) (Ici, error)
	UpdateIci(inputIci InputCreateUpdateIci, id int) (Ici, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetAllIci() ([]Ici, error) {
	listIci, listIciErr := s.repository.GetAllIci()
	return listIci, listIciErr
}

func (s *service) CreateIci(inputIci InputCreateUpdateIci) (Ici, error) {
	createIci, createIciErr := s.repository.CreateIci(inputIci)

	return createIci, createIciErr
}

func (s *service) UpdateIci(inputIci InputCreateUpdateIci, id int) (Ici, error) {
	updateIci, updateIciErr := s.repository.UpdateIci(inputIci, id)

	return updateIci, updateIciErr
}
