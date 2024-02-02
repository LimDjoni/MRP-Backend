package contract

type Service interface {
	GetListReportContractAll(page int, filterContract FilterAndSortContract, iupopkId int) (Pagination, error)
	GetDataContract(id int, iupopkId int) (Contract, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetListReportContractAll(page int, filterContract FilterAndSortContract, iupopkId int) (Pagination, error) {
	listReportContract, listReportContractErr := s.repository.GetListReportContractAll(page, filterContract, iupopkId)

	return listReportContract, listReportContractErr
}

func (s *service) GetDataContract(id int, iupopkId int) (Contract, error) {
	dataContract, dataContractErr := s.repository.GetDataContract(id, iupopkId)

	return dataContract, dataContractErr
}
