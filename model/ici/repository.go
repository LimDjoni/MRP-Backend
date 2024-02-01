package ici

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	GetAllIci() ([]Ici, error)
	CreateIci(inputIci InputCreateUpdateIci, IupopkId int) (Ici, error)
	UpdateIci(inputIci InputCreateUpdateIci, id int) (Ici, error)
	// DeleteIci(id int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAllIci() ([]Ici, error) {
	var icis []Ici

	errFind := r.db.Find(&icis).Error

	if errFind != nil {
		return icis, errFind
	}

	return icis, nil
}

func (r *repository) CreateIci(inputIci InputCreateUpdateIci, IupopkId int) (Ici, error) {
	var createIci Ici

	createIci.Date = inputIci.Date
	createIci.Level = inputIci.Level
	createIci.Avarage = inputIci.Avarage
	createIci.UnitPrice = inputIci.UnitPrice
	createIci.Currency = inputIci.Currency
	createIci.IupopkId = uint(IupopkId)

	errCreate := r.db.Create(&createIci).Error

	if errCreate != nil {
		return createIci, errCreate
	}

	return createIci, nil
}

func (r *repository) UpdateIci(inputIci InputCreateUpdateIci, id int) (Ici, error) {

	var updatedIci Ici
	errFind := r.db.Where("id = ?", id).First(&updatedIci).Error

	if errFind != nil {
		return updatedIci, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputIci)

	if errorMarshal != nil {
		return updatedIci, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedIci, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedIci).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedIci, updateErr
	}

	return updatedIci, nil
}
