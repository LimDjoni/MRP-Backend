package transaction

type Service interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDNAndMinerba(listData []int)(bool, error)
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

func (s *service) ListDataDNWithoutMinerba() ([]Transaction, error) {
	listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr := s.repository.ListDataDNWithoutMinerba()

	return listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr
}

func (s *service) CheckDataDNAndMinerba(listData []int)(bool, error) {
	checkData, checkDataErr := s.repository.CheckDataDNAndMinerba(listData)

	return checkData, checkDataErr
}
