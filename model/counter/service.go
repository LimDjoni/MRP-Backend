package counter

import "ajebackend/model/master/iupopk"

type Service interface {
	UpdateCounter() error
	CreateIupopk(input iupopk.InputIupopk) (iupopk.Iupopk, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) UpdateCounter() error {
	updateCounterErr := s.repository.UpdateCounter()

	return updateCounterErr
}

func (s *service) CreateIupopk(input iupopk.InputIupopk) (iupopk.Iupopk, error) {
	createdIupopk, createdIupopkErr := s.repository.CreateIupopk(input)

	return createdIupopk, createdIupopkErr
}
