package trader

type Service interface {
	ListTrader() ([]Trader, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListTrader() ([]Trader, error) {
	listTrader, listTraderErr := s.repository.ListTrader()

	return listTrader, listTraderErr
}
