package reportdmo

type Service interface {
	GetReportDmoWithPeriod(period string) (ReportDmo, error)
	GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportDmoWithPeriod(period string) (ReportDmo, error) {
	reportWithPeriod, reportWithPeriodErr := s.repository.GetReportDmoWithPeriod(period)

	return reportWithPeriod, reportWithPeriodErr
}

func (s *service) GetListReportDmoAll(page int, filterReportDmo FilterAndSortReportDmo) (Pagination, error) {
	listReportDmo, listReportDmoErr := s.repository.GetListReportDmoAll(page, filterReportDmo)

	return listReportDmo, listReportDmoErr
}
