package transactionshauling

import (
	"fmt"

	"ajebackend/model/transactionshauling/transactionispjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListStockRom(page int, iupopkId int) (Pagination, error)
	ListTransactionHauling(page int, iupopkId int) (Pagination, error)
	DetailStockRom(iupopkId int, stockRomId int) (transactiontoisp.TransactionToIsp, error)
	DetailTransactionHauling(iupopkId int, transactionHaulingId int) (transactionispjetty.TransactionIspJetty, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListStockRom(page int, iupopkId int) (Pagination, error) {
	var listStockRom []transactiontoisp.TransactionToIsp

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "created_at desc"

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	errFind := r.db.Preload(clause.Associations).Order(defaultSort).Where(queryFilter).Scopes(paginateData(listStockRom, &pagination, r.db, queryFilter)).Find(&listStockRom).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listStockRom

	return pagination, nil
}

func (r *repository) ListTransactionHauling(page int, iupopkId int) (Pagination, error) {
	var listTransactionHauling []transactionispjetty.TransactionIspJetty

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "created_at desc"

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	errFind := r.db.Preload(clause.Associations).Preload("TransactionToJetty.Truck").Preload("TransactionToJetty.Isp").Preload("TransactionToJetty.Site").Preload("TransactionToJetty.CreatedBy").Preload("TransactionToJetty.UpdatedBy").Preload("TransactionJetty.Jetty").Preload("TransactionJetty.CreatedBy").Preload("TransactionJetty.UpdatedBy").Order(defaultSort).Where(queryFilter).Scopes(paginateData(listTransactionHauling, &pagination, r.db, queryFilter)).Find(&listTransactionHauling).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listTransactionHauling

	return pagination, nil
}

func (r *repository) DetailStockRom(iupopkId int, stockRomId int) (transactiontoisp.TransactionToIsp, error) {
	var transactionStockRom transactiontoisp.TransactionToIsp

	errFind := r.db.Preload(clause.Associations).Where("id = ? and iupopk_id = ?", stockRomId, iupopkId).First(&transactionStockRom).Error

	if errFind != nil {
		return transactionStockRom, errFind
	}

	return transactionStockRom, nil
}

func (r *repository) DetailTransactionHauling(iupopkId int, transactionHaulingId int) (transactionispjetty.TransactionIspJetty, error) {
	var transactionHauling transactionispjetty.TransactionIspJetty

	errFind := r.db.Preload(clause.Associations).Preload("TransactionToJetty.Truck").Preload("TransactionToJetty.Isp").Preload("TransactionToJetty.Site").Preload("TransactionToJetty.CreatedBy").Preload("TransactionToJetty.UpdatedBy").Preload("TransactionJetty.Jetty").Preload("TransactionJetty.CreatedBy").Preload("TransactionJetty.UpdatedBy").Where("id = ? and iupopk_id = ?", transactionHaulingId, iupopkId).First(&transactionHauling).Error

	if errFind != nil {
		return transactionHauling, errFind
	}

	return transactionHauling, nil
}
