package reportdmo

type Service interface {
	GetReportDmoWithPeriod(period string, iupopkId int) (ReportDmo, error)
	GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportDmoWithPeriod(period string, iupopkId int) (ReportDmo, error) {
	reportWithPeriod, reportWithPeriodErr := s.repository.GetReportDmoWithPeriod(period, iupopkId)

	return reportWithPeriod, reportWithPeriodErr
}

func (s *service) GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo, iupopkId int) (Pagination, error) {
	listReportDmo, listReportDmoErr := s.repository.GetListReportDmoAll(page, filterReportDmo, iupopkId)

	return listReportDmo, listReportDmoErr
}
