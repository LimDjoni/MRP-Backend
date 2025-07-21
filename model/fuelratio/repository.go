package fuelratio

import (
	"encoding/json"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	CreateFuelRatio(fuelratios RegisterFuelRatioInput) (FuelRatio, error)
	FindFuelRatio() ([]FuelRatio, error)
	FindFuelRatioById(id uint) (FuelRatio, error)
	ListFuelRatio(page int, sortFilter SortFilterFuelRatio) (Pagination, error)
	FindFuelRatioExport(sortFilter SortFilterFuelRatioSummary) ([]SortFilterFuelRatioSummary, error)
	ListRangkuman(page int, sortFilter SortFilterFuelRatioSummary) (Pagination, error)
	UpdateFuelRatio(inputFuelRatio RegisterFuelRatioInput, id int) (FuelRatio, error)
	DeleteFuelRatio(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateFuelRatio(FuelRatioInput RegisterFuelRatioInput) (FuelRatio, error) {
	var newFuelRatio FuelRatio

	newFuelRatio.UnitId = FuelRatioInput.UnitId
	newFuelRatio.EmployeeId = FuelRatioInput.EmployeeId
	newFuelRatio.Shift = FuelRatioInput.Shift
	newFuelRatio.FirstHM = FuelRatioInput.FirstHM

	err := r.db.Create(&newFuelRatio).Error
	if err != nil {
		return newFuelRatio, err
	}

	return newFuelRatio, nil
}

func (r *repository) FindFuelRatio() ([]FuelRatio, error) {
	var fuelratios []FuelRatio

	errFind := r.db.
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").
		Preload("Employee").Find(&fuelratios).Error

	return fuelratios, errFind
}

func (r *repository) FindFuelRatioById(id uint) (FuelRatio, error) {
	var fuelratios FuelRatio

	errFind := r.db.
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").
		Preload("Employee").
		Where("id = ?", id).First(&fuelratios).Error
	return fuelratios, errFind
}

func (r *repository) ListFuelRatio(page int, sortFilter SortFilterFuelRatio) (Pagination, error) {
	var listFuelRatio []FuelRatio
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "fuel_ratios.id > 0"
	querySort := "fuel_ratios.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.UnitId != "" {
		queryFilter = queryFilter + " AND cast(u.unit_name AS TEXT) LIKE '%" + sortFilter.UnitId + "%'"
	}

	if sortFilter.EmployeeId != "" {
		queryFilter = queryFilter + " AND cast(o.firstname AS TEXT) LIKE '%" + sortFilter.EmployeeId + "%'"
	}

	if sortFilter.Shift != "" {
		queryFilter = queryFilter + " AND cast(shift AS TEXT) LIKE '%" + sortFilter.Shift + "%'"
	}

	if sortFilter.FirstHM != "" {
		queryFilter = queryFilter + " AND first_hm >= '" + sortFilter.FirstHM + "'"
	}

	if sortFilter.Status != "" {
		queryFilter = queryFilter + " AND status = " + sortFilter.Status
	}

	errFind := r.db.
		Joins("JOIN units u ON fuel_ratios.unit_id = u.id").
		Joins("JOIN employees o ON fuel_ratios.employee_id = o.id").
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").
		Preload("Employee").Where(queryFilter).Order(querySort).Scopes(paginateData(listFuelRatio, &pagination, r.db, queryFilter)).Find(&listFuelRatio).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listFuelRatio

	return pagination, nil
}

func (r *repository) FindFuelRatioExport(sortFilter SortFilterFuelRatioSummary) ([]SortFilterFuelRatioSummary, error) {
	var results []SortFilterFuelRatioSummary

	// WHERE clauses
	var filters []string
	var args []interface{}

	filters = append(filters, "fr.status = ?")
	args = append(args, true)

	if sortFilter.FirstHM != "" {
		filters = append(filters, "fr.first_hm >= ?")
		args = append(args, sortFilter.FirstHM)
	}

	if sortFilter.LastHM != "" {
		filters = append(filters, "fr.last_hm <= ?")
		args = append(args, sortFilter.LastHM)
	}

	if sortFilter.UnitName != "" {
		filters = append(filters, "u.unit_name ILIKE ?")
		args = append(args, "%"+sortFilter.UnitName+"%")
	}

	if sortFilter.Shift != "" {
		filters = append(filters, "fr.shift ILIKE ?")
		args = append(args, "%"+sortFilter.Shift+"%")
	}

	if sortFilter.TotalRefill != "" {
		filters = append(filters, "total_refill ILIKE ?")
		args = append(args, "%"+sortFilter.TotalRefill+"%")
	}

	if sortFilter.Consumption != "" {
		filters = append(filters, "o.consumption ILIKE ?")
		args = append(args, "%"+sortFilter.Consumption+"%")
	}

	if sortFilter.Tolerance != "" {
		filters = append(filters, "o.tolerance ILIKE ?")
		args = append(args, "%"+sortFilter.Tolerance+"%")
	}

	// Default sort
	querySort := "u.unit_name desc"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	// Build main query
	errFind := r.db.Table("fuel_ratios fr").
		Select(`
			fr.unit_id, 
			u.unit_name, 
			fr.shift, 
			SUM(fr.total_refill) AS total_refill, 
			ab.consumption, 
			ab.tolerance,  
			SUM((CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) AS duration,  
			SUM(EXTRACT(EPOCH FROM (CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) / 3600 * (ab.consumption * 3600)) / 3600 AS batas_bawah,
			SUM(EXTRACT(EPOCH FROM (CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) / 3600 * (ab.consumption * 3600 + (ab.consumption * 3600 * ab.tolerance / 100))) / 3600 AS batas_atas
		`).
		Joins("JOIN units u ON fr.unit_id = u.id").
		Joins("JOIN employees o ON fr.employee_id = o.id").
		Joins("JOIN alat_berats ab ON ab.brand_id = u.brand_id AND ab.heavy_equipment_id = u.heavy_equipment_id AND ab.series_id = u.series_id").
		Where(strings.Join(filters, " AND "), args...).
		Group("fr.unit_id, u.unit_name, fr.shift, ab.consumption, ab.tolerance").
		Order(querySort).Scan(&results).Error

	return results, errFind

}

func (r *repository) ListRangkuman(page int, sortFilter SortFilterFuelRatioSummary) (Pagination, error) {
	var results []SortFilterFuelRatioSummary
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page

	// WHERE clauses
	var filters []string
	var args []interface{}

	filters = append(filters, "fr.status = ?")
	args = append(args, true)

	if sortFilter.FirstHM != "" {
		filters = append(filters, "fr.first_hm >= ?")
		args = append(args, sortFilter.FirstHM)
	}

	if sortFilter.LastHM != "" {
		filters = append(filters, "fr.last_hm <= ?")
		args = append(args, sortFilter.LastHM)
	}

	if sortFilter.UnitName != "" {
		filters = append(filters, "u.unit_name ILIKE ?")
		args = append(args, "%"+sortFilter.UnitName+"%")
	}

	if sortFilter.Shift != "" {
		filters = append(filters, "fr.shift ILIKE ?")
		args = append(args, "%"+sortFilter.Shift+"%")
	}

	if sortFilter.TotalRefill != "" {
		filters = append(filters, "total_refill ILIKE ?")
		args = append(args, "%"+sortFilter.TotalRefill+"%")
	}

	if sortFilter.Consumption != "" {
		filters = append(filters, "o.consumption ILIKE ?")
		args = append(args, "%"+sortFilter.Consumption+"%")
	}

	if sortFilter.Tolerance != "" {
		filters = append(filters, "o.tolerance ILIKE ?")
		args = append(args, "%"+sortFilter.Tolerance+"%")
	}

	// Default sort
	querySort := "u.unit_name desc"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	// Build main query
	tx := r.db.Table("fuel_ratios fr").
		Select(`
			fr.unit_id, 
			u.unit_name, 
			fr.shift, 
			SUM(fr.total_refill) AS total_refill, 
			ab.consumption, 
			ab.tolerance,  
			SUM((CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) AS duration,  
			SUM(EXTRACT(EPOCH FROM (CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) / 3600 * (ab.consumption * 3600)) / 3600 AS batas_bawah,
			SUM(EXTRACT(EPOCH FROM (CAST(fr.last_hm AS time) - CAST(fr.first_hm AS time))) / 3600 * (ab.consumption * 3600 + (ab.consumption * 3600 * ab.tolerance / 100))) / 3600 AS batas_atas
		`).
		Joins("JOIN units u ON fr.unit_id = u.id").
		Joins("JOIN employees o ON fr.employee_id = o.id").
		Joins("JOIN alat_berats ab ON ab.brand_id = u.brand_id AND ab.heavy_equipment_id = u.heavy_equipment_id AND ab.series_id = u.series_id").
		Where(strings.Join(filters, " AND "), args...).
		Group("fr.unit_id, u.unit_name, fr.shift, ab.consumption, ab.tolerance").
		Order(querySort)

	// Count total for pagination
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return pagination, err
	}
	pagination.TotalRows = count

	totalPages := int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	if err := tx.Limit(pagination.Limit).Offset(offset).Scan(&results).Error; err != nil {
		return pagination, err
	}

	pagination.Data = results
	return pagination, nil
}

func (r *repository) UpdateFuelRatio(inputFuelRatio RegisterFuelRatioInput, id int) (FuelRatio, error) {

	var updatedFuelRatio FuelRatio
	errFind := r.db.Where("id = ?", id).First(&updatedFuelRatio).Error

	if errFind != nil {
		return updatedFuelRatio, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputFuelRatio)

	if errorMarshal != nil {
		return updatedFuelRatio, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedFuelRatio, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedFuelRatio).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedFuelRatio, updateErr
	}

	return updatedFuelRatio, nil
}

func (r *repository) DeleteFuelRatio(id uint) (bool, error) {
	tx := r.db.Begin()
	var fuelratios FuelRatio

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&fuelratios).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&fuelratios).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
