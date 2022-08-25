package minerbatransaction

import (
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"
	"gorm.io/gorm"
)

type MinerbaTransaction struct {
	gorm.Model
	MinerbaId uint `json:"minerba_id"`
	Minerba minerba.Minerba
	TransactionId uint `json:"transaction_id"`
	Transaction transaction.Transaction
}
