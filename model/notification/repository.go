package notification

import (
	"gorm.io/gorm"
)

type Repository interface {
	DeleteNotification(userId uint, iupopkId int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) DeleteNotification(userId uint, iupopkId int) (bool, error) {
	var listNotification []Notification

	errFind := r.db.Where("user_id = ? AND iupopk_id = ?", userId, iupopkId).Find(&listNotification).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("user_id = ? AND is_read = ?", userId, true).Delete(&listNotification).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}
