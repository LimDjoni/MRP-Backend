package fuelin

type Service interface {
	CreateFuelIn(fuelin RegisterFuelInInput) (FuelIn, error)
	FindFuelIn() ([]FuelIn, error)
	FindFuelInById(id uint) (FuelIn, error)
	GetListFuelIn(page int, sortFilter SortFilterFuelIn) (Pagination, error)
	UpdateFuelIn(inputFuelIn RegisterFuelInInput, id int) (FuelIn, error)
	DeleteFuelIn(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateFuelIn(fuelin RegisterFuelInInput) (FuelIn, error) {
	newFuelIn, err := s.repository.CreateFuelIn(fuelin)

	return newFuelIn, err
}

func (s *service) FindFuelIn() ([]FuelIn, error) {
	fuelin, err := s.repository.FindFuelIn()

	return fuelin, err
}

func (s *service) FindFuelInById(id uint) (FuelIn, error) {
	alatBerat, err := s.repository.FindFuelInById(id)

	return alatBerat, err
}

func (s *service) GetListFuelIn(page int, sortFilter SortFilterFuelIn) (Pagination, error) {
	listFuelIn, listFuelInErr := s.repository.ListFuelIn(page, sortFilter)

	return listFuelIn, listFuelInErr
}

func (s *service) UpdateFuelIn(inputFuelIn RegisterFuelInInput, id int) (FuelIn, error) {
	updateFuelIn, updateFuelInErr := s.repository.UpdateFuelIn(inputFuelIn, id)

	return updateFuelIn, updateFuelInErr
}

func (s *service) DeleteFuelIn(id uint) (bool, error) {
	isDeletedFuelIn, isDeletedFuelInErr := s.repository.DeleteFuelIn(id)

	return isDeletedFuelIn, isDeletedFuelInErr
}
