package transactionrequestreport

import (
	"ajebackend/model/masterreport"
)

type Service interface {
	CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error)
	DeleteTransactionReport() (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error) {
	transaction, transactionErr := s.repository.CreateTransactionReport(input, iupopkId)

	return transaction, transactionErr
}

func (s *service) UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error) {
	transaction, transactionErr := s.repository.UpdateTransactionReport(DnDocumentLink, LnDocumentLink, id, iupopkId)

	return transaction, transactionErr
}

func (s *service) DeleteTransactionReport() (bool, error) {
	isDeleted, isDeletedErr := s.repository.DeleteTransactionReport()

	return isDeleted, isDeletedErr
}
