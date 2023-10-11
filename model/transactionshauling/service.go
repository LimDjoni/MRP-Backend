package transactionshauling

import (
	"ajebackend/model/transactionshauling/transactionispjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"
)

type Service interface {
	ListStockRom(page int, iupopkId int) (Pagination, error)
	ListTransactionHauling(page int, iupopkId int) (Pagination, error)
	DetailStockRom(iupopkId int, stockRomId int) (transactiontoisp.TransactionToIsp, error)
	DetailTransactionHauling(iupopkId int, transactionHaulingId int) (transactionispjetty.TransactionIspJetty, error)
	SummaryJettyTransactionPerDay(iupopkId int) (SummaryJettyTransactionPerDay, error)
	SummaryInventoryStockRom(iupopkId int) ([]InventoryStockRom, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListStockRom(page int, iupopkId int) (Pagination, error) {
	list, listErr := s.repository.ListStockRom(page, iupopkId)

	return list, listErr
}

func (s *service) ListTransactionHauling(page int, iupopkId int) (Pagination, error) {
	list, listErr := s.repository.ListTransactionHauling(page, iupopkId)

	return list, listErr
}

func (s *service) DetailStockRom(iupopkId int, stockRomId int) (transactiontoisp.TransactionToIsp, error) {
	detail, detailErr := s.repository.DetailStockRom(iupopkId, stockRomId)

	return detail, detailErr
}

func (s *service) DetailTransactionHauling(iupopkId int, transactionHaulingId int) (transactionispjetty.TransactionIspJetty, error) {
	detail, detailErr := s.repository.DetailTransactionHauling(iupopkId, transactionHaulingId)

	return detail, detailErr
}

func (s *service) SummaryJettyTransactionPerDay(iupopkId int) (SummaryJettyTransactionPerDay, error) {
	summary, summaryErr := s.repository.SummaryJettyTransactionPerDay(iupopkId)

	return summary, summaryErr
}

func (s *service) SummaryInventoryStockRom(iupopkId int) ([]InventoryStockRom, error) {
	summary, summaryErr := s.repository.SummaryInventoryStockRom(iupopkId)

	return summary, summaryErr
}
