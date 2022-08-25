package logs

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Logs struct {
	gorm.Model
	Input datatypes.JSON `json:"input"`
	Message	datatypes.JSON `json:"message"`
	TransactionId *uint `json:"transaction_id"`
	Transaction *transaction.Transaction `json:"transaction"`
	MinerbaId *uint `json:"minerba_id"`
	Minerba *minerba.Minerba `json:"minerba"`
	DmoId *uint `json:"dmo_id"`
	Dmo *dmo.Dmo `json:"dmo"`
}
