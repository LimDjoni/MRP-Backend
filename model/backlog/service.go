package backlog

type Service interface {
	CreateBackLog(backlogs RegisterBackLogInput) (BackLog, error)
	FindBackLog() ([]BackLog, error)
	FindBackLogById(id uint) (BackLog, error)
	GetListBackLog(page int, sortFilter SortFilterBackLog) (Pagination, error)
	UpdateBackLog(inputBackLog RegisterBackLogInput, id int) (BackLog, error)
	DeleteBackLog(id uint) (bool, error)
	ListDashboardBackLog(dashboardSort SortFilterDashboardBacklog) (DashboardBackLog, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateBackLog(backlogs RegisterBackLogInput) (BackLog, error) {
	newBackLog, err := s.repository.CreateBackLog(backlogs)

	return newBackLog, err
}

func (s *service) FindBackLog() ([]BackLog, error) {
	backlogs, err := s.repository.FindBackLog()

	return backlogs, err
}

func (s *service) FindBackLogById(id uint) (BackLog, error) {
	alatBerat, err := s.repository.FindBackLogById(id)

	return alatBerat, err
}

func (s *service) GetListBackLog(page int, sortFilter SortFilterBackLog) (Pagination, error) {
	listBackLogs, listBackLogsErr := s.repository.ListBackLog(page, sortFilter)

	return listBackLogs, listBackLogsErr
}

func (s *service) UpdateBackLog(inputBackLog RegisterBackLogInput, id int) (BackLog, error) {
	updateBackLog, updateBackLogErr := s.repository.UpdateBackLog(inputBackLog, id)

	return updateBackLog, updateBackLogErr
}

func (s *service) DeleteBackLog(id uint) (bool, error) {
	isDeletedBackLog, isDeletedBackLogErr := s.repository.DeleteBackLog(id)

	return isDeletedBackLog, isDeletedBackLogErr
}

func (s *service) ListDashboardBackLog(dashboardSort SortFilterDashboardBacklog) (DashboardBackLog, error) {
	backlogs, backlogsErr := s.repository.ListDashboardBackLog(dashboardSort)

	return backlogs, backlogsErr
}
