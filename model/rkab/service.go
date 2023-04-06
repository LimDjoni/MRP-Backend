package rkab

type Service interface {
	ListRkab(page int, sortFilter SortFilterRkab, iupopkId int) (Pagination, error)
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
