package backlog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBackLog(backlogs RegisterBackLogInput) (BackLog, error)
	FindBackLog() ([]BackLog, error)
	FindBackLogById(id uint) (BackLog, error)
	ListBackLog(page int, sortFilter SortFilterBackLog) (Pagination, error)
	UpdateBackLog(inputBackLog RegisterBackLogInput, id int) (BackLog, error)
	DeleteBackLog(id uint) (bool, error)
	ListDashboardBackLog(dashboardSort SortFilterDashboardBacklog) (DashboardBackLog, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func nullIfEmpty(s *string) *string {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	return s
}

func dateIfNotZero(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func (r *repository) CreateBackLog(BackLogInput RegisterBackLogInput) (BackLog, error) {
	var newBackLog BackLog

	newBackLog.UnitId = BackLogInput.UnitId
	newBackLog.HMBreakdown = BackLogInput.HMBreakdown
	newBackLog.Problem = BackLogInput.Problem
	newBackLog.Component = BackLogInput.Component
	newBackLog.PartNumber = BackLogInput.PartNumber
	newBackLog.PartDescription = BackLogInput.PartDescription
	newBackLog.QtyOrder = BackLogInput.QtyOrder
	newBackLog.HMReady = BackLogInput.HMReady
	newBackLog.PPNumber = BackLogInput.PPNumber
	newBackLog.PONumber = BackLogInput.PONumber
	newBackLog.Status = BackLogInput.Status
	newBackLog.DateOfInspection = BackLogInput.DateOfInspection

	newBackLog.PlanReplaceRepair = nullIfEmpty(BackLogInput.PlanReplaceRepair)

	err := r.db.Create(&newBackLog).Error
	if err != nil {
		return newBackLog, err
	}

	return newBackLog, nil
}

func (r *repository) FindBackLog() ([]BackLog, error) {
	var backlogs []BackLog

	errFind := r.db.
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").Find(&backlogs).Error

	return backlogs, errFind
}

func (r *repository) FindBackLogById(id uint) (BackLog, error) {
	var backlogs BackLog

	errFind := r.db.
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").
		Where("id = ?", id).First(&backlogs).Error
	return backlogs, errFind
}

func (r *repository) ListBackLog(page int, sortFilter SortFilterBackLog) (Pagination, error) {
	var listBackLog []BackLog
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	queryFilter := "bl.id > 0 AND bl.deleted_at IS NULL"
	querySort := "bl.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.DateOfInspection != "" {
		parsed, err := time.Parse("2006-01-02", sortFilter.DateOfInspection)
		if err == nil {
			dateOnly := parsed.Format("2006-01-02")
			queryFilter = queryFilter + " AND bl.date_of_inspection = '" + dateOnly + "'"
		}
	}

	if sortFilter.PlanReplaceRepair != "" {
		parsed, err := time.Parse("2006-01-02", sortFilter.PlanReplaceRepair)
		if err == nil {
			dateOnly := parsed.Format("2006-01-02")
			queryFilter = queryFilter + " AND bl.plan_replace_repair = '" + dateOnly + "'"
		}
	}

	if sortFilter.BrandName != "" {
		queryFilter = queryFilter + " AND CONCAT(b.brand_name, ' ', s.series_name) ILIKE '%" + sortFilter.BrandName + "%'"
	}
	if sortFilter.UnitId != "" {
		queryFilter = queryFilter + " AND CAST(bl.unit_id AS TEXT) ILIKE '%" + sortFilter.UnitId + "%'"
	}
	if sortFilter.Problem != "" {
		queryFilter = queryFilter + " AND CAST(bl.problem AS TEXT) ILIKE '%" + sortFilter.Problem + "%'"
	}
	if sortFilter.Component != "" {
		queryFilter = queryFilter + " AND CAST(bl.component AS TEXT) ILIKE '%" + sortFilter.Component + "%'"
	}
	if sortFilter.PartNumber != "" {
		queryFilter = queryFilter + " AND CAST(bl.part_number AS TEXT) ILIKE '%" + sortFilter.PartNumber + "%'"
	}
	if sortFilter.PartDescription != "" {
		queryFilter = queryFilter + " AND CAST(bl.part_description AS TEXT) ILIKE '%" + sortFilter.PartDescription + "%'"
	}
	if sortFilter.QtyOrder != "" {
		queryFilter = queryFilter + " AND CAST(bl.qty_order AS TEXT) ILIKE '%" + sortFilter.QtyOrder + "%'"
	}
	if sortFilter.PPNumber != "" {
		queryFilter = queryFilter + " AND CAST(bl.pp_number AS TEXT) ILIKE '%" + sortFilter.PPNumber + "%'"
	}
	if sortFilter.PONumber != "" {
		queryFilter = queryFilter + " AND CAST(bl.po_number AS TEXT) ILIKE '%" + sortFilter.PONumber + "%'"
	}
	if sortFilter.Status != "" {
		queryFilter = queryFilter + " AND CAST(bl.status AS TEXT) ILIKE '%" + sortFilter.Status + "%'"
	}

	errFind := r.db.
		Unscoped().
		Table("back_logs bl").
		Select(`bl.id, bl.unit_id, bl.hm_breakdown, bl.problem, bl.component, 
         bl.part_number, bl.part_description, bl.qty_order, bl.date_of_inspection, 
         bl.plan_replace_repair, bl.hm_ready, bl.pp_number, bl.po_number, bl.status, 
         bl.created_at, bl.updated_at, bl.deleted_at,
         CASE 
			WHEN bl.plan_replace_repair IS NULL THEN 0
			ELSE CAST(bl.plan_replace_repair AS DATE) - CAST(bl.date_of_inspection AS DATE) 
		END AS aging_backlog_by_date`).
		Joins("JOIN units u ON bl.unit_id = u.id").
		Joins("JOIN brands b ON b.id = u.brand_id").
		Joins("JOIN series s ON s.id = u.series_id").
		Preload("Unit").
		Preload("Unit.Brand").
		Preload("Unit.HeavyEquipment").
		Preload("Unit.Series").
		Preload("Unit.Series.Brand").
		Preload("Unit.Series.HeavyEquipment").Where(queryFilter).Order(querySort).Scopes(paginateDataPage(listBackLog, &pagination, r.db, queryFilter)).Find(&listBackLog).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listBackLog

	return pagination, nil
}

func (r *repository) UpdateBackLog(inputBackLog RegisterBackLogInput, id int) (BackLog, error) {

	var updatedBackLog BackLog
	errFind := r.db.Where("id = ?", id).First(&updatedBackLog).Error

	if errFind != nil {
		return updatedBackLog, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputBackLog)

	if errorMarshal != nil {
		return updatedBackLog, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedBackLog, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedBackLog).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedBackLog, updateErr
	}

	return updatedBackLog, nil
}

func (r *repository) DeleteBackLog(id uint) (bool, error) {
	tx := r.db.Begin()
	var backlogs BackLog

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&backlogs).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&backlogs).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) ListDashboardBackLog(dashboardSort SortFilterDashboardBacklog) (DashboardBackLog, error) {
	var dashboardBackLogs DashboardBackLog
	var totalBackLogCount int64

	queryFilter := "back_logs.id > 0 AND back_logs.deleted_at IS NULL"

	if dashboardSort.Year != "" {
		queryFilter = queryFilter + " AND EXTRACT(YEAR FROM CAST(back_logs.date_of_inspection AS DATE)) = " + fmt.Sprintf("%s", dashboardSort.Year)
	}

	// === Queries ===
	r.db.Model(&BackLog{}).Where(queryFilter).Count(&totalBackLogCount)
	dashboardBackLogs.TotalBackLog = uint(totalBackLogCount)

	type Totals struct {
		Total1 uint
		Total2 uint
		Total3 uint
		Total4 uint
	}

	var totals Totals
	err := r.db.Model(&BackLog{}).
		Select(`
		COUNT(
			CASE 
				WHEN back_logs.plan_replace_repair IS NULL THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 THEN 1
			END
		) AS total1,
		
		COUNT(
			CASE 
				WHEN back_logs.plan_replace_repair IS NOT NULL 
				AND (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 THEN 1
			END
		) AS total2,
		
		COUNT(
			CASE 
				WHEN back_logs.plan_replace_repair IS NOT NULL 
				AND (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 THEN 1
			END
		) AS total3,
		
		COUNT(
			CASE 
				WHEN back_logs.plan_replace_repair IS NOT NULL 
				AND (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 THEN 1
			END
		) AS total4
	`).
		Where(queryFilter).
		Scan(&totals).Error
	if err != nil {
		return dashboardBackLogs, err
	}

	dashboardBackLogs.Total1 = totals.Total1
	dashboardBackLogs.Total2 = totals.Total2
	dashboardBackLogs.Total3 = totals.Total3
	dashboardBackLogs.Total4 = totals.Total4

	var backlogSummary BacklogSummary
	errSummary := r.db.Model(&BackLog{}).
		Select(`
		COUNT(CASE 
			WHEN back_logs.status LIKE 'PENDING'  
			THEN 1 END) AS pending,
		COUNT(CASE 
			WHEN back_logs.status LIKE 'OPEN'
			THEN 1 END) AS open,
		COUNT(CASE 
			WHEN back_logs.status LIKE 'CLOSED'
			THEN 1 END) AS closed,
		COUNT(CASE 
			WHEN back_logs.status LIKE 'CANCELLED'
			THEN 1 END) AS cancelled,
		COUNT(CASE 
			WHEN back_logs.status LIKE 'REJECTED'
			THEN 1 END) AS rejected
	`).
		Where(queryFilter).
		Scan(&backlogSummary).Error

	if errSummary != nil {
		return dashboardBackLogs, errSummary
	}
	dashboardBackLogs.Summary = backlogSummary

	var rawAging rawAging

	errAging := r.db.Model(&BackLog{}).
		Select(`
			COUNT(CASE 
				WHEN back_logs.plan_replace_repair IS NULL AND back_logs.status='PENDING' THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 AND back_logs.status='PENDING'
				THEN 1 END) AS pending0_5,

			COUNT(CASE 
				WHEN back_logs.plan_replace_repair IS NULL AND back_logs.status='OPEN' THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 AND back_logs.status='OPEN'
				THEN 1 END) AS open0_5, 

			COUNT(CASE 
				WHEN back_logs.plan_replace_repair IS NULL AND back_logs.status='CLOSED' THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 AND back_logs.status='CLOSED'
				THEN 1 END) AS closed0_5,  

			COUNT(CASE 
				WHEN back_logs.plan_replace_repair IS NULL AND back_logs.status='CANCELLED' THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 AND back_logs.status='CANCELLED'
				THEN 1 END) AS cancelled0_5,   

			COUNT(CASE 
				WHEN back_logs.plan_replace_repair IS NULL AND back_logs.status='REJECTED' THEN 1
				WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 0 AND 5 AND back_logs.status='REJECTED'
				THEN 1 END) AS rejected0_5,    

			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 AND back_logs.status='PENDING' THEN 1 END) AS pending6_15,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 AND back_logs.status='OPEN' THEN 1 END) AS open6_15,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 AND back_logs.status='CLOSED' THEN 1 END) AS closed6_15,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 AND back_logs.status='CANCELLED' THEN 1 END) AS cancelled6_15,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 6 AND 15 AND back_logs.status='REJECTED' THEN 1 END) AS rejected6_15,

			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 AND back_logs.status='PENDING' THEN 1 END) AS pending16_30,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 AND back_logs.status='OPEN' THEN 1 END) AS open16_30,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 AND back_logs.status='CLOSED' THEN 1 END) AS closed16_30,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 AND back_logs.status='CANCELLED' THEN 1 END) AS cancelled16_30,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) BETWEEN 16 AND 30 AND back_logs.status='REJECTED' THEN 1 END) AS rejected16_30,

			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 AND back_logs.status='PENDING' THEN 1 END) AS pending30plus,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 AND back_logs.status='OPEN' THEN 1 END) AS open30plus,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 AND back_logs.status='CLOSED' THEN 1 END) AS closed30plus,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 AND back_logs.status='CANCELLED' THEN 1 END) AS cancelled30plus,
			COUNT(CASE WHEN (CAST(back_logs.plan_replace_repair AS DATE) - CAST(back_logs.date_of_inspection AS DATE)) > 30 AND back_logs.status='REJECTED' THEN 1 END) AS rejected30plus
		`).
		Where(queryFilter).
		Scan(&rawAging).Error

	if errAging != nil {
		return dashboardBackLogs, errAging
	}

	// After scanning from the DB:
	dashboardBackLogs.AgingSummary = AgingSummary{
		AgingTotal1: BacklogSummary{
			Pending:   rawAging.Pending0_5,
			Open:      rawAging.Open0_5,
			Closed:    rawAging.Closed0_5,
			Cancelled: rawAging.Cancelled0_5,
			Rejected:  rawAging.Rejected0_5,
		},
		AgingTotal2: BacklogSummary{
			Pending:   rawAging.Pending6_15,
			Open:      rawAging.Open6_15,
			Closed:    rawAging.Closed6_15,
			Cancelled: rawAging.Cancelled6_15,
			Rejected:  rawAging.Rejected6_15,
		},
		AgingTotal3: BacklogSummary{
			Pending:   rawAging.Pending16_30,
			Open:      rawAging.Open16_30,
			Closed:    rawAging.Closed16_30,
			Cancelled: rawAging.Cancelled16_30,
			Rejected:  rawAging.Rejected16_30,
		},
		AgingTotal4: BacklogSummary{
			Pending:   rawAging.Pending30plus,
			Open:      rawAging.Open30plus,
			Closed:    rawAging.Closed30plus,
			Cancelled: rawAging.Cancelled30plus,
			Rejected:  rawAging.Rejected30plus,
		},
	}

	return dashboardBackLogs, nil
}
