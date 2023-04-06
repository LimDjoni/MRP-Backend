package allmaster

import "ajebackend/model/master/iupopk"

type Service interface {
	ListMasterData() (MasterData, error)
	FindIupopk(iupopkId int) (iupopk.Iupopk, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListMasterData() (MasterData, error) {
	listMasterData, listMasterDataErr := s.repository.ListMasterData()

	return listMasterData, listMasterDataErr
}

func (s *service) FindIupopk(iupopkId int) (iupopk.Iupopk, error) {
	findIupopk, findIupopkErr := s.repository.FindIupopk(iupopkId)

	return findIupopk, findIupopkErr
}
