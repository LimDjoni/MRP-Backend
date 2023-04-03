package insw

type Service interface {
	ListInsw(page int, sortFilter SortFilterInsw, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListInsw(page int, sortFilter SortFilterInsw, iupopkId int) (Pagination, error) {
	listInsw, listInswErr := s.repository.ListInsw(page, sortFilter, iupopkId)

	return listInsw, listInswErr
}
