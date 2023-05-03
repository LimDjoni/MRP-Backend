package cafassignment

type Service interface {
	ListCafAssignment(page int, sortFilter SortFilterCafAssignment, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListCafAssignment(page int, sortFilter SortFilterCafAssignment, iupopkId int) (Pagination, error) {
	listCafAssignment, listCafAssignmentErr := s.repository.ListCafAssignment(page, sortFilter, iupopkId)

	return listCafAssignment, listCafAssignmentErr
}
