package notification

type Service interface {
	DeleteNotification(userId uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DeleteNotification(userId uint) (bool, error) {
	isDeletedNotification, isDeletedNotificationErr := s.repository.DeleteNotification(userId)

	return isDeletedNotification, isDeletedNotificationErr
}
