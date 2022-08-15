package transaction

type Service interface {
	CreateTransactionDN (inputTransactionDN DataTransactionInput) (Transaction, error)
	ListDataDN (page int) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	DeleteTransaction(id int) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateTransactionDN (inputTransactionDN DataTransactionInput) (Transaction, error) {
	transaction, transactionErr := s.repository.CreateTransactionDN(inputTransactionDN)

	return transaction, transactionErr
}

func (s *service) ListDataDN(page int) (Pagination, error) {
	listDN, listDNErr := s.repository.ListDataDN(page)

	return listDN, listDNErr
}

func (s *service) DetailTransactionDN(id int) (Transaction, error) {
	detailTransactionDN, detailTransactionDNErr := s.repository.DetailTransactionDN(id)

	return detailTransactionDN, detailTransactionDNErr
}

func (s *service) DeleteTransaction(id int) (bool, error) {
	deleteTransaction, deleteTransactionErr := s.repository.DeleteTransaction(id)

	return deleteTransaction, deleteTransactionErr
}