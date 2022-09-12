package minerba

type Service interface {
	GetReportMinerbaWithPeriod(period string) (Minerba, error)
	GetListReportMinerbaAll(page int) (Pagination, error)
	GetDataMinerba(id int)(Minerba, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportMinerbaWithPeriod(period string) (Minerba, error) {
	reportMinerba, reportMinerbaErr := s.repository.GetReportMinerbaWithPeriod(period)

	return reportMinerba, reportMinerbaErr
}

func (s *service) GetListReportMinerbaAll(page int) (Pagination, error) {
	listReportMinerba, listReportMinerbaErr := s.repository.GetListReportMinerbaAll(page)

	return listReportMinerba, listReportMinerbaErr
}

func (s *service) GetDataMinerba(id int)(Minerba, error) {
	dataMinerba, dataMinerbaErr := s.repository.GetDataMinerba(id)

	return dataMinerba, dataMinerbaErr
}

