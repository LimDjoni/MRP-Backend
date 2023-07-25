package trader

import (
	"ajebackend/model/master/company"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListTrader() ([]Trader, error)
	CheckListTrader(list []int) ([]Trader, error)
	CheckEndUser(id int) (Trader, error)
	CreateTrader(inputTrader InputCreateUpdateTrader) (Trader, error)
	UpdateTrader(inputTrader InputCreateUpdateTrader, id int) (Trader, error)
	DeleteTrader(id int) (bool, error)
	ListTraderWithCompanyId(id int) ([]Trader, error)
	DetailTrader(id int) (Trader, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListTrader() ([]Trader, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Preload(clause.Associations).Preload("Company.IndustryType.CategoryIndustryType").Order("trader_name asc").Find(&listTrader).Error

	return listTrader, getListTraderErr
}

func (r *repository) CheckListTrader(list []int) ([]Trader, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Preload(clause.Associations).Where("id IN ?", list).Find(&listTrader).Error

	if getListTraderErr != nil {
		return listTrader, getListTraderErr
	}

	if len(listTrader) != len(list) {
		return listTrader, errors.New("record not found")
	}

	return listTrader, nil
}

func (r *repository) CheckEndUser(id int) (Trader, error) {
	var endUser Trader

	getEndUserErr := r.db.Preload(clause.Associations).Where("id = ?", id).First(&endUser).Error

	if getEndUserErr != nil {
		return endUser, getEndUserErr
	}

	return endUser, nil
}

func (r *repository) CreateTrader(inputTrader InputCreateUpdateTrader) (Trader, error) {
	var createdTrader Trader

	var findCompany company.Company
	errFindCompany := r.db.Where("id = ?", inputTrader.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return createdTrader, errFindCompany
	}

	createdTrader.TraderName = inputTrader.TraderName
	createdTrader.Position = inputTrader.Position
	var email string
	if inputTrader.Email != nil {
		email = strings.ToLower(*inputTrader.Email)
	}
	if email != "" {
		createdTrader.Email = &email
	}
	createdTrader.CompanyId = uint(inputTrader.CompanyId)

	errCreate := r.db.Create(&createdTrader).Error

	if errCreate != nil {
		return createdTrader, errCreate
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", createdTrader.ID).First(&createdTrader).Error

	if errFind != nil {
		return createdTrader, errFind
	}

	return createdTrader, nil
}

func (r *repository) UpdateTrader(inputTrader InputCreateUpdateTrader, id int) (Trader, error) {
	var updatedTrader Trader

	inputTrader.TraderName = inputTrader.TraderName
	inputTrader.Position = inputTrader.Position
	var email string
	if inputTrader.Email != nil {
		email = strings.ToLower(*inputTrader.Email)
	}
	if email != "" {
		inputTrader.Email = &email
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", id).First(&updatedTrader).Error

	if errFind != nil {
		return updatedTrader, errFind
	}

	var findCompany company.Company

	errFindCompany := r.db.Where("id = ?", inputTrader.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return updatedTrader, errFindCompany
	}

	dataInput, errorMarshal := json.Marshal(inputTrader)

	if errorMarshal != nil {
		return updatedTrader, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedTrader, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedTrader).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedTrader, updateErr
	}

	return updatedTrader, nil
}

func (r *repository) DeleteTrader(id int) (bool, error) {
	var deletedTrader Trader

	errFind := r.db.Where("id = ?", id).First(&deletedTrader).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ?", id).Delete(&deletedTrader).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) ListTraderWithCompanyId(id int) ([]Trader, error) {

	var listTraderWithCompanyId []Trader

	errFind := r.db.Preload(clause.Associations).Where("company_id = ?", id).Find(&listTraderWithCompanyId).Error

	if errFind != nil {
		return listTraderWithCompanyId, errFind
	}

	return listTraderWithCompanyId, nil
}

func (r *repository) DetailTrader(id int) (Trader, error) {

	var detailTrader Trader

	errFind := r.db.Preload(clause.Associations).Where("id = ?", id).First(&detailTrader).Error

	if errFind != nil {
		return detailTrader, errFind
	}

	return detailTrader, nil
}
