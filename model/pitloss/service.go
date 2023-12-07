package pitloss

type Service interface {
	DetailJettyBalance(id int, iupopkId int) (OutputJettyBalancePitLossDetail, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DetailJettyBalance(id int, iupopkId int) (OutputJettyBalancePitLossDetail, error) {
	data, dataErr := s.repository.DetailJettyBalance(id, iupopkId)

	return data, dataErr
}
