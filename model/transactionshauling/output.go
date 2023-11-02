package transactionshauling

import (
	"ajebackend/model/master/isp"
)

type SummaryJettyTransactionPerDay struct {
	Ritase int     `json:"ritase"`
	Tonase float64 `json:"tonase"`
}

type InventoryStockRom struct {
	IspId          uint    `json:"isp_id"`
	Isp            isp.Isp `json:"isp"`
	Stock          float64 `json:"stock"`
	CountInTransit int     `json:"count_in_transit"`
}

type SumTransactionJetty struct {
	IspId    uint    `json:"isp_id"`
	Quantity float64 `json:"quantity"`
}

type CountInTransit struct {
	IspId uint `json:"isp_id"`
	Count int  `json:"count"`
}
