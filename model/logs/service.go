package logs

type Service interface {
	CreateLogs (log Logs)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateLogs (log Logs) {
	s.repository.CreateLogs(log)
	return
}
