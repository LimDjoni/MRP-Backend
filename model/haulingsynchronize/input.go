package haulingsynchronize

import (
	"ajebackend/model/master/contractor"
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/truck"
	"ajebackend/model/transactionshauling/transactionjetty"
	"ajebackend/model/transactionshauling/transactiontoisp"
	"ajebackend/model/transactionshauling/transactiontojetty"
)

type SynchronizeInputMaster struct {
	Contractor []contractor.Contractor `json:"contractor"`
	Isp        []isp.Isp               `json:"isp"`
	Iupopk     []iupopk.Iupopk         `json:"iupopk"`
	Jetty      []jetty.Jetty           `json:"jetty"`
	Pit        []pit.Pit               `json:"pit"`
	Truck      []truck.Truck           `json:"truck"`
}

type SynchronizeInputTransactionIsp struct {
	TransactionToIsp   []transactiontoisp.TransactionToIsp     `json:"transaction_to_isp"`
	TransactionToJetty []transactiontojetty.TransactionToJetty `json:"transaction_to_jetty"`
	SynchronizeTime    string                                  `json:"syncrhonize_time"`
	IupopkId           uint                                    `json:"iupopk_id"`
}

type SynchronizeInputTransactionJetty struct {
	TransactionJetty []transactionjetty.TransactionJetty `json:"transaction_jetty"`
	SynchronizeTime  string                              `json:"syncrhonize_time"`
	IupopkId         uint                                `json:"iupopk_id"`
	Contractor       []contractor.Contractor             `json:"contractor"`
	Isp              []isp.Isp                           `json:"isp"`
	Jetty            []jetty.Jetty                       `json:"jetty"`
	Pit              []pit.Pit                           `json:"pit"`
	Truck            []truck.Truck                       `json:"truck"`
}
