package stockfuel

type Service interface {
	ListStockFuel(sortFilter StockFuelSummary) ([]StockFuelSummary, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListStockFuel(sortFilter StockFuelSummary) ([]StockFuelSummary, error) {
	listFuelRatios, listFuelRatiosErr := s.repository.ListStockFuel(sortFilter)

	return listFuelRatios, listFuelRatiosErr
}
