package seeding

import (
	"ajebackend/model/transaction"

	"gorm.io/gorm"
)

func UpdateTransactionsQuantity(db *gorm.DB) {
	var transactionsDn transaction.Transaction

	db.Raw("UPDATE transactions SET quantity_unloading = quantity WHERE quantity_unloading IS NULL").Scan(&transactionsDn)

	return
}
