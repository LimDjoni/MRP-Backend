package history

import (
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"
)

type Service interface {
	CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransactionDN(id int, userId uint) (bool, error)
	UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocumentTransactionDN (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error)
	CreateMinerba (period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error)
	DeleteMinerba (idMinerba int, userId uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func(s *service) CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	transaction, transactionErr := s.repository.CreateTransactionDN(inputTransactionDN, userId)

	return transaction, transactionErr
}

func(s *service) DeleteTransactionDN(id int, userId uint) (bool, error) {
	isDeletedTransaction, isDeletedTransactionErr := s.repository.DeleteTransactionDN(id, userId)

	return isDeletedTransaction, isDeletedTransactionErr
}

func (s *service) UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	updateTransaction, updateTransactionErr := s.repository.UpdateTransactionDN(idTransaction, inputEditTransactionDN, userId)

	return updateTransaction, updateTransactionErr
}

func (s *service) UploadDocumentTransactionDN (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error) {
	uploadedDocument, uploadedDocumentErr := s.repository.UploadDocumentTransactionDN(idTransaction, urlS3, userId, documentType)

	return uploadedDocument, uploadedDocumentErr
}

func (s *service) CreateMinerba (period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error) {
	createdMinerba, createdMinerbaErr := s.repository.CreateMinerba(period, baseIdNumber, updateTransaction, userId)

	return createdMinerba, createdMinerbaErr
}

func (s *service) DeleteMinerba (idMinerba int, userId uint) (bool, error) {
	isDeletedMinerba, isDeletedMinerbaErr := s.repository.DeleteMinerba(idMinerba, userId)

	return isDeletedMinerba, isDeletedMinerbaErr
}
