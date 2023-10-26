package transactionshauling

import (
	"fmt"
	"time"

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
	SummaryJettyTransactionPerDay(iupopkId int) (SummaryJettyTransactionPerDay, error)
	SummaryInventoryStockRom(iupopkId int) (Summary, error)
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

	errFind := r.db.Preload(clause.Associations).Preload("TransactionToJetty.Truck.Contractor").Preload("TransactionToJetty.Pit").Preload("TransactionToJetty.Isp").Preload("TransactionToJetty.CreatedBy").Preload("TransactionToJetty.UpdatedBy").Preload("TransactionJetty.Jetty").Preload("TransactionJetty.CreatedBy").Preload("TransactionJetty.UpdatedBy").Order(defaultSort).Where(queryFilter).Scopes(paginateData(listTransactionHauling, &pagination, r.db, queryFilter)).Find(&listTransactionHauling).Error

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

	errFind := r.db.Preload(clause.Associations).Preload("TransactionToJetty.Truck.Contractor").Preload("TransactionToJetty.Pit").Preload("TransactionToJetty.Isp").Preload("TransactionToJetty.CreatedBy").Preload("TransactionToJetty.UpdatedBy").Preload("TransactionJetty.Jetty").Preload("TransactionJetty.CreatedBy").Preload("TransactionJetty.UpdatedBy").Where("id = ? and iupopk_id = ?", transactionHaulingId, iupopkId).First(&transactionHauling).Error

	if errFind != nil {
		return transactionHauling, errFind
	}

	return transactionHauling, nil
}

func (r *repository) SummaryJettyTransactionPerDay(iupopkId int) (SummaryJettyTransactionPerDay, error) {
	var summary SummaryJettyTransactionPerDay

	year, month, date := time.Now().Date()

	startDate := fmt.Sprintf("%v-%v-%vT00:00:00", year, int(month), date)
	endDate := fmt.Sprintf("%v-%v-%vT23:59:59", year, int(month), date)

	errFind := r.db.Table("transaction_jetties").Select("COUNT(*) as ritase, SUM(nett_quantity) as tonase").Where("clock_in_date >= ? and clock_in_date <= ? and iupopk_id = ?", startDate, endDate, iupopkId).Scan(&summary).Error

	if errFind != nil {
		return summary, errFind
	}

	return summary, nil
}

func (r *repository) SummaryInventoryStockRom(iupopkId int) (Summary, error) {
	var summary Summary

	var inventory []InventoryStockRom

	errFind := r.db.Preload(clause.Associations).Table("isps").Select("quantity as stock, id as isp_id").Where("iupopk_id = ?", iupopkId).Group("id").Find(&inventory).Error

	if errFind != nil {
		return summary, errFind
	}

	var sumTransactionJetty []SumTransactionJetty
	var countInTransit []CountInTransit

	errFindSum := r.db.Table("transaction_jetties").Select("SUM(nett_quantity) as quantity, isp_id").Where("iupopk_id = ? ", iupopkId).Group("isp_id").Find(&sumTransactionJetty).Error

	if errFindSum != nil {
		return summary, errFindSum
	}

	errFindCount := r.db.Table("transaction_isp_jetties tij").Select("Count(*) as count, ttj.isp_id").Joins("LEFT JOIN transaction_to_jetties ttj on ttj.id = tij.transaction_to_jetty_id").Where("tij.iupopk_id = ? and tij.transaction_jetty_id IS NULL", iupopkId).Group("ttj.isp_id").Find(&countInTransit).Error

	if errFindCount != nil {
		return summary, errFindCount
	}

	var newInventory []InventoryStockRom

	for _, v := range inventory {
		for _, vSum := range sumTransactionJetty {
			if v.IspId == vSum.IspId {
				v.Stock -= vSum.Quantity
				break
			}
		}

		for _, vCount := range countInTransit {
			if v.IspId == vCount.IspId {
				v.CountInTransit = vCount.Count
				break
			}
		}

		newInventory = append(newInventory, v)
	}

	summary.InventoryStockRom = newInventory

	var inventoryStockJetty []InventoryStockJetty

	errFindJetty := r.db.Preload(clause.Associations).Table("jetties").Select("quantity as stock, id as jetty_id").Where("iupopk_id = ?", iupopkId).Group("id").Find(&inventoryStockJetty).Error

	if errFindJetty != nil {
		return summary, errFindJetty
	}

	var sumTransaction []SumTransaction

	errFindSumJetty := r.db.Table("transactions").Select("SUM(quantity) as quantity, loading_port_id as jetty_id").Where("seller_id = ? ", iupopkId).Group("loading_port_id").Find(&sumTransaction).Error

	if errFindSumJetty != nil {
		return summary, errFindSumJetty
	}

	var sumStockJetty []SumTransaction

	errFindSumStockJetty := r.db.Table("transaction_jetties").Select("SUM(nett_quantity) as quantity, jetty_id").Where("iupopk_id = ? ", iupopkId).Group("jetty_id").Find(&sumStockJetty).Error

	if errFindSumStockJetty != nil {
		return summary, errFindSumStockJetty
	}

	var newInventoryStockJetty []InventoryStockJetty

	for _, v := range inventoryStockJetty {
		for _, vSum := range sumTransaction {
			if v.JettyId == vSum.JettyId {
				v.Stock -= vSum.Quantity
				break
			}
		}

		for _, vStock := range sumStockJetty {
			if v.JettyId == vStock.JettyId {
				v.Stock += vStock.Quantity
				break
			}
		}

		newInventoryStockJetty = append(newInventoryStockJetty, v)
	}

	summary.InventoryStockJetty = newInventoryStockJetty

	return summary, nil
}
