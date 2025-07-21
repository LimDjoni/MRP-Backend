package alatberat

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateAlatBerat(alatberat RegisterAlatBeratInput) (AlatBerat, error)
	FindAlatBerat() ([]AlatBerat, error)
	FindAlatBeratById(id uint) (AlatBerat, error)
	ListAlatBerat(page int, sortFilter SortFilterAlatBerat) (Pagination, error)
	FindConsumption(brandId uint, heavyEquipmentId uint, seriesId uint) (AlatBerat, error)
	UpdateAlatBerat(inputAlatBerat RegisterAlatBeratInput, id int) (AlatBerat, error)
	DeleteAlatBerat(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateAlatBerat(AlatBeratInput RegisterAlatBeratInput) (AlatBerat, error) {
	var newAlataBerat AlatBerat

	newAlataBerat.BrandId = AlatBeratInput.BrandId
	newAlataBerat.HeavyEquipmentId = AlatBeratInput.HeavyEquipmentId
	newAlataBerat.SeriesId = AlatBeratInput.SeriesId
	newAlataBerat.Consumption = AlatBeratInput.Consumption
	newAlataBerat.Tolerance = AlatBeratInput.Tolerance

	err := r.db.Create(&newAlataBerat).Error
	if err != nil {
		return newAlataBerat, err
	}

	return newAlataBerat, nil
}

func (r *repository) FindAlatBerat() ([]AlatBerat, error) {
	var alatBerat []AlatBerat

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Preload("Series").Find(&alatBerat).Error

	return alatBerat, errFind
}

func (r *repository) FindAlatBeratById(id uint) (AlatBerat, error) {
	var alatBerat AlatBerat

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Preload("Series").
		Where("id = ?", id).First(&alatBerat).Error
	return alatBerat, errFind
}

func (r *repository) ListAlatBerat(page int, sortFilter SortFilterAlatBerat) (Pagination, error) {
	var listAlatBerat []AlatBerat
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "alat_berats.id > 0"
	querySort := "alat_berats.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.BrandId != "" {
		queryFilter = queryFilter + " AND alat_berats.brand_id = " + sortFilter.BrandId
	}

	if sortFilter.HeavyEquipmentId != "" {
		queryFilter = queryFilter + " AND cast(he.heavy_equipment_name AS TEXT) LIKE '%" + sortFilter.HeavyEquipmentId + "%'"
	}

	if sortFilter.SeriesId != "" {
		queryFilter = queryFilter + " AND cast(s.series_name AS TEXT) LIKE '%" + sortFilter.SeriesId + "%'"
	}

	if sortFilter.Consumption != "" {
		queryFilter = queryFilter + " AND alat_berats.consumption = " + sortFilter.Consumption
	}

	if sortFilter.Tolerance != "" {
		queryFilter = queryFilter + " AND alat_berats.tolerance = " + sortFilter.Tolerance
	}

	errFind := r.db.Joins("JOIN heavy_equipments he ON alat_berats.heavy_equipment_id = he.id").Joins("JOIN series s ON alat_berats.series_id = s.id").Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listAlatBerat, &pagination, r.db, queryFilter)).Find(&listAlatBerat).Error
	if errFind != nil {

		return pagination, errFind

	}

	pagination.Data = listAlatBerat

	return pagination, nil
}

func (r *repository) FindConsumption(brandId uint, heavyEquipmentId uint, seriesId uint) (AlatBerat, error) {
	var alatBerat AlatBerat

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Preload("Series").
		Where("brand_id = ? AND heavy_equipment_id = ? AND series_id = ?", brandId, heavyEquipmentId, seriesId).First(&alatBerat).Error
	return alatBerat, errFind
}

func (r *repository) UpdateAlatBerat(inputAlatBerat RegisterAlatBeratInput, id int) (AlatBerat, error) {

	var updatedAlatBerat AlatBerat
	errFind := r.db.Where("id = ?", id).First(&updatedAlatBerat).Error

	if errFind != nil {
		return updatedAlatBerat, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputAlatBerat)

	if errorMarshal != nil {
		return updatedAlatBerat, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedAlatBerat, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedAlatBerat).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedAlatBerat, updateErr
	}

	return updatedAlatBerat, nil
}

func (r *repository) DeleteAlatBerat(id uint) (bool, error) {
	tx := r.db.Begin()
	var alatBerat AlatBerat

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&alatBerat).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&alatBerat).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
