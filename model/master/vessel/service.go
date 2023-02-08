package vessel

type Service interface {
	CheckVessel(vesselName string) (bool, error)
	GetVessel() ([]Vessel, error)
	CreateVessel(vesselName string) (Vessel, error)
	DetailVessel(id int) (Vessel, error)
	UpdateVessel(vesselName string, id int) (Vessel, error)
	DeleteVessel(id int) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CheckVessel(vesselName string) (bool, error) {
	masterVessel, masterVesselErr := s.repository.CheckVessel(vesselName)

	return masterVessel, masterVesselErr
}
func (s *service) GetVessel() ([]Vessel, error) {
	masterVessel, masterVesselErr := s.repository.GetVessel()

	return masterVessel, masterVesselErr
}
func (s *service) CreateVessel(vesselName string) (Vessel, error) {
	masterVessel, masterVesselErr := s.repository.CreateVessel(vesselName)

	return masterVessel, masterVesselErr
}
func (s *service) DetailVessel(id int) (Vessel, error) {
	masterVessel, masterVesselErr := s.repository.DetailVessel(id)

	return masterVessel, masterVesselErr
}
func (s *service) UpdateVessel(vesselName string, id int) (Vessel, error) {
	masterVessel, masterVesselErr := s.repository.UpdateVessel(vesselName, id)

	return masterVessel, masterVesselErr
}
func (s *service) DeleteVessel(id int) (bool, error) {
	masterVessel, masterVesselErr := s.repository.DeleteVessel(id)

	return masterVessel, masterVesselErr
}
