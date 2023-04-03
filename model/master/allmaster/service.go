package allmaster

type Service interface {
	ListMasterData() (MasterData, error)
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
