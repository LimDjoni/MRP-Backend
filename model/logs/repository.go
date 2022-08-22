package logs

import "gorm.io/gorm"

type Repository interface {
	CreateLogs (log Logs)
}

type repository struct {
	db *gorm.DB
}

func NewRepository (db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateLogs (log Logs) {
	r.db.Create(&log)

	return
}
