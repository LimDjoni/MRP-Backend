package trader

type Service interface {
	ListTrader() ([]Trader, error)
	CheckListTrader(list []uint) (bool, error)
	CheckEndUser(id uint) (bool, error)
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

func (s *service) CheckListTrader(list []uint) (bool, error) {
	isListValid, isListValidErr := s.repository.CheckListTrader(list)

	return isListValid, isListValidErr
}

func (s *service) CheckEndUser(id uint) (bool, error) {
	isEndUserValid, isEndUserValidErr := s.repository.CheckEndUser(id)

	return isEndUserValid, isEndUserValidErr
}
