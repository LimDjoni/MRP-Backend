package notification

type Service interface {
	CreateNotification(input InputNotification, userId uint) (Notification, error)
	GetNotification(userId uint) ([]Notification, error)
	UpdateReadNotification(userId uint) ([]Notification, error)
	DeleteNotification(userId uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateNotification(input InputNotification, userId uint) (Notification, error) {
	createdNotification, createdNotificationErr := s.repository.CreateNotification(input, userId)

	return createdNotification, createdNotificationErr
}

func (s *service) GetNotification(userId uint) ([]Notification, error) {
	listNotification, listNotificationErr := s.repository.GetNotification(userId)

	return listNotification, listNotificationErr
}

func (s *service) UpdateReadNotification(userId uint) ([]Notification, error) {
	updatedNotification, updatedNotificationErr := s.repository.UpdateReadNotification(userId)

	return updatedNotification, updatedNotificationErr
}

func (s *service) DeleteNotification(userId uint) (bool, error) {
	isDeletedNotification, isDeletedNotificationErr := s.repository.DeleteNotification(userId)

	return isDeletedNotification, isDeletedNotificationErr
}
