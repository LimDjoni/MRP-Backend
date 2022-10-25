package notification

import (
	"gorm.io/gorm"
)

type Repository interface {
	UpdateReadNotification(userId uint) ([]Notification, error)
	DeleteNotification(userId uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) UpdateReadNotification(userId uint) ([]Notification, error) {
	var listNotification []Notification

	errFind := r.db.Where("user_id = ?", userId).Find(&listNotification).Error

	if errFind != nil {
		return listNotification, errFind
	}

	updErr := r.db.Model(&listNotification).Update("is_read", true).Error

	if updErr != nil {
		return listNotification, updErr
	}

	return listNotification, nil
}

func (r *repository) DeleteNotification(userId uint) (bool, error) {
	var listNotification []Notification

	errFind := r.db.Where("user_id = ?", userId).Find(&listNotification).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("user_id = ? AND is_read = ?", userId, true).Delete(&listNotification).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}
