package haulingsynchronize

import (
	"ajebackend/model/transactionshauling/transactionispjetty"
	"ajebackend/model/transactionshauling/transactionjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"
	"ajebackend/model/transactionshauling/transactiontojetty"
	"fmt"
	"strings"

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

	var inputTransactionToIsp []transactiontoisp.InputTransactionToIsp
	var inputTransactionToJetty []transactiontojetty.InputTransactionToJetty

	inputTransactionToIsp = syncData.TransactionToIsp
	inputTransactionToJetty = syncData.TransactionToJetty

	tx := r.db.Begin()

	if len(inputTransactionToIsp) > 0 {
		errCreateToIsp := tx.Model(&transactionToIsp).Create(&inputTransactionToIsp).Error

		if errCreateToIsp != nil {
			tx.Rollback()
			return false, errCreateToIsp
		}
	}

	if len(inputTransactionToJetty) > 0 {
		errCreateToJetty := tx.Model(&transactionToJetty).Create(&inputTransactionToJetty).Error

		if errCreateToJetty != nil {
			tx.Rollback()
			return false, errCreateToJetty
		}

	}

	var transactionIspJetties []map[string]interface{}

	if len(transactionToJetty) > 0 {
		for _, v := range transactionToJetty {
			splitId := strings.Split(v.IdNumber, "PHU-")

			temp := make(map[string]interface{})
			temp["transaction_to_jetties"] = v.ID
			temp["iupopk_id"] = syncData.IupopkId
			temp["id_number"] = "HAU-" + splitId[1]

			transactionIspJetties = append(transactionIspJetties, temp)
		}
	}

	if len(transactionIspJetties) > 0 {
		errCreateIspJetty := tx.Model(&transactionIspJetty).Create(&transactionIspJetties).Error

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

	var inputTransactionJetty []transactionjetty.InputTransactionJetty

	inputTransactionJetty = syncData.TransactionJetty

	tx := r.db.Begin()

	if len(inputTransactionJetty) > 0 {
		errCreateJetty := tx.Model(&transactionJetty).Create(&inputTransactionJetty).Error

		if errCreateJetty != nil {
			tx.Rollback()
			return false, errCreateJetty
		}
	}

	var transactionIspJetty []transactionispjetty.TransactionIspJetty

	errFindIspJetty := tx.Preload(clause.Associations).Where("transaction_jetty_id IS NULL").Order("created_at asc").Find(&transactionIspJetty).Error

	if errFindIspJetty != nil {
		tx.Rollback()
		return false, errFindIspJetty
	}

	if len(transactionIspJetty) > 0 {
		for _, v := range transactionIspJetty {
			var tempTransactionJetty transactionjetty.TransactionJetty

			var rawQuery string

			if v.TransactionToJetty.PitId == nil && v.TransactionToJetty.IspId == nil {
				continue
			}

			if v.TransactionToJetty.PitId != nil {
				rawQuery = fmt.Sprintf(`select tj.* from transaction_jetties tj
	LEFT JOIN transaction_isp_jetties tij on tij.transaction_jetty_id = tj.id
	where truck_id = %v and isp_id IS NULL and pit_id = %v and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = %v and tj.seam = '%v' ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckId,
					*v.TransactionToJetty.PitId, syncData.IupopkId, v.TransactionToJetty.JettyId, v.TransactionToJetty.Seam)
			}

			if v.TransactionToJetty.IspId != nil {
				rawQuery = fmt.Sprintf(`select tj.* from transaction_jetties tj
	LEFT JOIN transaction_isp_jetties tij on tij.transaction_jetty_id = tj.id
	where truck_id = %v and isp_id = %v and pit_id IS NULL and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = %v and tj.seam = '%v' ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckId,
					*v.TransactionToJetty.IspId, syncData.IupopkId, v.TransactionToJetty.JettyId, v.TransactionToJetty.Seam)
			}

			errFindTransactionJetty := tx.Raw(rawQuery).First(&tempTransactionJetty).Error

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
