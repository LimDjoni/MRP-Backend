package minerba

type Service interface {
	GetReportMinerbaWithPeriod(period string, iupopkId int) (Minerba, error)
	GetListReportMinerbaAll(page int, filterMinerba FilterAndSortMinerba, iupopkId int) (Pagination, error)
	GetDataMinerba(id int, iupopkId int) (Minerba, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportMinerbaWithPeriod(period string, iupopkId int) (Minerba, error) {
	reportMinerba, reportMinerbaErr := s.repository.GetReportMinerbaWithPeriod(period, iupopkId)

	return reportMinerba, reportMinerbaErr
}

func (s *service) GetListReportMinerbaAll(page int, filterMinerba FilterAndSortMinerba, iupopkId int) (Pagination, error) {
	listReportMinerba, listReportMinerbaErr := s.repository.GetListReportMinerbaAll(page, filterMinerba, iupopkId)

	return listReportMinerba, listReportMinerbaErr
}

func (s *service) GetDataMinerba(id int, iupopkId int) (Minerba, error) {
	dataMinerba, dataMinerbaErr := s.repository.GetDataMinerba(id, iupopkId)

	return dataMinerba, dataMinerbaErr
}
