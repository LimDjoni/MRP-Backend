package haulingsynchronize

import (
	"ajebackend/model/master/contractor"
	"ajebackend/model/master/isp"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
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

	var transactionIspJetties []transactionispjetty.TransactionIspJetty

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

	if len(transactionJetty) > 0 {
		errCreateJetty := tx.Create(&transactionJetty).Error

		if errCreateJetty != nil {
			tx.Rollback()
			return false, errCreateJetty
		}

		for _, v := range transactionJetty {
			var prod production.Production

			if v.IspId == nil {
				errFind := tx.Where("production_date = ? AND pit_id = ? AND isp_id IS NULL AND jetty_id = ?", strings.Split(v.ClockInDate, "T")[0], v.PitId, v.JettyId).First(&prod).Error

				if errFind != nil {
					prod.Quantity = v.NettQuantity
					prod.RitaseQuantity = 1
					prod.IspId = v.IspId
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
				errFind := tx.Where("production_date = ? AND pit_id is NULL AND isp_id = ? AND jetty_id = ?", strings.Split(v.ClockInDate, "T")[0], v.IspId, v.JettyId).First(&prod).Error

				if errFind != nil {
					prod.Quantity = v.NettQuantity
					prod.RitaseQuantity = 1
					prod.IspId = v.IspId
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
			}
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
	where truck_id = %v and isp_id IS NULL and pit_id = %v and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = %v and tj.seam = '%v' and tj.gar = %v ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckId,
					*v.TransactionToJetty.PitId, syncData.IupopkId, v.TransactionToJetty.JettyId, v.TransactionToJetty.Seam, v.TransactionToJetty.Gar)
			}

			if v.TransactionToJetty.IspId != nil {
				rawQuery = fmt.Sprintf(`select tj.* from transaction_jetties tj
	LEFT JOIN transaction_isp_jetties tij on tij.transaction_jetty_id = tj.id
	where truck_id = %v and isp_id = %v and pit_id IS NULL and tj.iupopk_id = %v and tij.id IS NULL and tj.jetty_id = %v ORDER BY tj.created_at asc`, v.TransactionToJetty.TruckId,
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

	if len(syncData.Truck) > 0 {
		for _, v := range syncData.Truck {
			var tempTruck truck.Truck

			errFind := tx.Where("id = ?", v.ID).First(&tempTruck).Error

			if errFind == nil {
				tempTruck = v
				if v.Rfid == nil || *v.Rfid == "" {
					tempTruck.Rfid = nil
				}
				errUpd := tx.Save(&tempTruck).Error

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

	if len(syncData.Contractor) > 0 {
		for _, v := range syncData.Contractor {
			var tempContractor contractor.Contractor

			errFind := tx.Where("id = ?", v.ID).First(&tempContractor).Error

			if errFind == nil {
				tempContractor = v

				errUpd := tx.Save(&tempContractor).Error

				if errUpd != nil {
					tx.Rollback()
					return false, errUpd
				}
			} else {
				errCreate := tx.Create(&v).Error

				if errCreate != nil {
					tx.Rollback()
					return false, errCreate
				}
			}
		}
	}

	if len(syncData.Pit) > 0 {
		for _, v := range syncData.Pit {
			var tempPit pit.Pit

			errFind := tx.Where("id = ?", v.ID).First(&tempPit).Error

			if errFind == nil {
				tempPit = v

				errUpd := tx.Save(&tempPit).Error

				if errUpd != nil {
					tx.Rollback()
					return false, errUpd
				}
			} else {
				errCreate := tx.Create(&v).Error

				if errCreate != nil {
					tx.Rollback()
					return false, errCreate
				}
			}
		}
	}

	if len(syncData.Isp) > 0 {
		for _, v := range syncData.Isp {
			var tempIsp isp.Isp

			errFind := tx.Where("id = ?", v.ID).First(&tempIsp).Error

			if errFind == nil {
				tempIsp = v

				errUpd := tx.Save(&tempIsp).Error

				if errUpd != nil {
					tx.Rollback()
					return false, errUpd
				}
			} else {
				errCreate := tx.Create(&v).Error

				if errCreate != nil {
					tx.Rollback()
					return false, errCreate
				}
			}
		}
	}

	if len(syncData.Jetty) > 0 {
		for _, v := range syncData.Jetty {
			var tempJetty jetty.Jetty

			errFind := tx.Where("id = ?", v.ID).First(&tempJetty).Error

			if errFind == nil {
				tempJetty = v

				errUpd := tx.Save(&tempJetty).Error

				if errUpd != nil {
					tx.Rollback()
					return false, errUpd
				}
			} else {
				errCreate := tx.Create(&v).Error

				if errCreate != nil {
					tx.Rollback()
					return false, errCreate
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
