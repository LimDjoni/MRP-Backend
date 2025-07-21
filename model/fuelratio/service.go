package fuelratio

type Service interface {
	CreateFuelRatio(fuelratios RegisterFuelRatioInput) (FuelRatio, error)
	FindFuelRatio() ([]FuelRatio, error)
	FindFuelRatioById(id uint) (FuelRatio, error)
	GetListFuelRatio(page int, sortFilter SortFilterFuelRatio) (Pagination, error)
	FindFuelRatioExport(sortFilter SortFilterFuelRatioSummary) ([]SortFilterFuelRatioSummary, error)
	ListRangkuman(page int, sortFilter SortFilterFuelRatioSummary) (Pagination, error)
	UpdateFuelRatio(inputFuelRatio RegisterFuelRatioInput, id int) (FuelRatio, error)
	DeleteFuelRatio(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateFuelRatio(fuelratios RegisterFuelRatioInput) (FuelRatio, error) {
	newFuelRatio, err := s.repository.CreateFuelRatio(fuelratios)

	return newFuelRatio, err
}

func (s *service) FindFuelRatio() ([]FuelRatio, error) {
	fuelratios, err := s.repository.FindFuelRatio()

	return fuelratios, err
}

func (s *service) FindFuelRatioById(id uint) (FuelRatio, error) {
	alatBerat, err := s.repository.FindFuelRatioById(id)

	return alatBerat, err
}

func (s *service) GetListFuelRatio(page int, sortFilter SortFilterFuelRatio) (Pagination, error) {
	listFuelRatios, listFuelRatiosErr := s.repository.ListFuelRatio(page, sortFilter)

	return listFuelRatios, listFuelRatiosErr
}

func (s *service) FindFuelRatioExport(sortFilter SortFilterFuelRatioSummary) ([]SortFilterFuelRatioSummary, error) {
	listFuelRatios, listFuelRatiosErr := s.repository.FindFuelRatioExport(sortFilter)

	return listFuelRatios, listFuelRatiosErr
}

func (s *service) ListRangkuman(page int, sortFilter SortFilterFuelRatioSummary) (Pagination, error) {
	listFuelRatios, listFuelRatiosErr := s.repository.ListRangkuman(page, sortFilter)

	return listFuelRatios, listFuelRatiosErr
}

func (s *service) UpdateFuelRatio(inputFuelRatio RegisterFuelRatioInput, id int) (FuelRatio, error) {
	updateFuelRatio, updateFuelRatioErr := s.repository.UpdateFuelRatio(inputFuelRatio, id)

	return updateFuelRatio, updateFuelRatioErr
}

func (s *service) DeleteFuelRatio(id uint) (bool, error) {
	isDeletedFuelRatio, isDeletedFuelRatioErr := s.repository.DeleteFuelRatio(id)

	return isDeletedFuelRatio, isDeletedFuelRatioErr
}
