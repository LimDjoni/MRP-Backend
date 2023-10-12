package haulingsynchronize

import (
	"ajebackend/model/master/contractor"
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/site"
	"ajebackend/model/master/truck"
	"ajebackend/model/transactionshauling/transactionispjetty"
	"ajebackend/model/transactionshauling/transactionjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"
	"ajebackend/model/transactionshauling/transactiontojetty"
	"ajebackend/model/user"
	"ajebackend/model/useriupopk"
)

type SynchronizeInputMaster struct {
	Contractor []contractor.Contractor `json:"contractor"`
	Isp        []isp.Isp               `json:"isp"`
	Iupopk     []iupopk.Iupopk         `json:"iupopk"`
	Jetty      []jetty.Jetty           `json:"jetty"`
	Pit        []pit.Pit               `json:"pit"`
	Site       []site.Site             `json:"site"`
	Truck      []truck.Truck           `json:"truck"`
	User       []user.User             `json:"user"`
	UserIupopk []useriupopk.UserIupopk `json:"user_iupopk"`
}

type SynchronizeInputTransactionIsp struct {
	TransactionIspJetty []transactionispjetty.TransactionIspJetty `json:"transaction_isp_jetty"`
	TransactionToIsp    []transactiontoisp.TransactionToIsp       `json:"transaction_to_isp"`
	TransactionToJetty  []transactiontojetty.TransactionToJetty   `json:"transaction_to_jetty"`
	SynchronizeTime     string                                    `json:"syncrhonize_time"`
	IupopkId            uint                                      `json:"iupopk_id"`
}

type SynchronizeInputTransactionJetty struct {
	TransactionJetty []transactionjetty.TransactionJetty `json:"transaction_jetty"`
	SynchronizeTime  string                              `json:"syncrhonize_time"`
	IupopkId         uint                                `json:"iupopk_id"`
}
