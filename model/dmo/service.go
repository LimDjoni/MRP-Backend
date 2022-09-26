package dmo

type Service interface {
	GetReportDmoWithPeriod(period string) (Dmo, error)
	GetListReportDmoAll(page int) (Pagination, error)
	GetDataDmo(id int) (Dmo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportDmoWithPeriod(period string) (Dmo, error) {
	reportDmo, reportDmoErr := s.repository.GetReportDmoWithPeriod(period)

	return reportDmo, reportDmoErr
}

func (s *service) GetListReportDmoAll(page int) (Pagination, error) {
	listReportDmo, listReportDmoErr := s.repository.GetListReportDmoAll(page)

	return listReportDmo, listReportDmoErr
}

func (s *service) GetDataDmo(id int) (Dmo, error) {
	dataDMo, dataDMoErr := s.repository.GetDataDmo(id)

	return dataDMo, dataDMoErr
}
