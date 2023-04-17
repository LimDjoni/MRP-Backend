package cafassignmentenduser

type Service interface {
	DetailCafAssignment(id int, iupopkId int) (DetailCafAssignment, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DetailCafAssignment(id int, iupopkId int) (DetailCafAssignment, error) {
	detailCafAssignment, detailCafAssignmentErr := s.repository.DetailCafAssignment(id, iupopkId)

	return detailCafAssignment, detailCafAssignmentErr
}
