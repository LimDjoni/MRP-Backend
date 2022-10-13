package notification

import (
"gorm.io/gorm"
)

type Repository interface {
	CreateNotification(input InputNotification, userId uint) (Notification, error)
	GetNotification(userId uint) ([]Notification, error)
	UpdateReadNotification(userId uint) ([]Notification, error)
	DeleteNotification(userId uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateNotification(input InputNotification, userId uint) (Notification, error) {
	var createdNotification Notification

	createdNotification.Status = input.Status
	createdNotification.Type = input.Type
	createdNotification.Period = input.Period
	createdNotification.UserId = userId
	createdNotification.IsRead = false

	errCreate := r.db.Create(&createdNotification).Error

	if errCreate != nil {
		return createdNotification, errCreate
	}

	return createdNotification, nil
}

func (r *repository) GetNotification(userId uint) ([]Notification, error) {
	var listNotification []Notification

	errFind := r.db.Where("user_id = ?", userId).Find(&listNotification).Error

	if errFind != nil {
		return listNotification, errFind
	}

	return listNotification, nil
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
