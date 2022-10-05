package trader

import (
	"ajebackend/model/company"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Repository interface {
	ListTrader() ([]Trader, error)
	CheckListTrader(list []uint) (bool, error)
	CheckEndUser(id uint) (bool, error)
	CreateTrader(inputTrader InputTrader) (Trader, error)
	UpdateTrader(inputTrader InputTrader, id int) (Trader, error)
	DeleteTrader(id int) (bool, error)
	ListTraderWithCompanyId(id int) ([]Trader, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListTrader() ([]Trader, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Preload(clause.Associations).Find(&listTrader).Error

	return listTrader, getListTraderErr
}

func (r *repository) CheckListTrader(list []uint) (bool, error) {
	var listTrader []Trader

	getListTraderErr := r.db.Where("id IN ?", list).Find(&listTrader).Error

	if getListTraderErr != nil {
		return false, getListTraderErr
	}

	if len(listTrader) != len(list) {
		return false, errors.New("record not found")
	}

	return true, nil
}

func (r *repository) CheckEndUser(id uint) (bool, error) {
	var endUser Trader

	getEndUserErr := r.db.Where("id = ?", id).First(&endUser).Error

	if getEndUserErr != nil {
		return false, getEndUserErr
	}

	return true, nil
}

func (r *repository) CreateTrader(inputTrader InputTrader) (Trader, error) {
	var createdTrader Trader

	var findCompany company.Company
	errFindCompany := r.db.Where("id = ?", inputTrader.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return createdTrader, errFindCompany
	}

	createdTrader.TraderName = strings.ToUpper(inputTrader.TraderName)
	createdTrader.Position = strings.ToUpper(inputTrader.Position)
	createdTrader.CompanyId = inputTrader.CompanyId

	errCreate := r.db.Create(&createdTrader).Error

	if errCreate != nil {
		return createdTrader, errCreate
	}

	return createdTrader, nil
}

func (r *repository) UpdateTrader(inputTrader InputTrader, id int) (Trader, error) {
	var updatedTrader Trader

	inputTrader.TraderName = strings.ToUpper(inputTrader.TraderName)
	inputTrader.Position = strings.ToUpper(inputTrader.Position)

	errFind := r.db.Where("id = ?", id).First(&updatedTrader).Error

	if errFind != nil {
		return  updatedTrader, errFind
	}

	var findCompany company.Company

	errFindCompany := r.db.Where("id = ?", inputTrader.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return  updatedTrader, errFindCompany
	}

	dataInput, errorMarshal := json.Marshal(inputTrader)

	if errorMarshal != nil {
		return  updatedTrader, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return  updatedTrader, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedTrader).Updates(dataInputMapString).Error

	if updateErr != nil {
		return  updatedTrader, updateErr
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

	errFind := r.db.Where("company_id = ?", id).Find(&listTraderWithCompanyId).Error

	if errFind != nil {
		return listTraderWithCompanyId, errFind
	}

	return listTraderWithCompanyId, nil
}
