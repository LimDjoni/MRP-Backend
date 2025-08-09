package adjuststock

type Service interface {
	CreateAdjustStock(adjuststock RegisterAdjustStockInput) (AdjustStock, error)
	FindAdjustStock() ([]AdjustStock, error)
	FindAdjustStockById(id uint) (AdjustStock, error)
	ListAdjustStock(page int, sortFilter SortFilterAdjustStock) (Pagination, error)
	UpdateAdjustStock(inputAdjustStock RegisterAdjustStockInput, id int) (AdjustStock, error)
	DeleteAdjustStock(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateAdjustStock(adjuststock RegisterAdjustStockInput) (AdjustStock, error) {
	newAdjustStock, err := s.repository.CreateAdjustStock(adjuststock)

	return newAdjustStock, err
}

func (s *service) FindAdjustStock() ([]AdjustStock, error) {
	adjuststock, err := s.repository.FindAdjustStock()

	return adjuststock, err
}

func (s *service) FindAdjustStockById(id uint) (AdjustStock, error) {
	adjustStock, err := s.repository.FindAdjustStockById(id)

	return adjustStock, err
}

func (s *service) ListAdjustStock(page int, sortFilter SortFilterAdjustStock) (Pagination, error) {
	adjustStock, err := s.repository.ListAdjustStock(page, sortFilter)

	return adjustStock, err
}

func (s *service) UpdateAdjustStock(inputAdjustStock RegisterAdjustStockInput, id int) (AdjustStock, error) {
	updateAdjustStock, updateAdjustStockErr := s.repository.UpdateAdjustStock(inputAdjustStock, id)

	return updateAdjustStock, updateAdjustStockErr
}

func (s *service) DeleteAdjustStock(id uint) (bool, error) {
	isDeletedAdjustStock, isDeletedAdjustStockErr := s.repository.DeleteAdjustStock(id)

	return isDeletedAdjustStock, isDeletedAdjustStockErr
}
