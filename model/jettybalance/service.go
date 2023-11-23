package jettybalance

type Service interface {
	ListJettyBalance(page int, sortFilter SortFilterJettyBalance, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListJettyBalance(page int, sortFilter SortFilterJettyBalance, iupopkId int) (Pagination, error) {
	data, dataErr := s.repository.ListJettyBalance(page, sortFilter, iupopkId)

	return data, dataErr
}
