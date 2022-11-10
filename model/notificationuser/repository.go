package notificationuser

import (
	"ajebackend/model/notification"
	"ajebackend/model/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateNotification(input notification.InputNotification, userId uint) (notification.Notification, error)
	GetNotification(userId uint) ([]NotificationUser, error)
	UpdateReadNotification(userId uint) ([]NotificationUser, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateNotification(input notification.InputNotification, userId uint) (notification.Notification, error) {
	var createdNotification notification.Notification

	createdNotification.Status = input.Status
	createdNotification.Type = input.Type
	createdNotification.Period = input.Period
	createdNotification.Document = input.Document
	createdNotification.EndUser = input.EndUser
	createdNotification.UserId = userId

	tx := r.db.Begin()

	var allUsers []user.User

	errUser := tx.Find(&allUsers).Error

	if errUser != nil {
		tx.Rollback()
		return createdNotification, errUser
	}

	errCreate := tx.Create(&createdNotification).Error

	if errCreate != nil {
		tx.Rollback()
		return createdNotification, errCreate
	}

	var notificationUser []NotificationUser
	for _, v := range allUsers  {
		var notif NotificationUser
		notif.NotificationId = createdNotification.ID
		notif.UserId = v.ID

		notificationUser = append(notificationUser, notif)
	}

	errCreateNotificationUserErr := tx.Create(&notificationUser).Error

	if errCreateNotificationUserErr != nil {
		tx.Rollback()
		return createdNotification, errCreateNotificationUserErr
	}

	tx.Commit()
	return createdNotification, nil
}

func (r *repository) GetNotification(userId uint) ([]NotificationUser, error) {
	var listNotification []NotificationUser

	errFind := r.db.Preload("Notification.User").Preload(clause.Associations).Order("id desc").Where("user_id = ?", userId).Find(&listNotification).Error

	if errFind != nil {
		return listNotification, errFind
	}

	return listNotification, nil
}

func (r *repository) UpdateReadNotification(userId uint) ([]NotificationUser, error) {
	var listNotification []NotificationUser

	errFind := r.db.Preload("Notification.User").Preload(clause.Associations).Order("id desc").Where("user_id = ?", userId).Find(&listNotification).Error

	if errFind != nil {
		return listNotification, errFind
	}

	updErr := r.db.Model(&listNotification).Update("is_read", true).Error

	if updErr != nil {
		return listNotification, updErr
	}

	return listNotification, nil
}
