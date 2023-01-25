package dmovessel

import "ajebackend/model/groupingvesseldn"

type Service interface {
	GetDataDmoVessel(id uint) ([]DmoVessel, error)
	ListGroupingVesselWithoutDmo() ([]groupingvesseldn.GroupingVesselDn, error)
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

func (s *service) ListGroupingVesselWithoutDmo() ([]groupingvesseldn.GroupingVesselDn, error) {
	listGroupingVesselWithoutDmo, listGroupingVesselWithoutDmoErr := s.repository.ListGroupingVesselWithoutDmo()

	return listGroupingVesselWithoutDmo, listGroupingVesselWithoutDmoErr
}
