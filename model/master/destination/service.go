package destination

type Service interface {
	GetDestination() ([]Destination, error)
	GetDestinationByName(name string) (Destination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetDestination() ([]Destination, error) {
	getDestinations, getDestinationsErr := s.repository.GetDestination()

	return getDestinations, getDestinationsErr
}

func (s *service) GetDestinationByName(name string) (Destination, error) {
	getDestinationByName, getDestinationByNameErr := s.repository.GetDestinationByName(name)

	return getDestinationByName, getDestinationByNameErr
}
