package notificationuser

import "ajebackend/model/notification"

type Service interface {
	CreateNotification(input notification.InputNotification, userId uint, iupopkId int) (notification.Notification, error)
	GetNotification(userId uint, iupopkId int) ([]NotificationUser, error)
	UpdateReadNotification(userId uint, iupopkId int) ([]NotificationUser, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateNotification(input notification.InputNotification, userId uint, iupopkId int) (notification.Notification, error) {
	createdNotification, createdNotificationErr := s.repository.CreateNotification(input, userId, iupopkId)

	return createdNotification, createdNotificationErr
}

func (s *service) GetNotification(userId uint, iupopkId int) ([]NotificationUser, error) {
	listNotification, listNotificationErr := s.repository.GetNotification(userId, iupopkId)

	return listNotification, listNotificationErr
}

func (s *service) UpdateReadNotification(userId uint, iupopkId int) ([]NotificationUser, error) {
	updateReadNotification, updateReadNotificationErr := s.repository.UpdateReadNotification(userId, iupopkId)

	return updateReadNotification, updateReadNotificationErr
}
