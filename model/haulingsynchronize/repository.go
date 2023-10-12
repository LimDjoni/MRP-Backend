package haulingsynchronize

import (
	"ajebackend/model/transactionshauling/transactionispjetty"
	"ajebackend/model/transactionshauling/transactionjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"
	"ajebackend/model/transactionshauling/transactiontojetty"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error)
	SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error) {
	var transactionIspJetty []transactionispjetty.TransactionIspJetty
	var transactionToIsp []transactiontoisp.TransactionToIsp
	var transactionToJetty []transactiontojetty.TransactionToJetty

	transactionIspJetty = syncData.TransactionIspJetty
	transactionToIsp = syncData.TransactionToIsp
	transactionToJetty = syncData.TransactionToJetty

	tx := r.db.Begin()

	if len(transactionToIsp) > 0 {
		errCreateToIsp := tx.Create(&transactionToIsp).Error

		if errCreateToIsp != nil {
			tx.Rollback()
			return false, errCreateToIsp
		}
	}

	if len(transactionToJetty) > 0 {
		errCreateToJetty := tx.Create(&transactionToJetty).Error

		if errCreateToJetty != nil {
			tx.Rollback()
			return false, errCreateToJetty
		}
	}

	if len(transactionIspJetty) > 0 {
		errCreateIspJetty := tx.Create(&transactionIspJetty).Error

		if errCreateIspJetty != nil {
			tx.Rollback()
			return false, errCreateIspJetty
		}
	}

	var haulingSync HaulingSynchronize

	errFindSynchronize := tx.Where("iupopk_id = ?", syncData.IupopkId).First(&haulingSync).Error

	if errFindSynchronize != nil {
		tx.Rollback()
		return false, errFindSynchronize
	}

	errUpdSynchronize := tx.Table("hauling_synchronizes").Where("id = ?", haulingSync.ID).Update("last_synchronize_isp", syncData.SynchronizeTime).Error

	if errUpdSynchronize != nil {
		tx.Rollback()
		return false, errUpdSynchronize
	}

	tx.Commit()
	return true, nil
}

func (r *repository) SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error) {

	var transactionJetty []transactionjetty.TransactionJetty

	transactionJetty = syncData.TransactionJetty

	tx := r.db.Begin()

	if len(transactionJetty) > 0 {
		errCreateJetty := tx.Create(&transactionJetty).Error

		if errCreateJetty != nil {
			tx.Rollback()
			return false, errCreateJetty
		}

		var transactionIspJetty []transactionispjetty.TransactionIspJetty

		errFindIspJetty := tx.Preload(clause.Associations).Where("transaction_jetty_id IS NULL").Order("created_at asc").Find(&transactionIspJetty).Error

		if errFindIspJetty != nil {
			tx.Rollback()
			return false, errFindIspJetty
		}

		for _, v := range transactionIspJetty {
			var tempTransactionJetty transactionjetty.TransactionJetty

			errFindTransactionJetty := tx.Where("truck_id = ? and isp_id = ? and site_id = ?", v.TransactionToJetty.TruckId,
				v.TransactionToJetty.IspId,
				v.TransactionToJetty.SiteId).Order("created_at asc").First(&tempTransactionJetty).Error

			if errFindTransactionJetty == nil {
				errUpdIspJetty := tx.Table("transaction_isp_jetties").Where("id = ?", v.ID).Update("transaction_jetty_id", tempTransactionJetty.ID).Error

				if errUpdIspJetty != nil {
					tx.Rollback()
					return false, errUpdIspJetty
				}
			}
		}
	}

	var haulingSync HaulingSynchronize

	errFindSynchronize := tx.Where("iupopk_id = ?", syncData.IupopkId).First(&haulingSync).Error

	if errFindSynchronize != nil {
		tx.Rollback()
		return false, errFindSynchronize
	}

	errUpdSynchronize := tx.Table("hauling_synchronizes").Where("id = ?", haulingSync.ID).Update("last_synchronize_jetty", syncData.SynchronizeTime).Error

	if errUpdSynchronize != nil {
		tx.Rollback()
		return false, errUpdSynchronize
	}

	tx.Commit()
	return true, nil
}
