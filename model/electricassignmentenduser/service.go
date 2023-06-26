package electricassignmentenduser

type Service interface {
	DetailElectricAssignment(id int, iupopkId int) (DetailElectricAssignment, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DetailElectricAssignment(id int, iupopkId int) (DetailElectricAssignment, error) {
	detailElectricAssignment, detailElectricAssignmentErr := s.repository.DetailElectricAssignment(id, iupopkId)

	return detailElectricAssignment, detailElectricAssignmentErr
}
