package coareport

import "ajebackend/model/transaction"

type Service interface {
	GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	listTransaction, listTransactionErr := s.repository.GetTransactionCoaReport(dateFrom, dateTo, iupopkId)

	return listTransaction, listTransactionErr
}
