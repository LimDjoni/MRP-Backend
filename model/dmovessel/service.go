package dmovessel

type Service interface {
	GetDataDmoVessel(id uint) ([]DmoVessel, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetDataDmoVessel(id uint) ([]DmoVessel, error) {
	listDmoVessel, listDmoVesselErr := s.repository.GetDataDmoVessel(id)

	return listDmoVessel, listDmoVesselErr
}
