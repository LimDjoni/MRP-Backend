package company

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListCompany() ([]Company, error)
	CreateCompany(inputCompany InputCreateUpdateCompany) (Company, error)
	UpdateCompany(inputCompany InputCreateUpdateCompany, id int) (Company, error)
	DeleteCompany(id int) (bool, error)
	DetailCompany(id int) (Company, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListCompany() ([]Company, error) {
	var listCompany []Company

	getListCompanyErr := r.db.Preload(clause.Associations).Order("company_name asc").Find(&listCompany).Error

	return listCompany, getListCompanyErr
}

func (r *repository) CreateCompany(inputCompany InputCreateUpdateCompany) (Company, error) {
	var createdCompany Company

	createdCompany.CompanyName = inputCompany.CompanyName
	createdCompany.Address = inputCompany.Address
	createdCompany.Province = inputCompany.Province
	createdCompany.PhoneNumber = inputCompany.PhoneNumber
	createdCompany.FaxNumber = inputCompany.FaxNumber
	createdCompany.IndustryType = inputCompany.IndustryType

	errCreate := r.db.Create(&createdCompany).Error

	if errCreate != nil {
		return createdCompany, errCreate
	}
	return createdCompany, nil
}

func (r *repository) UpdateCompany(inputCompany InputCreateUpdateCompany, id int) (Company, error) {
	var updatedCompany Company

	errFind := r.db.Where("id = ?", id).First(&updatedCompany).Error

	if errFind != nil {
		return updatedCompany, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputCompany)

	if errorMarshal != nil {
		return updatedCompany, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedCompany, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedCompany).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedCompany, updateErr
	}

	return updatedCompany, nil
}

func (r *repository) DeleteCompany(id int) (bool, error) {
	var deletedCompany Company

	errFind := r.db.Where("id = ?", id).First(&deletedCompany).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ?", id).Delete(&deletedCompany).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) DetailCompany(id int) (Company, error) {
	var detailCompany Company

	errFind := r.db.Where("id = ?", id).First(&detailCompany).Error

	if errFind != nil {
		return detailCompany, errFind
	}

	return detailCompany, nil
}
