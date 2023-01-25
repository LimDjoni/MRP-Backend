package seeding

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"

	"gorm.io/gorm"
)

func UpdateNaming(db *gorm.DB) {
	var transactionsDn transaction.Transaction

	var minerbas minerba.Minerba
	var dmos dmo.Dmo

	db.Raw("UPDATE transactions SET id_number = REPLACE(id_number, 'DN', 'TDN') WHERE id_number LIKE 'DN%'").Scan(&transactionsDn)

	db.Raw("UPDATE minerbas SET id_number = REPLACE(id_number, 'LM', 'LSD') WHERE id_number LIKE 'LM%'").Scan(&minerbas)

	db.Raw("UPDATE dmos SET id_number = REPLACE(id_number, 'SR', 'LDO') WHERE id_number LIKE 'SR%'").Scan(&dmos)
	return
}
