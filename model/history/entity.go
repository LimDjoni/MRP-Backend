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
	DmoId *uint `json:"dmo_id" gorm:"constraint:OnDelete:CASCADE;"`
	Dmo *dmo.Dmo
	TransactionId *uint `json:"transaction_id" gorm:"constraint:OnDelete:CASCADE;"`
	Transaction *transaction.Transaction
	UserId uint `json:"user_id"`
	User user.User
	MinerbaId *uint `json:"minerba_id" gorm:"constraint:OnDelete:CASCADE;"`
	Minerba *minerba.Minerba
	Status string `json:"status"`
	BeforeData datatypes.JSON `json:"before_data"`
	AfterData datatypes.JSON `json:"after_data"`
}
