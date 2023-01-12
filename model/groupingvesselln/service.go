package groupingvesselln

type Service interface {
	ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error) {
	listGroupingVesselLn, listGroupingVesselLnErr := s.repository.ListGroupingVesselLn(page, sortFilter)

	return listGroupingVesselLn, listGroupingVesselLnErr
}
