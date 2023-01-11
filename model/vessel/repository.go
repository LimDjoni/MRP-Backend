package vessel

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	CheckVessel(vesselName string) (bool, error)
	GetVessel() ([]Vessel, error)
	CreateVessel(vesselName string) (Vessel, error)
	DetailVessel(id int) (Vessel, error)
	UpdateVessel(vesselName string, id int) (Vessel, error)
	DeleteVessel(id int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CheckVessel(vesselName string) (bool, error) {
	var vessel Vessel

	errFind := r.db.Where("name = ?", vesselName).First(&vessel).Error

	if errFind != nil {
		return false, errFind
	}
	return true, nil
}

func (r *repository) GetVessel() ([]Vessel, error) {
	var listVessel []Vessel

	errFind := r.db.Order("name asc").Find(&listVessel).Error

	return listVessel, errFind
}

func (r *repository) CreateVessel(vesselName string) (Vessel, error) {
	var vessel Vessel

	vessel.Name = vesselName

	errCreate := r.db.Create(&vessel).Error

	if errCreate != nil {
		return vessel, errCreate
	}
	return vessel, nil
}

func (r *repository) DetailVessel(id int) (Vessel, error) {
	var vessel Vessel

	errFind := r.db.Where("id = ?", id).First(&vessel).Error

	if errFind != nil {
		return vessel, errFind
	}
	return vessel, nil
}

func (r *repository) UpdateVessel(inputVessel InputVessel, id int) (Vessel, error) {
	var updatedVessel Vessel

	errFind := r.db.Where("id = ?", id).First(&updatedVessel).Error

	if errFind != nil {
		return updatedVessel, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputVessel)

	if errorMarshal != nil {
		return updatedVessel, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedVessel, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedVessel).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedVessel, updateErr
	}

	return updatedVessel, nil
}

func (r *repository) DeleteVessel(id int) (bool, error) {
	var vessel Vessel

	errFind := r.db.Where("id = ?", id).First(&vessel).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ?", id).Delete(&vessel).Error

	if errDelete != nil {
		return false, errDelete
	}
	return true, nil
}
