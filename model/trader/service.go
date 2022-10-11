package trader

type Service interface {
	ListTrader() ([]Trader, error)
	CheckListTrader(list []int) ([]Trader, error)
	CheckEndUser(id int) (Trader, error)
	CreateTrader(inputTrader InputCreateUpdateTrader) (Trader, error)
	UpdateTrader(inputTrader InputCreateUpdateTrader, id int) (Trader, error)
	DeleteTrader(id int) (bool, error)
	ListTraderWithCompanyId(id int) ([]Trader, error)
	DetailTrader(id int) (Trader, error)
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

func (s *service) CheckListTrader(list []int) ([]Trader, error) {
	isListValid, isListValidErr := s.repository.CheckListTrader(list)

	return isListValid, isListValidErr
}

func (s *service) CheckEndUser(id int) (Trader, error) {
	isEndUserValid, isEndUserValidErr := s.repository.CheckEndUser(id)

	return isEndUserValid, isEndUserValidErr
}

func (s *service) CreateTrader(inputTrader InputCreateUpdateTrader) (Trader, error) {
	createTrader, createTraderErr := s.repository.CreateTrader(inputTrader)

	return createTrader, createTraderErr
}

func (s *service) UpdateTrader(inputTrader InputCreateUpdateTrader, id int) (Trader, error) {
	updateTrader, updateTraderErr := s.repository.UpdateTrader(inputTrader, id)

	return updateTrader, updateTraderErr
}

func (s *service) DeleteTrader(id int) (bool, error) {
	deleteTrader, deleteTraderErr := s.repository.DeleteTrader(id)

	return deleteTrader, deleteTraderErr
}

func (s *service) ListTraderWithCompanyId(id int) ([]Trader, error) {
	listTrader, listTraderErr := s.repository.ListTraderWithCompanyId(id)

	return listTrader, listTraderErr
}

func (s *service) DetailTrader(id int) (Trader, error) {
	detailTrader, detailTraderErr := s.repository.DetailTrader(id)

	return detailTrader, detailTraderErr
}
