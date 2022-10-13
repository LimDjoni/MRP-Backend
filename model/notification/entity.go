package notification

import (
	"ajebackend/model/user"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Status string `json:"status"`
	Type string `json:"type"`
	Period string `json:"period"`
	IsRead bool `json:"is_read"`
	UserId uint `json:"user_id"`
	User user.User
}
