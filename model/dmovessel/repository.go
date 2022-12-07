package dmovessel

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetDataDmoVessel(id uint) ([]DmoVessel, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetDataDmoVessel(id uint) ([]DmoVessel, error) {
	var listDmoVessel []DmoVessel

	errFind := r.db.Where("dmo_id = ?", id).Find(&listDmoVessel).Error

	return listDmoVessel, errFind
}
