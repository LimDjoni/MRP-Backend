package production

type Service interface {
	GetListProduction(page int, filter FilterListProduction) (Pagination, error)
	DetailProduction(id int) (Production, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func(s *service) GetListProduction(page int, filter FilterListProduction) (Pagination, error) {
	getListProduction, getListProductionErr := s.repository.GetListProduction(page, filter)

	return getListProduction, getListProductionErr
}

func(s *service) DetailProduction(id int) (Production, error) {
	detailProduction, detailProductionErr := s.repository.DetailProduction(id)

	return detailProduction, detailProductionErr
}
