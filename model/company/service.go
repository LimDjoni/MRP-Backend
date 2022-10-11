package company

type Service interface {
	ListCompany() ([]Company, error)
	CreateCompany(inputCompany InputCreateUpdateCompany) (Company, error)
	UpdateCompany(inputCompany InputCreateUpdateCompany, id int) (Company, error)
	DeleteCompany(id int) (bool, error)
	DetailCompany(id int) (Company, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListCompany() ([]Company, error) {
	listCompany, listCompanyErr := s.repository.ListCompany()

	return listCompany, listCompanyErr
}

func (s *service) CreateCompany(inputCompany InputCreateUpdateCompany) (Company, error) {
	createCompany, createCompanyErr := s.repository.CreateCompany(inputCompany)

	return createCompany, createCompanyErr
}

func (s *service) UpdateCompany(inputCompany InputCreateUpdateCompany, id int) (Company, error) {
	updateCompany, updateCompanyErr := s.repository.UpdateCompany(inputCompany, id)

	return updateCompany, updateCompanyErr
}

func (s *service) DeleteCompany(id int) (bool, error) {
	deleteCompany, deleteCompanyErr := s.repository.DeleteCompany(id)

	return deleteCompany, deleteCompanyErr
}

func (s *service) DetailCompany(id int) (Company, error) {
	detailCompany, detailCompanyErr := s.repository.DetailCompany(id)

	return detailCompany, detailCompanyErr
}
