package dmo

type Service interface {
	GetReportDmoWithPeriod(period string, iupopkId int) (Dmo, error)
	GetListReportDmoAll(page int, filterDmo FilterAndSortDmo, iupopkId int) (Pagination, error)
	GetDataDmo(id int, iupopkId int) (Dmo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetReportDmoWithPeriod(period string, iupopkId int) (Dmo, error) {
	reportDmo, reportDmoErr := s.repository.GetReportDmoWithPeriod(period, iupopkId)

	return reportDmo, reportDmoErr
}

func (s *service) GetListReportDmoAll(page int, filterDmo FilterAndSortDmo, iupopkId int) (Pagination, error) {
	listReportDmo, listReportDmoErr := s.repository.GetListReportDmoAll(page, filterDmo, iupopkId)

	return listReportDmo, listReportDmoErr
}

func (s *service) GetDataDmo(id int, iupopkId int) (Dmo, error) {
	dataDMo, dataDMoErr := s.repository.GetDataDmo(id, iupopkId)

	return dataDMo, dataDMoErr
}
