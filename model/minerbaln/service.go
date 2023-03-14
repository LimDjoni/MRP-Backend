package minerbaln

type Service interface {
	GetReportMinerbaLnWithPeriod(period string, iupopkId int) (MinerbaLn, error)
	GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn, iupopkId int) (Pagination, error)
	GetDataMinerbaLn(id int, iupopkId int) (MinerbaLn, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportMinerbaLnWithPeriod(period string, iupopkId int) (MinerbaLn, error) {
	reportMinerbaLn, reportMinerbaLnErr := s.repository.GetReportMinerbaLnWithPeriod(period, iupopkId)

	return reportMinerbaLn, reportMinerbaLnErr
}

func (s *service) GetListReportMinerbaLnAll(page int, filterMinerbaLn FilterAndSortMinerbaLn, iupopkId int) (Pagination, error) {
	listReportMinerbaLn, listReportMinerbaLnErr := s.repository.GetListReportMinerbaLnAll(page, filterMinerbaLn, iupopkId)

	return listReportMinerbaLn, listReportMinerbaLnErr
}

func (s *service) GetDataMinerbaLn(id int, iupopkId int) (MinerbaLn, error) {
	dataMinerbaLn, dataMinerbaLnErr := s.repository.GetDataMinerbaLn(id, iupopkId)

	return dataMinerbaLn, dataMinerbaLnErr
}
