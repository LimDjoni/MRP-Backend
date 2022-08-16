package history

import "ajebackend/model/transaction"

type Service interface {
	CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint) (bool, error)
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

func(s *service) DeleteTransaction(id int, userId uint) (bool, error) {
	isDeletedTransaction, isDeletedTransactionErr := s.repository.DeleteTransaction(id, userId)

	return isDeletedTransaction, isDeletedTransactionErr
}
