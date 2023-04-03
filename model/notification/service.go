package notification

type Service interface {
	DeleteNotification(userId uint, iupopkId int) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DeleteNotification(userId uint, iupopkId int) (bool, error) {
	isDeletedNotification, isDeletedNotificationErr := s.repository.DeleteNotification(userId, iupopkId)

	return isDeletedNotification, isDeletedNotificationErr
}
