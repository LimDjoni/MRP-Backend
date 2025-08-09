package fuelin

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	CreateFuelIn(fuelin RegisterFuelInInput) (FuelIn, error)
	FindFuelIn() ([]FuelIn, error)
	FindFuelInById(id uint) (FuelIn, error)
	ListFuelIn(page int, sortFilter SortFilterFuelIn) (Pagination, error)
	UpdateFuelIn(inputFuelIn RegisterFuelInInput, id int) (FuelIn, error)
	DeleteFuelIn(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateFuelIn(fuelin RegisterFuelInInput) (FuelIn, error) {
	var newFuelIn FuelIn

	newFuelIn.Date = fuelin.Date
	newFuelIn.Vendor = fuelin.Vendor
	newFuelIn.Code = fuelin.Code
	newFuelIn.NomorSuratJalan = fuelin.NomorSuratJalan
	newFuelIn.NomorPlatMobil = fuelin.NomorPlatMobil
	newFuelIn.Qty = fuelin.Qty
	newFuelIn.QtyNow = fuelin.QtyNow
	newFuelIn.Driver = fuelin.Driver
	newFuelIn.TujuanAwal = fuelin.TujuanAwal

	err := r.db.Create(&newFuelIn).Error
	if err != nil {
		return newFuelIn, err
	}

	return newFuelIn, nil
}

func (r *repository) FindFuelIn() ([]FuelIn, error) {
	var fuelIn []FuelIn

	errFind := r.db.Find(&fuelIn).Error

	return fuelIn, errFind
}

func (r *repository) FindFuelInById(id uint) (FuelIn, error) {
	var fuelIn FuelIn

	errFind := r.db.
		Where("id = ?", id).First(&fuelIn).Error
	return fuelIn, errFind
}

func (r *repository) ListFuelIn(page int, sortFilter SortFilterFuelIn) (Pagination, error) {
	var listFuelIn []FuelIn
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "id > 0"
	querySort := "id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.Vendor != "" {
		queryFilter = queryFilter + " AND cast(vendor AS TEXT) LIKE '%" + sortFilter.Vendor + "%'"
	}

	if sortFilter.Code != "" {
		queryFilter = queryFilter + " AND cast(code AS TEXT) LIKE '%" + sortFilter.Code + "%'"
	}

	if sortFilter.NomorSuratJalan != "" {
		queryFilter = queryFilter + " AND cast(nomor_surat_jalan AS TEXT) ILIKE '%" + sortFilter.NomorSuratJalan + "%'"
	}

	if sortFilter.NomorPlatMobil != "" {
		queryFilter = queryFilter + " AND cast(nomor_plat_mobil AS TEXT) ILIKE '%" + sortFilter.NomorPlatMobil + "%'"
	}

	if sortFilter.Qty != "" {
		queryFilter = queryFilter + " AND qty = " + sortFilter.Qty
	}

	if sortFilter.QtyNow != "" {
		queryFilter = queryFilter + " AND qty_now = " + sortFilter.QtyNow
	}

	if sortFilter.Driver != "" {
		queryFilter = queryFilter + " AND cast(driver AS TEXT) ILIKE '%" + sortFilter.Driver + "%'"
	}

	if sortFilter.TujuanAwal != "" {
		queryFilter = queryFilter + " AND cast(tujuan_awal AS TEXT) LIKE '%" + sortFilter.TujuanAwal + "%'"
	}

	errFind := r.db.Where(queryFilter).Order(querySort).Scopes(paginateData(listFuelIn, &pagination, r.db, queryFilter)).Find(&listFuelIn).Error
	if errFind != nil {

		return pagination, errFind

	}

	pagination.Data = listFuelIn

	return pagination, nil
}

func (r *repository) UpdateFuelIn(inputFuelIn RegisterFuelInInput, id int) (FuelIn, error) {

	var updatedFuelIn FuelIn
	errFind := r.db.Where("id = ?", id).First(&updatedFuelIn).Error

	if errFind != nil {
		return updatedFuelIn, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputFuelIn)

	if errorMarshal != nil {
		return updatedFuelIn, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedFuelIn, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedFuelIn).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedFuelIn, updateErr
	}

	return updatedFuelIn, nil
}

func (r *repository) DeleteFuelIn(id uint) (bool, error) {
	tx := r.db.Begin()
	var fuelIn FuelIn

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&fuelIn).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&fuelIn).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
