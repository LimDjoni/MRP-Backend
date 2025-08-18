package unit

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateUnit(units RegisterUnitInput) (Unit, error)
	FindUnit() ([]Unit, error)
	FindUnitById(id uint) (Unit, error)
	ListUnit(page int, sortFilter SortFilterUnit) (Pagination, error)
	UpdateUnit(inputUnit RegisterUnitInput, id int) (Unit, error)
	DeleteUnit(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUnit(UnitInput RegisterUnitInput) (Unit, error) {
	var newUnit Unit

	newUnit.UnitName = UnitInput.UnitName
	newUnit.BrandId = UnitInput.BrandId
	newUnit.HeavyEquipmentId = UnitInput.HeavyEquipmentId
	newUnit.SeriesId = UnitInput.SeriesId

	err := r.db.Create(&newUnit).Error
	if err != nil {
		return newUnit, err
	}

	return newUnit, nil
}

func (r *repository) FindUnit() ([]Unit, error) {
	var units []Unit

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Preload("Series").Order("units.unit_name").Find(&units).Error

	return units, errFind
}

func (r *repository) FindUnitById(id uint) (Unit, error) {
	var units Unit

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Preload("Series").
		Where("id = ?", id).First(&units).Error
	return units, errFind
}

func (r *repository) ListUnit(page int, sortFilter SortFilterUnit) (Pagination, error) {
	var listUnit []Unit
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "units.id > 0"
	querySort := "units.unit_name"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.UnitName != "" {
		queryFilter = queryFilter + " AND cast(unit_name AS TEXT) LIKE '%" + sortFilter.UnitName + "%'"
	}

	if sortFilter.BrandId != "" {
		queryFilter = queryFilter + " AND units.brand_id = " + sortFilter.BrandId
	}

	if sortFilter.HeavyEquipmentId != "" {
		queryFilter = queryFilter + " AND cast(he.heavy_equipment_name AS TEXT) LIKE '%" + sortFilter.HeavyEquipmentId + "%'"
	}

	if sortFilter.SeriesId != "" {
		queryFilter = queryFilter + " AND cast(s.series_name AS TEXT) LIKE '%" + sortFilter.SeriesId + "%'"
	}

	errFind := r.db.Joins("JOIN heavy_equipments he ON units.heavy_equipment_id = he.id").Joins("JOIN series s ON units.series_id = s.id").Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listUnit, &pagination, r.db, queryFilter)).Find(&listUnit).Error
	if errFind != nil {

		return pagination, errFind

	}

	pagination.Data = listUnit

	return pagination, nil
}

func (r *repository) UpdateUnit(inputUnit RegisterUnitInput, id int) (Unit, error) {

	var updatedUnit Unit
	errFind := r.db.Where("id = ?", id).First(&updatedUnit).Error

	if errFind != nil {
		return updatedUnit, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputUnit)

	if errorMarshal != nil {
		return updatedUnit, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedUnit, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedUnit).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedUnit, updateErr
	}

	return updatedUnit, nil
}

func (r *repository) DeleteUnit(id uint) (bool, error) {
	tx := r.db.Begin()
	var units Unit

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&units).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&units).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
