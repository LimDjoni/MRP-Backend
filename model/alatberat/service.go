package alatberat

type Service interface {
	CreateAlatBerat(alatberat RegisterAlatBeratInput) (AlatBerat, error)
	FindAlatBerat() ([]AlatBerat, error)
	FindAlatBeratById(id uint) (AlatBerat, error)
	GetListAlatBerat(page int, sortFilter SortFilterAlatBerat) (Pagination, error)
	FindConsumption(brandId uint, heavyEquipmentId uint, seriesId uint) (AlatBerat, error)
	UpdateAlatBerat(inputAlatBerat RegisterAlatBeratInput, id int) (AlatBerat, error)
	DeleteAlatBerat(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateAlatBerat(alatberat RegisterAlatBeratInput) (AlatBerat, error) {
	newAlatBerat, err := s.repository.CreateAlatBerat(alatberat)

	return newAlatBerat, err
}

func (s *service) FindAlatBerat() ([]AlatBerat, error) {
	alatberat, err := s.repository.FindAlatBerat()

	return alatberat, err
}

func (s *service) FindAlatBeratById(id uint) (AlatBerat, error) {
	alatBerat, err := s.repository.FindAlatBeratById(id)

	return alatBerat, err
}

func (s *service) GetListAlatBerat(page int, sortFilter SortFilterAlatBerat) (Pagination, error) {
	listAlatBerat, listAlatBeratErr := s.repository.ListAlatBerat(page, sortFilter)

	return listAlatBerat, listAlatBeratErr
}

func (s *service) FindConsumption(brandId uint, heavyEquipmentId uint, seriesId uint) (AlatBerat, error) {
	listAlatBerat, listAlatBeratErr := s.repository.FindConsumption(brandId, heavyEquipmentId, seriesId)

	return listAlatBerat, listAlatBeratErr
}

func (s *service) UpdateAlatBerat(inputAlatBerat RegisterAlatBeratInput, id int) (AlatBerat, error) {
	updateAlatBerat, updateAlatBeratErr := s.repository.UpdateAlatBerat(inputAlatBerat, id)

	return updateAlatBerat, updateAlatBeratErr
}

func (s *service) DeleteAlatBerat(id uint) (bool, error) {
	isDeletedAlatBerat, isDeletedAlatBeratErr := s.repository.DeleteAlatBerat(id)

	return isDeletedAlatBerat, isDeletedAlatBeratErr
}
