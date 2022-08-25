package transaction

type Service interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error) {
	listDN, listDNErr := s.repository.ListDataDN(page, sortFilter)

	return listDN, listDNErr
}

func (s *service) DetailTransactionDN(id int) (Transaction, error) {
	detailTransactionDN, detailTransactionDNErr := s.repository.DetailTransactionDN(id)

	return detailTransactionDN, detailTransactionDNErr
}
