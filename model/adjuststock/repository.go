package adjuststock

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAdjustStock(adjuststock RegisterAdjustStockInput) (AdjustStock, error)
	FindAdjustStock() ([]AdjustStock, error)
	FindAdjustStockById(id uint) (AdjustStock, error)
	ListAdjustStock(page int, sortFilter SortFilterAdjustStock) (Pagination, error)
	UpdateAdjustStock(inputAdjustStock RegisterAdjustStockInput, id int) (AdjustStock, error)
	DeleteAdjustStock(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateAdjustStock(AdjustStockInput RegisterAdjustStockInput) (AdjustStock, error) {
	var newAdjustStock AdjustStock

	newAdjustStock.Date = AdjustStockInput.Date
	newAdjustStock.Stock = AdjustStockInput.Stock

	err := r.db.Create(&newAdjustStock).Error
	if err != nil {
		return newAdjustStock, err
	}

	return newAdjustStock, nil
}

func (r *repository) FindAdjustStock() ([]AdjustStock, error) {
	var adjustStock []AdjustStock

	errFind := r.db.Find(&adjustStock).Error

	return adjustStock, errFind
}

func (r *repository) FindAdjustStockById(id uint) (AdjustStock, error) {
	var adjustStock AdjustStock

	errFind := r.db.
		Where("id = ?", id).First(&adjustStock).Error
	return adjustStock, errFind
}

func (r *repository) ListAdjustStock(page int, sortFilter SortFilterAdjustStock) (Pagination, error) {
	var listAdjustStock []AdjustStock
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "id > 0"
	querySort := "id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.Stock != "" {
		queryFilter = queryFilter + " AND stock = " + sortFilter.Stock
	}

	if sortFilter.Date != "" {
		queryFilter = queryFilter + " AND cast(date AS TEXT) LIKE '%" + sortFilter.Date + "%'"
	}

	errFind := r.db.Where(queryFilter).Order(querySort).Scopes(paginateData(listAdjustStock, &pagination, r.db, queryFilter)).Find(&listAdjustStock).Error
	if errFind != nil {

		return pagination, errFind

	}

	pagination.Data = listAdjustStock

	return pagination, nil
}

func (r *repository) UpdateAdjustStock(inputAdjustStock RegisterAdjustStockInput, id int) (AdjustStock, error) {

	var updatedAdjustStock AdjustStock
	errFind := r.db.Where("id = ?", id).First(&updatedAdjustStock).Error

	if errFind != nil {
		return updatedAdjustStock, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputAdjustStock)

	if errorMarshal != nil {
		return updatedAdjustStock, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedAdjustStock, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedAdjustStock).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedAdjustStock, updateErr
	}

	return updatedAdjustStock, nil
}

func (r *repository) DeleteAdjustStock(id uint) (bool, error) {
	tx := r.db.Begin()
	var adjustStock AdjustStock

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&adjustStock).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&adjustStock).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
