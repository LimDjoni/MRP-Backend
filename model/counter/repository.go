package counter

import (
	"ajebackend/model/master/iupopk"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	UpdateCounter() error
	CreateIupopk(input iupopk.InputIupopk) (iupopk.Iupopk, error)
	GetCounter(iupopkId int) (Counter, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) UpdateCounter() error {

	_, _, day := time.Now().Date()

	if day != 1 {
		return errors.New("Date should be first date of month")
	}

	counterMap := map[string]interface{}{
		"transaction_dn":      1,
		"transaction_ln":      1,
		"grouping_mv_dn":      1,
		"grouping_mv_ln":      1,
		"sp3medn":             1,
		"sp3meln":             1,
		"ba_end_user":         1,
		"dmo":                 1,
		"production":          1,
		"coa_report":          1,
		"coa_report_ln":       1,
		"insw":                1,
		"rkab":                1,
		"electric_assignment": 1,
		"caf_assignment":      1,
		"royalty_recon":       1,
		"royalty_report":      1,
	}

	updateCounterErr := r.db.Model(&Counter{}).Where("created_at <= ?", time.Now()).Updates(counterMap).Error

	if updateCounterErr != nil {
		return updateCounterErr
	}

	return nil
}

func (r *repository) CreateIupopk(input iupopk.InputIupopk) (iupopk.Iupopk, error) {
	var createdIupopk iupopk.Iupopk

	createdIupopk.Name = input.Name
	createdIupopk.Address = input.Address
	createdIupopk.Province = input.Province
	createdIupopk.Email = input.Email
	createdIupopk.PhoneNumber = input.PhoneNumber
	createdIupopk.FaxNumber = input.FaxNumber
	createdIupopk.DirectorName = input.DirectorName
	createdIupopk.Position = input.Position
	createdIupopk.Code = input.Code

	tx := r.db.Begin()

	createIupopkErr := tx.Create(&createdIupopk).Error

	if createIupopkErr != nil {
		tx.Rollback()
		return createdIupopk, createIupopkErr
	}

	// var haulingSynchronize haulingsynchronize.HaulingSynchronize

	// haulingSynchronize.IupopkId = createdIupopk.ID

	// createHaulingSyncErr := tx.Create(&haulingSynchronize).Error

	// if createHaulingSyncErr != nil {
	// 	tx.Rollback()
	// 	return createdIupopk, createHaulingSyncErr
	// }

	var createdCounter Counter

	createdCounter.IupopkId = createdIupopk.ID
	createdCounter.TransactionDn = 1
	createdCounter.TransactionLn = 1
	createdCounter.GroupingMvDn = 1
	createdCounter.GroupingMvLn = 1
	createdCounter.Sp3medn = 1
	createdCounter.Sp3meln = 1
	createdCounter.BaEndUser = 1
	createdCounter.Dmo = 1
	createdCounter.Production = 1
	createdCounter.Insw = 1
	createdCounter.CoaReport = 1
	createdCounter.CoaReportLn = 1
	createdCounter.Rkab = 1
	createdCounter.ElectricAssignment = 1
	createdCounter.CafAssignment = 1
	createdCounter.RoyaltyRecon = 1
	createdCounter.RoyaltyReport = 1
	createdCounter.BastFormat = input.BastFormat

	createCounterErr := tx.Create(&createdCounter).Error

	if createCounterErr != nil {
		tx.Rollback()
		return createdIupopk, createCounterErr
	}

	tx.Commit()
	return createdIupopk, nil
}

func (r *repository) GetCounter(iupopkId int) (Counter, error) {
	var counter Counter

	errFind := r.db.Where("iupopk_id = ?", iupopkId).First(&counter).Error

	return counter, errFind
}
