package production

type Service interface {
	GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error)
	DetailProduction(id int, iupopkId int) (Production, error)
	SummaryProduction(year string, iupopkId int) (OutputSummaryProduction, error)
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

func (s *service) SummaryProduction(year string, iupopkId int) (OutputSummaryProduction, error) {
	summary, summaryErr := s.repository.SummaryProduction(year, iupopkId)

	return summary, summaryErr
}
