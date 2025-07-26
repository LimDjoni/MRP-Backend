package fuelratio

import (
	"encoding/json"
	"math"
	"strings"
	"time"

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

func nullIfEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func dateIfNotZero(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func (r *repository) CreateFuelRatio(FuelRatioInput RegisterFuelRatioInput) (FuelRatio, error) {
	var newFuelRatio FuelRatio

	newFuelRatio.UnitId = FuelRatioInput.UnitId
	newFuelRatio.OperatorName = FuelRatioInput.OperatorName
	newFuelRatio.Shift = FuelRatioInput.Shift
	newFuelRatio.FirstHM = FuelRatioInput.FirstHM
	newFuelRatio.LastHM = FuelRatioInput.LastHM
	newFuelRatio.TotalRefill = FuelRatioInput.TotalRefill
	newFuelRatio.Status = FuelRatioInput.Status

	newFuelRatio.Tanggal = nullIfEmpty(FuelRatioInput.Tanggal)
	newFuelRatio.TanggalAwal = nullIfEmpty(FuelRatioInput.TanggalAwal)
	newFuelRatio.TanggalAkhir = nullIfEmpty(FuelRatioInput.TanggalAkhir)

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
		Preload("Unit.Series.HeavyEquipment").Find(&fuelratios).Error

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

	if sortFilter.OperatorName != "" {
		queryFilter = queryFilter + " AND LOWER(operator_name) LIKE '%" + strings.ToLower(sortFilter.OperatorName) + "%'"
	}

	if sortFilter.Shift != "" {
		queryFilter = queryFilter + " AND cast(shift AS TEXT) LIKE '%" + sortFilter.Shift + "%'"
	}

	if sortFilter.Tanggal != "" {
		queryFilter = queryFilter + " AND tanggal = '" + sortFilter.Tanggal + "'"
	}

	if sortFilter.Status != "" {
		queryFilter = queryFilter + " AND status = " + sortFilter.Status
	}

	errFind := r.db.
		Joins("JOIN units u ON fuel_ratios.unit_id = u.id").
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").Where(queryFilter).Order(querySort).Scopes(paginateDataPage(listFuelRatio, &pagination, r.db)).Find(&listFuelRatio).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listFuelRatio

	return pagination, nil
}

func (r *repository) FindFuelRatioExport(sortFilter SortFilterFuelRatioSummary) ([]SortFilterFuelRatioSummary, error) {
	var results []SortFilterFuelRatioSummary

	// Filters
	var filters []string
	var args []interface{}

	filters = append(filters, "fr.status = ?")
	args = append(args, true)

	if sortFilter.TanggalAkhir != "" {
		filters = append(filters, "fr.tanggal_akhir >= ?")
		args = append(args, sortFilter.TanggalAkhir)
	}
	if sortFilter.TanggalAwal != "" {
		filters = append(filters, "fr.tanggal_awal <= ?")
		args = append(args, sortFilter.TanggalAwal)
	}
	if sortFilter.UnitName != "" {
		filters = append(filters, "u.unit_name ILIKE ?")
		args = append(args, "%"+sortFilter.UnitName+"%")
	}
	if sortFilter.Shift != "" {
		filters = append(filters, "fr.shift ILIKE ?")
		args = append(args, "%"+sortFilter.Shift+"%")
	}
	if sortFilter.Consumption != "" {
		filters = append(filters, "CAST(ab.consumption AS TEXT) ILIKE ?")
		args = append(args, "%"+sortFilter.Consumption+"%")
	}
	if sortFilter.Tolerance != "" {
		filters = append(filters, "CAST(ab.tolerance AS TEXT) ILIKE ?")
		args = append(args, "%"+sortFilter.Tolerance+"%")
	}

	querySort := "u.unit_name desc"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	// 👇 alat_berats subquery with DISTINCT ON
	subAlatBerat := `
		(SELECT DISTINCT ON (brand_id, heavy_equipment_id, series_id) *
		FROM alat_berats) ab`

	// 👇 Build the subQuery
	subQuery := r.db.Table("fuel_ratios fr").
		Select(`
			fr.unit_id, 
			u.unit_name, 
			fr.shift, 
			SUM(fr.total_refill) AS total_refill, 
			ab.consumption, 
			ab.tolerance, 
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600
				ELSE fr.last_hm - fr.first_hm
			END) AS duration,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN (EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600) * ab.consumption
				ELSE (fr.last_hm - fr.first_hm) * ab.consumption
			END) AS batas_bawah,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN (EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600) * 
					(ab.consumption + (ab.consumption * ab.tolerance / 100))
				ELSE (fr.last_hm - fr.first_hm) * 
					(ab.consumption + (ab.consumption * ab.tolerance / 100))
			END) AS batas_atas,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600
				ELSE fr.last_hm - fr.first_hm
			END) / NULLIF(ab.consumption, 0) AS total_konsumsi_bbm
		`).
		Joins("JOIN units u ON fr.unit_id = u.id").
		Joins("JOIN "+subAlatBerat+" ON ab.brand_id = u.brand_id AND ab.heavy_equipment_id = u.heavy_equipment_id AND ab.series_id = u.series_id").
		Where(strings.Join(filters, " AND "), args...).
		Group("fr.unit_id, u.unit_name, fr.shift, ab.consumption, ab.tolerance")

	// Wrap subquery
	tx := r.db.Table("(?) AS sub", subQuery)

	// Apply post-subquery filters
	if sortFilter.TotalKonsumsiBBM != "" {
		tx = tx.Where("CAST(total_konsumsi_bbm AS TEXT) ILIKE ?", "%"+sortFilter.TotalKonsumsiBBM+"%")
	}
	if sortFilter.TotalRefill != "" {
		tx = tx.Where("CAST(total_refill AS TEXT) ILIKE ?", "%"+sortFilter.TotalRefill+"%")
	}
	if sortFilter.Duration != "" {
		tx = tx.Where("CAST(duration AS TEXT) ILIKE ?", "%"+sortFilter.Duration+"%")
	}

	// Final query execution without pagination
	if err := tx.Order(querySort).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *repository) ListRangkuman(page int, sortFilter SortFilterFuelRatioSummary) (Pagination, error) {
	var results []SortFilterFuelRatioSummary
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page

	var filters []string
	var args []interface{}

	filters = append(filters, "fr.status = ?")
	args = append(args, true)

	if sortFilter.TanggalAkhir != "" {
		filters = append(filters, "fr.tanggal_akhir >= ?")
		args = append(args, sortFilter.TanggalAkhir)
	}
	if sortFilter.TanggalAwal != "" {
		filters = append(filters, "fr.tanggal_awal <= ?")
		args = append(args, sortFilter.TanggalAwal)
	}
	if sortFilter.UnitName != "" {
		filters = append(filters, "u.unit_name ILIKE ?")
		args = append(args, "%"+sortFilter.UnitName+"%")
	}
	if sortFilter.Shift != "" {
		filters = append(filters, "fr.shift ILIKE ?")
		args = append(args, "%"+sortFilter.Shift+"%")
	}
	if sortFilter.Consumption != "" {
		filters = append(filters, "CAST(ab.consumption AS TEXT) ILIKE ?")
		args = append(args, "%"+sortFilter.Consumption+"%")
	}
	if sortFilter.Tolerance != "" {
		filters = append(filters, "CAST(ab.tolerance AS TEXT) ILIKE ?")
		args = append(args, "%"+sortFilter.Tolerance+"%")
	}

	querySort := "u.unit_name desc"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	// 👇 alat_berats subquery with DISTINCT ON
	subAlatBerat := `
		(SELECT DISTINCT ON (brand_id, heavy_equipment_id, series_id) *
		FROM alat_berats) ab`

	// 👇 Build the subQuery
	subQuery := r.db.Table("fuel_ratios fr").
		Select(`
			fr.unit_id, 
			u.unit_name, 
			fr.shift, 
			SUM(fr.total_refill) AS total_refill, 
			ab.consumption, 
			ab.tolerance, 
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600
				ELSE fr.last_hm - fr.first_hm
			END) AS duration,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN (EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600) * ab.consumption
				ELSE (fr.last_hm - fr.first_hm) * ab.consumption
			END) AS batas_bawah,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN (EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600) * 
					(ab.consumption + (ab.consumption * ab.tolerance / 100))
				ELSE (fr.last_hm - fr.first_hm) * 
					(ab.consumption + (ab.consumption * ab.tolerance / 100))
			END) AS batas_atas,
			SUM(CASE
				WHEN fr.tanggal_awal IS NOT NULL AND fr.tanggal_akhir IS NOT NULL AND fr.tanggal_awal != '' AND fr.tanggal_akhir != ''
				THEN EXTRACT(EPOCH FROM (fr.tanggal_akhir::timestamp - fr.tanggal_awal::timestamp)) / 3600
				ELSE fr.last_hm - fr.first_hm
			END) / NULLIF(ab.consumption, 0) AS total_konsumsi_bbm
		`).
		Joins("JOIN units u ON fr.unit_id = u.id").
		Joins("JOIN "+subAlatBerat+" ON ab.brand_id = u.brand_id AND ab.heavy_equipment_id = u.heavy_equipment_id AND ab.series_id = u.series_id").
		Where(strings.Join(filters, " AND "), args...).
		Group("fr.unit_id, u.unit_name, fr.shift, ab.consumption, ab.tolerance")

	// 👇 Wrap subquery
	tx := r.db.Table("(?) AS sub", subQuery)

	if sortFilter.TotalKonsumsiBBM != "" {
		tx = tx.Where("CAST(total_konsumsi_bbm AS TEXT) ILIKE ?", "%"+sortFilter.TotalKonsumsiBBM+"%")
	}
	if sortFilter.TotalRefill != "" {
		tx = tx.Where("CAST(total_refill AS TEXT) ILIKE ?", "%"+sortFilter.TotalRefill+"%")
	}
	if sortFilter.Duration != "" {
		tx = tx.Where("CAST(duration AS TEXT) ILIKE ?", "%"+sortFilter.Duration+"%")
	}

	tx = tx.Order(querySort)

	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return pagination, err
	}

	pagination.TotalRows = count
	pagination.TotalPages = int(math.Ceil(float64(count) / float64(pagination.Limit)))

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
