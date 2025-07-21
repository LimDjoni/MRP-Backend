package employee

type Service interface {
	CreateEmployee(employees RegisterEmployeeInput) (Employee, error)
	FindEmployee(empCode uint) ([]Employee, error)
	FindEmployeeById(id uint) (Employee, error)
	FindEmployeeByDepartmentId(userId uint) ([]Employee, error)
	GetListEmployee(page int, sortFilter SortFilterEmployee) (Pagination, error)
	UpdateEmployee(inputEmployee UpdateEmployeeInput, id int) (Employee, error)
	DeleteEmployee(id uint) (bool, error)
	ListDashboard(empCode uint, dashboardSort SortFilterDashboardEmployee) (DashboardEmployee, error)
	ListDashboardTurnover(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeTurnOver, error)
	ListDashboardKontrak(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeKontrak, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateEmployee(employees RegisterEmployeeInput) (Employee, error) {
	newEmployee, err := s.repository.CreateEmployee(employees)

	return newEmployee, err
}

func (s *service) FindEmployee(empCode uint) ([]Employee, error) {
	employees, err := s.repository.FindEmployee(empCode)

	return employees, err
}

func (s *service) FindEmployeeById(id uint) (Employee, error) {
	employees, err := s.repository.FindEmployeeById(id)

	return employees, err
}

func (s *service) FindEmployeeByDepartmentId(userId uint) ([]Employee, error) {
	employees, err := s.repository.FindEmployeeByDepartmentId(userId)

	return employees, err
}

func (s *service) GetListEmployee(page int, sortFilter SortFilterEmployee) (Pagination, error) {
	listListEmployee, listListEmployeeErr := s.repository.ListEmployee(page, sortFilter)

	return listListEmployee, listListEmployeeErr
}

func (s *service) UpdateEmployee(inputEmployee UpdateEmployeeInput, id int) (Employee, error) {
	updateEmployee, updateEmployeeErr := s.repository.UpdateEmployee(inputEmployee, id)

	return updateEmployee, updateEmployeeErr
}

func (s *service) DeleteEmployee(id uint) (bool, error) {
	isDeletedEmployee, isDeletedEmployeeErr := s.repository.DeleteEmployee(id)

	return isDeletedEmployee, isDeletedEmployeeErr
}

func (s *service) ListDashboard(empCode uint, dashboardSort SortFilterDashboardEmployee) (DashboardEmployee, error) {
	listListDashboard, listListDashboardErr := s.repository.ListDashboard(empCode, dashboardSort)

	return listListDashboard, listListDashboardErr
}

func (s *service) ListDashboardTurnover(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeTurnOver, error) {
	listListDashboard, listListDashboardErr := s.repository.ListDashboardTurnover(empCode, dashboardSort)

	return listListDashboard, listListDashboardErr
}

func (s *service) ListDashboardKontrak(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeKontrak, error) {
	listListDashboard, listListDashboardErr := s.repository.ListDashboardKontrak(empCode, dashboardSort)

	return listListDashboard, listListDashboardErr
}
