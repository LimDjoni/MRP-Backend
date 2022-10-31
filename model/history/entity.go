package history

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/production"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	DmoId *uint `json:"dmo_id"`
	Dmo *dmo.Dmo `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	TransactionId *uint `json:"transaction_id"`
	Transaction *transaction.Transaction `json:"transaction" gorm:"constraint:OnDelete:CASCADE;"`
	UserId uint `json:"user_id"`
	User user.User
	MinerbaId *uint `json:"minerba_id"`
	Minerba *minerba.Minerba `json:"minerba" gorm:"constraint:OnDelete:CASCADE;"`
	ProductionId *uint `json:"production_id"`
	Production production.Production `json:"production" gorm:"constraint:OnDelete:CASCADE;"`
	Status string `json:"status"`
	BeforeData datatypes.JSON `json:"before_data"`
	AfterData datatypes.JSON `json:"after_data"`
}
