package destination

import "gorm.io/gorm"

type Repository interface {
	GetDestination() ([]Destination, error)
	GetDestinationByName(name string) (Destination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetDestination() ([]Destination, error) {
	var destinations []Destination

	errFind := r.db.Find(&destinations).Error

	if errFind != nil {
		return destinations, errFind
	}

	return destinations, nil
}

func (r *repository) GetDestinationByName(name string) (Destination, error) {
	var destination Destination

	errFind := r.db.Where("name = ?", name).First(&destination).Error

	if errFind != nil {
		return destination, errFind
	}

	return destination, nil
}
