package dmovessel

import "ajebackend/model/groupingvesseldn"

type Service interface {
	GetDataDmoVessel(id uint, iupopkId int) ([]DmoVessel, error)
	ListGroupingVesselWithoutDmo(iupopkId int) ([]groupingvesseldn.GroupingVesselDn, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetDataDmoVessel(id uint, iupopkId int) ([]DmoVessel, error) {
	listDmoVessel, listDmoVesselErr := s.repository.GetDataDmoVessel(id, iupopkId)

	return listDmoVessel, listDmoVesselErr
}

func (s *service) ListGroupingVesselWithoutDmo(iupopkId int) ([]groupingvesseldn.GroupingVesselDn, error) {
	listGroupingVesselWithoutDmo, listGroupingVesselWithoutDmoErr := s.repository.ListGroupingVesselWithoutDmo(iupopkId)

	return listGroupingVesselWithoutDmo, listGroupingVesselWithoutDmoErr
}
