package production

type Service interface {
	GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error)
	DetailProduction(id int, iupopkId int) (Production, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error) {
	getListProduction, getListProductionErr := s.repository.GetListProduction(page, filter, iupopkId)

	return getListProduction, getListProductionErr
}

func (s *service) DetailProduction(id int, iupopkId int) (Production, error) {
	detailProduction, detailProductionErr := s.repository.DetailProduction(id, iupopkId)

	return detailProduction, detailProductionErr
}
