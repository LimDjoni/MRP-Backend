package groupingvesselln

type Service interface {
	ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error)
	ListGroupingVesselLnWithPeriod(month string, year int) ([]GroupingVesselLn, error)
	DetailInsw(id int) (DetailInsw, error)
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

func (s *service) ListGroupingVesselLnWithPeriod(month string, year int) ([]GroupingVesselLn, error) {
	listGroupingVesselLnWithoutInsw, listGroupingVesselLnWithoutInswErr := s.repository.ListGroupingVesselLnWithPeriod(month, year)

	return listGroupingVesselLnWithoutInsw, listGroupingVesselLnWithoutInswErr
}

func (s *service) DetailInsw(id int) (DetailInsw, error) {
	detailInsw, detailInswErr := s.repository.DetailInsw(id)

	return detailInsw, detailInswErr
}
