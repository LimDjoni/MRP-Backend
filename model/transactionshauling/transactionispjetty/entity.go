package transactionispjetty

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/transactionshauling/transactionjetty"
	"ajebackend/model/transactionshauling/transactiontojetty"

	"gorm.io/gorm"
)

type TransactionIspJetty struct {
	gorm.Model
	IdNumber             string                                `json:"id_number"`
	TransactionJettyId   *uint                                 `json:"transaction_jetty_id"`
	TransactionJetty     *transactionjetty.TransactionJetty    `json:"transaction_jetty"`
	TransactionToJettyId uint                                  `json:"transaction_to_jetty_id"`
	TransactionToJetty   transactiontojetty.TransactionToJetty `json:"transaction_to_jetty"`
	IupopkId             uint                                  `json:"iupopk_id"`
	Iupopk               iupopk.Iupopk                         `json:"iupopk"`
}
