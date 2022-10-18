package seeding

import (
	"ajebackend/model/transaction"
	"gorm.io/gorm"
)

func UpdateTransactionsRoyalty(db *gorm.DB) {
	db.Unscoped().Model(&transaction.Transaction{}).Where("dp_royalty_ntpn = '' AND payment_dp_royalty_ntpn = '' AND dp_royalty_billing_code = '' AND payment_dp_royalty_billing_code = ''").Updates(map[string]interface{}{
		"dp_royalty_ntpn" : nil,
		"dp_royalty_billing_code" : nil,
		"payment_dp_royalty_ntpn" : nil,
		"payment_dp_royalty_billing_code" : nil,
	})

	db.Unscoped().Model(&transaction.Transaction{}).Where("is_not_claim is NULL").Update("is_not_claim", false)

	return
}
