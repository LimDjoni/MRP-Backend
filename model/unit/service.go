package unit

type Service interface {
	CreateUnit(units RegisterUnitInput) (Unit, error)
	FindUnit() ([]Unit, error)
	FindUnitById(id uint) (Unit, error)
	GetListUnit(page int, sortFilter SortFilterUnit) (Pagination, error)
	UpdateUnit(inputUnit RegisterUnitInput, id int) (Unit, error)
	DeleteUnit(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateUnit(units RegisterUnitInput) (Unit, error) {
	newUnit, err := s.repository.CreateUnit(units)

	return newUnit, err
}

func (s *service) FindUnit() ([]Unit, error) {
	units, err := s.repository.FindUnit()

	return units, err
}

func (s *service) FindUnitById(id uint) (Unit, error) {
	alatBerat, err := s.repository.FindUnitById(id)

	return alatBerat, err
}

func (s *service) GetListUnit(page int, sortFilter SortFilterUnit) (Pagination, error) {
	listReportMinerbaLn, listReportMinerbaLnErr := s.repository.ListUnit(page, sortFilter)

	return listReportMinerbaLn, listReportMinerbaLnErr
}

func (s *service) UpdateUnit(inputUnit RegisterUnitInput, id int) (Unit, error) {
	updateUnit, updateUnitErr := s.repository.UpdateUnit(inputUnit, id)

	return updateUnit, updateUnitErr
}

func (s *service) DeleteUnit(id uint) (bool, error) {
	isDeletedUnit, isDeletedUnitErr := s.repository.DeleteUnit(id)

	return isDeletedUnit, isDeletedUnitErr
}
