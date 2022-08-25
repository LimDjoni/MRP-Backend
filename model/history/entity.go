package history

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	DmoId *uint `json:"dmo_id"`
	Dmo *dmo.Dmo
	TransactionId *uint `json:"transaction_id"`
	Transaction *transaction.Transaction
	UserId uint `json:"user_id"`
	User user.User
	MinerbaId *uint `json:"minerba_id"`
	Minerba *minerba.Minerba
	Status string `json:"status"`
	BeforeData datatypes.JSON `json:"before_data"`
	AfterData datatypes.JSON `json:"after_data"`
}
