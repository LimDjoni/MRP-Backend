package minerbaln

type Service interface {
	GetReportMinerbaLnWithPeriod(period string) (MinerbaLn, error)
	GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn) (Pagination, error)
	GetDataMinerbaLn(id int) (MinerbaLn, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportMinerbaLnWithPeriod(period string) (MinerbaLn, error) {
	reportMinerbaLn, reportMinerbaLnErr := s.repository.GetReportMinerbaLnWithPeriod(period)

	return reportMinerbaLn, reportMinerbaLnErr
}

func (s *service) GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn) (Pagination, error) {
	listReportMinerbaLn, listReportMinerbaLnErr := s.repository.GetListReportMinerbaLnAll(page, filterMinerbaLn)

	return listReportMinerbaLn, listReportMinerbaLnErr
}

func (s *service) GetDataMinerbaLn(id int) (MinerbaLn, error) {
	dataMinerbaLn, dataMinerbaLnErr := s.repository.GetDataMinerbaLn(id)

	return dataMinerbaLn, dataMinerbaLnErr
}
