package notification

type Service interface {
	UpdateReadNotification(userId uint) ([]Notification, error)
	DeleteNotification(userId uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) UpdateReadNotification(userId uint) ([]Notification, error) {
	updatedNotification, updatedNotificationErr := s.repository.UpdateReadNotification(userId)

	return updatedNotification, updatedNotificationErr
}

func (s *service) DeleteNotification(userId uint) (bool, error) {
	isDeletedNotification, isDeletedNotificationErr := s.repository.DeleteNotification(userId)

	return isDeletedNotification, isDeletedNotificationErr
}
