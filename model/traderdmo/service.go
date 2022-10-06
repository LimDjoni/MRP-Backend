package traderdmo

type Service interface {
	DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error) {
	dmoListWithTraderId, dmoListWithTraderIdErr := s.repository.DmoIdListWithTraderId(idTrader)

	return dmoListWithTraderId, dmoListWithTraderIdErr
}
