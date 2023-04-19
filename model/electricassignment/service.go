package electricassignment

type Service interface {
	ListElectricAssignment(page int, sortFilter SortFilterElectricAssignment, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListElectricAssignment(page int, sortFilter SortFilterElectricAssignment, iupopkId int) (Pagination, error) {
	listElectricAssignment, listElectricAssignmentErr := s.repository.ListElectricAssignment(page, sortFilter, iupopkId)

	return listElectricAssignment, listElectricAssignmentErr
}
