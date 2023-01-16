package groupingvesseldn

type Service interface {
	ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn) (Pagination, error) {
	listGroupingVesselDn, listGroupingVesselDnErr := s.repository.ListGroupingVesselDn(page, sortFilter)

	return listGroupingVesselDn, listGroupingVesselDnErr
}
