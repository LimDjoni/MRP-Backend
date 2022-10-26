package traderdmo

import (
	"ajebackend/model/trader"
)

type Service interface {
	DmoIdListWithTraderId(idTrader int) ([]TraderDmo, error)
	TraderListWithDmoId(idDmo int) ([]trader.Trader, trader.Trader, error)
	GetTraderEndUserDmo(idDmo int) (trader.Trader, error)
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

func (s *service) TraderListWithDmoId(idDmo int) ([]trader.Trader, trader.Trader, error) {
	listTrader, endUser, dmoListWithDmoIdErr := s.repository.TraderListWithDmoId(idDmo)

	return listTrader, endUser, dmoListWithDmoIdErr
}

func (s *service) GetTraderEndUserDmo(idDmo int) (trader.Trader, error) {
	traderEndUserDmo, traderEndUserDmoErr := s.repository.GetTraderEndUserDmo(idDmo)

	return traderEndUserDmo, traderEndUserDmoErr
}
