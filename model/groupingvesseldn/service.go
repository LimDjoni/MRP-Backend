package groupingvesseldn

type Service interface {
	ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListGroupingVesselDn(page int, sortFilter SortFilterGroupingVesselDn, iupopkId int) (Pagination, error) {
	listGroupingVesselDn, listGroupingVesselDnErr := s.repository.ListGroupingVesselDn(page, sortFilter, iupopkId)

	return listGroupingVesselDn, listGroupingVesselDnErr
}
