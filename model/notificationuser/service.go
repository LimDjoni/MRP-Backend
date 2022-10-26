package notificationuser

import "ajebackend/model/notification"

type Service interface {
	CreateNotification(input notification.InputNotification, userId uint) (notification.Notification, error)
	GetNotification(userId uint) ([]NotificationUser, error)
	UpdateReadNotification(userId uint) ([]NotificationUser, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateNotification(input notification.InputNotification, userId uint) (notification.Notification, error) {
	createdNotification, createdNotificationErr := s.repository.CreateNotification(input, userId)

	return createdNotification, createdNotificationErr
}

func (s *service) GetNotification(userId uint) ([]NotificationUser, error) {
	listNotification, listNotificationErr := s.repository.GetNotification(userId)

	return listNotification, listNotificationErr
}

func (s *service) UpdateReadNotification(userId uint) ([]NotificationUser, error) {
	updateReadNotification, updateReadNotificationErr := s.repository.UpdateReadNotification(userId)

	return updateReadNotification, updateReadNotificationErr
}
