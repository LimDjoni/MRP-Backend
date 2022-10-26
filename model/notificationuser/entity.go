package notificationuser

import (
	"ajebackend/model/notification"
	"ajebackend/model/user"
	"gorm.io/gorm"
)

type NotificationUser struct {
	gorm.Model
	NotificationId uint `json:"notification_id"`
	Notification notification.Notification `json:"notification" gorm:"constraint:OnDelete:CASCADE;"`
	IsRead bool `json:"is_read"`
	UserId uint `json:"user_id"`
	User user.User `json:"user"`
}
