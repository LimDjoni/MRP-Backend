package haulingsynchronize

import (
	"ajebackend/model/master/truck"
	"ajebackend/model/production"
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
	UpdateSyncMasterIsp(iupopkId uint, dateTime string) (bool, error)
	UpdateSyncMasterJetty(iupopkId uint, dateTime string) (bool, error)
	GetSyncMasterDataIsp(iupopkId uint) (MasterDataIsp, error)
	GetSyncMasterDataJetty(iupopkId uint) (MasterDataJetty, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error) {
	var transactionToIsp []transactiontoisp.TransactionToIsp
	var transactionToJetty []transactiontojetty.TransactionToJetty

	transactionToIsp = syncData.TransactionToIsp
	transactionToJetty = syncData.TransactionToJetty

	tx := r.db.Begin()

	// Create Transaction Transaction To Isp (Stock Rom)
	if len(transactionToIsp) > 0 {
		errCreateToIsp := tx.Create(&transactionToIsp).Error

		if errCreateToIsp != nil {
			tx.Rollback()
			return false, errCreateToIsp
		}
	}

	// Create Transaction To Jetty
	if len(transactionToJetty) > 0 {
		errCreateToJetty := tx.Create(&transactionToJetty).Error

		if errCreateToJetty != nil {
			tx.Rollback()
			return false, errCreateToJetty
		}
	}

	var transactionIspJetties []transactionispjetty.TransactionIspJetty

	// Create Data Transaction Isp Jetties
	if len(transactionToJetty) > 0 {
		for _, v := range transactionToJetty {
			splitId := strings.Split(v.IdNumber, "PHU-")

			var temp transactionispjetty.TransactionIspJetty
			temp.TransactionToJettyId = v.ID
			temp.IupopkId = syncData.IupopkId
			temp.IdNumber = "HAU-" + splitId[1]

			transactionIspJetties = append(transactionIspJetties, temp)
		}
	}

	// Create Transaction Isp Jetties Database
	if len(transactionIspJetties) > 0 {
		errCreateIspJetty := tx.Create(&transactionIspJetties).Error

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

	// Create / Update Truck Database
	if len(syncData.Truck) > 0 {
		for _, v := range syncData.Truck {
			var tempTruck truck.Truck

			errFind := tx.Where("code = ?", v.Code).First(&tempTruck).Error

			if errFind == nil {
				tempTruck = v
				if v.Rfid == nil || *v.Rfid == "" {
					tempTruck.Rfid = nil
				}
				errUpd := tx.Where("code = ?", v.Code).Updates(&tempTruck).Error

				if errUpd != nil {
					tx.Rollback()
					return false, errUpd
				}
			} else {
				if v.Rfid == nil || *v.Rfid == "" {
					v.Rfid = nil
				}
				errCreate := tx.Create(&v).Error

				if errCreate != nil {
					tx.Rollback()
					return false, errCreate
				}
			}
		}
	}

	// Create Transaction Jetty
	if len(transactionJetty) > 0 {
		errCreateJetty := tx.Create(&transactionJetty).Error

		if errCreateJetty != nil {
			tx.Rollback()
			return false, errCreateJetty
		}

		for _, v := range transactionJetty {
			var prod production.Production

			if v.IspId == nil {
				errFind := tx.Where("production_date = ? AND pit_id = ? AND isp_id IS NULL AND jetty_id = ?", strings.Split(v.ClockInDate, "T")[0], v.PitId, &v.JettyId).First(&prod).Error

				if errFind != nil {
					prod.Quantity = v.NettQuantity
					prod.RitaseQuantity = 1
					prod.PitId = v.PitId
					prod.JettyId = &v.JettyId
					prod.IupopkId = syncData.IupopkId
					prod.ProductionDate = strings.Split(v.ClockInDate, "T")[0]

					errCreateProd := tx.Create(&prod).Error

					if errCreateProd != nil {
						tx.Rollback()
						return false, errCreateProd
					}
				} else {
					errUpdProd := tx.Table("productions").Where("id = ?", prod.ID).Updates(map[string]interface{}{"quantity": prod.Quantity + v.NettQuantity, "ritase_quantity": prod.RitaseQuantity + 1}).Error

					if errUpdProd != nil {
						tx.Rollback()
						return false, errUpdProd
					}
				}
			} else if v.PitId == nil {
				errFind := tx.Where("production_date = ? AND pit_id is NULL AND isp_id = ? AND jetty_id = ?", strings.Split(v.ClockInDate, "T")[0], v.IspId, &v.JettyId).First(&prod).Error

				if errFind != nil {
					prod.Quantity = v.NettQuantity
					prod.RitaseQuantity = 1
					prod.IspId = v.IspId
					prod.JettyId = &v.JettyId
					prod.IupopkId = syncData.IupopkId
					prod.ProductionDate = strings.Split(v.ClockInDate, "T")[0]

					errCreateProd := tx.Create(&prod).Error

					if errCreateProd != nil {
						tx.Rollback()
						return false, errCreateProd
					}
				} else {
					errUpdProd := tx.Table("productions").Where("id = ?", prod.ID).Updates(map[string]interface{}{"quantity": prod.Quantity + v.NettQuantity, "ritase_quantity": prod.RitaseQuantity + 1}).Error

					if errUpdProd != nil {
						tx.Rollback()
						return false, errUpdProd
					}
				}
			}
		}
	}

	var transactionIspJetty []transactionispjetty.TransactionIspJetty

	tx.Preload(clause.Associations).Where("transaction_jetty_id IS NULL").Order("created_at asc").Find(&transactionIspJetty)

	// Connect Transaction to Jetty & Transaction Jetty in Transaction Isp Jetty

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
	where truck_code = %v and isp_id IS NULL and pit_id = '%v' and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = '%v' and tj.seam = '%v' and tj.gar = %v ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckCode,
					*v.TransactionToJetty.PitId, syncData.IupopkId, v.TransactionToJetty.JettyId, v.TransactionToJetty.Seam, v.TransactionToJetty.Gar)
			}

			if v.TransactionToJetty.IspId != nil {
				rawQuery = fmt.Sprintf(`select tj.* from transaction_jetties tj
	LEFT JOIN transaction_isp_jetties tij on tij.transaction_jetty_id = tj.id
	where truck_code = %v and isp_id = '%v' and pit_id IS NULL and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = '%v' ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckCode,
					*v.TransactionToJetty.IspId, syncData.IupopkId, v.TransactionToJetty.JettyId)
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

func (r *repository) UpdateSyncMasterIsp(iupopkId uint, dateTime string) (bool, error) {
	var haulingSync HaulingSynchronize

	errFirst := r.db.Where("iupopk_id = ?", iupopkId).First(&haulingSync).Error

	if errFirst != nil {
		return false, errFirst
	}

	// Update latest time synchronize master to ISP
	errUpdate := r.db.Model(&haulingSync).Update("last_synchronize_master_to_isp", dateTime).Error

	if errUpdate != nil {
		return false, errUpdate
	}

	return true, nil
}

func (r *repository) UpdateSyncMasterJetty(iupopkId uint, dateTime string) (bool, error) {
	var haulingSync HaulingSynchronize

	errFirst := r.db.Where("iupopk_id = ?", iupopkId).First(&haulingSync).Error

	if errFirst != nil {
		return false, errFirst
	}

	// Update latest time synchronize master to Jetty

	errUpdate := r.db.Model(&haulingSync).Update("last_synchronize_master_to_jetty", dateTime).Error

	if errUpdate != nil {
		return false, errUpdate
	}

	return true, nil
}

func (r *repository) GetSyncMasterDataIsp(iupopkId uint) (MasterDataIsp, error) {
	var masterDataIsp MasterDataIsp

	var haulingSync HaulingSynchronize

	errFirst := r.db.Where("iupopk_id = ?", iupopkId).First(&haulingSync).Error

	if errFirst != nil {
		return masterDataIsp, errFirst
	}

	if haulingSync.LastSynchronizeMasterToIsp != nil {
		errFindContractor := r.db.Table("contractors").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Contractor).Error

		if errFindContractor != nil {
			return masterDataIsp, errFindContractor
		}

		errFindIsp := r.db.Table("isps").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Isp).Error

		if errFindIsp != nil {
			return masterDataIsp, errFindIsp
		}

		errFindIupopk := r.db.Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Iupopk).Error

		if errFindIupopk != nil {
			return masterDataIsp, errFindIupopk
		}

		errFindJetty := r.db.Table("jetties").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Jetty).Error

		if errFindJetty != nil {
			return masterDataIsp, errFindJetty
		}

		errFindPit := r.db.Table("pits").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Pit).Error

		if errFindPit != nil {
			return masterDataIsp, errFindPit
		}

		errFindTruck := r.db.Table("trucks").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Truck).Error

		if errFindTruck != nil {
			return masterDataIsp, errFindTruck
		}

		errFindRole := r.db.Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.Role).Error

		if errFindRole != nil {
			return masterDataIsp, errFindRole
		}

		errFindUser := r.db.Where("updated_at >= ? AND is_ho = false", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.User).Error

		if errFindUser != nil {
			return masterDataIsp, errFindUser
		}

		errFindUserIupopk := r.db.Table("user_iupopks").Joins("left join users u on u.id = user_iupopks.user_id").Where("user_iupopks.updated_at >= ? AND u.is_ho = false", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.UserIupopk).Error

		if errFindUserIupopk != nil {
			return masterDataIsp, errFindUserIupopk
		}

		errFindUserRole := r.db.Table("user_roles").Joins("left join users u on u.id = user_roles.user_id").Where("user_roles.updated_at >= ? AND u.is_ho = false", haulingSync.LastSynchronizeMasterToIsp).Find(&masterDataIsp.UserRole).Error

		if errFindUserRole != nil {
			return masterDataIsp, errFindUserRole
		}
	} else {
		errFindContractor := r.db.Table("contractors").Find(&masterDataIsp.Contractor).Error

		if errFindContractor != nil {
			return masterDataIsp, errFindContractor
		}

		errFindIsp := r.db.Table("isps").Find(&masterDataIsp.Isp).Error

		if errFindIsp != nil {
			return masterDataIsp, errFindIsp
		}

		errFindIupopk := r.db.Find(&masterDataIsp.Iupopk).Error

		if errFindIupopk != nil {
			return masterDataIsp, errFindIupopk
		}

		errFindJetty := r.db.Table("jetties").Find(&masterDataIsp.Jetty).Error

		if errFindJetty != nil {
			return masterDataIsp, errFindJetty
		}

		errFindPit := r.db.Table("pits").Find(&masterDataIsp.Pit).Error

		if errFindPit != nil {
			return masterDataIsp, errFindPit
		}

		errFindTruck := r.db.Table("trucks").Find(&masterDataIsp.Truck).Error

		if errFindTruck != nil {
			return masterDataIsp, errFindTruck
		}

		errFindRole := r.db.Find(&masterDataIsp.Role).Error

		if errFindRole != nil {
			return masterDataIsp, errFindRole
		}

		errFindUser := r.db.Where("is_ho = false").Find(&masterDataIsp.User).Error

		if errFindUser != nil {
			return masterDataIsp, errFindUser
		}

		errFindUserIupopk := r.db.Table("user_iupopks").Joins("left join users u on u.id = user_iupopks.user_id").Where("u.is_ho = false").Find(&masterDataIsp.UserIupopk).Error

		if errFindUserIupopk != nil {
			return masterDataIsp, errFindUserIupopk
		}

		errFindUserRole := r.db.Table("user_roles").Joins("left join users u on u.id = user_roles.user_id").Where("u.is_ho = false").Find(&masterDataIsp.UserRole).Error

		if errFindUserRole != nil {
			return masterDataIsp, errFindUserRole
		}
	}

	return masterDataIsp, nil
}

func (r *repository) GetSyncMasterDataJetty(iupopkId uint) (MasterDataJetty, error) {
	var masterDataJetty MasterDataJetty

	var haulingSync HaulingSynchronize

	errFirst := r.db.Where("iupopk_id = ?", iupopkId).First(&haulingSync).Error

	if errFirst != nil {
		return masterDataJetty, errFirst
	}

	if haulingSync.LastSynchronizeMasterToJetty != nil {
		errFindContractor := r.db.Table("contractors").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Contractor).Error

		if errFindContractor != nil {
			return masterDataJetty, errFindContractor
		}

		errFindIsp := r.db.Table("isps").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Isp).Error

		if errFindIsp != nil {
			return masterDataJetty, errFindIsp
		}

		errFindIupopk := r.db.Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Iupopk).Error

		if errFindIupopk != nil {
			return masterDataJetty, errFindIupopk
		}

		errFindJetty := r.db.Table("jetties").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Jetty).Error

		if errFindJetty != nil {
			return masterDataJetty, errFindJetty
		}

		errFindPit := r.db.Table("pits").Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Pit).Error

		if errFindPit != nil {
			return masterDataJetty, errFindPit
		}

		errFindRole := r.db.Where("updated_at >= ?", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.Role).Error

		if errFindRole != nil {
			return masterDataJetty, errFindRole
		}

		errFindUser := r.db.Where("updated_at >= ? AND is_ho = false", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.User).Error

		if errFindUser != nil {
			return masterDataJetty, errFindUser
		}

		errFindUserIupopk := r.db.Table("user_iupopks").Joins("left join users u on u.id = user_iupopks.user_id").Where("user_iupopks.updated_at >= ? AND u.is_ho = false", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.UserIupopk).Error

		if errFindUserIupopk != nil {
			return masterDataJetty, errFindUserIupopk
		}

		errFindUserRole := r.db.Table("user_roles").Joins("left join users u on u.id = user_roles.user_id").Where("user_roles.updated_at >= ? AND u.is_ho = false", haulingSync.LastSynchronizeMasterToJetty).Find(&masterDataJetty.UserRole).Error

		if errFindUserRole != nil {
			return masterDataJetty, errFindUserRole
		}
	} else {
		errFindContractor := r.db.Table("contractors").Find(&masterDataJetty.Contractor).Error

		if errFindContractor != nil {
			return masterDataJetty, errFindContractor
		}

		errFindIsp := r.db.Table("isps").Find(&masterDataJetty.Isp).Error

		if errFindIsp != nil {
			return masterDataJetty, errFindIsp
		}

		errFindIupopk := r.db.Find(&masterDataJetty.Iupopk).Error

		if errFindIupopk != nil {
			return masterDataJetty, errFindIupopk
		}

		errFindJetty := r.db.Table("jetties").Find(&masterDataJetty.Jetty).Error

		if errFindJetty != nil {
			return masterDataJetty, errFindJetty
		}

		errFindPit := r.db.Table("pits").Find(&masterDataJetty.Pit).Error

		if errFindPit != nil {
			return masterDataJetty, errFindPit
		}

		errFindRole := r.db.Find(&masterDataJetty.Role).Error

		if errFindRole != nil {
			return masterDataJetty, errFindRole
		}

		errFindUser := r.db.Where("is_ho = false").Find(&masterDataJetty.User).Error

		if errFindUser != nil {
			return masterDataJetty, errFindUser
		}

		errFindUserIupopk := r.db.Table("user_iupopks").Joins("left join users u on u.id = user_iupopks.user_id").Where("u.is_ho = false").Find(&masterDataJetty.UserIupopk).Error

		if errFindUserIupopk != nil {
			return masterDataJetty, errFindUserIupopk
		}

		errFindUserRole := r.db.Table("user_roles").Joins("left join users u on u.id = user_roles.user_id").Where("u.is_ho = false").Find(&masterDataJetty.UserRole).Error

		if errFindUserRole != nil {
			return masterDataJetty, errFindUserRole
		}
	}

	return masterDataJetty, nil
}
