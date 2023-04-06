package rkab

type Service interface {
	ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error)
	DetailRkabWithYear(year int, iupopkId int) ([]Rkab, error)
	DetailRkabWithId(id int, iupopkId int) (Rkab, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error) {
	listRkab, listRkabErr := s.repository.ListRkab(page, sortFilter, iupopkId)

	return listRkab, listRkabErr
}

func (s *service) DetailRkabWithYear(year int, iupopkId int) ([]Rkab, error) {
	detailRkab, detailRkabErr := s.repository.DetailRkabWithYear(year, iupopkId)

	return detailRkab, detailRkabErr
}

func (s *service) DetailRkabWithId(id int, iupopkId int) (Rkab, error) {
	detailRkab, detailRkabErr := s.repository.DetailRkabWithId(id, iupopkId)

	return detailRkab, detailRkabErr
}
