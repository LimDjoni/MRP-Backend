package transactionshauling

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/jetty"
)

type SummaryJettyTransactionPerDay struct {
	Ritase int     `json:"ritase"`
	Tonase float64 `json:"tonase"`
}

type InventoryStockRom struct {
	IspId uint    `json:"isp_id"`
	Isp   isp.Isp `json:"isp"`
	Stock float64 `json:"stock"`
}

type SumTransactionJetty struct {
	IspId    uint    `json:"isp_id"`
	Quantity float64 `json:"quantity"`
}

type InventoryStockJetty struct {
	JettyId uint        `json:"jetty_id"`
	Jetty   jetty.Jetty `json:"jetty"`
	Stock   float64     `json:"stock"`
}

type SumTransaction struct {
	JettyId  uint    `json:"jetty_id"`
	Quantity float64 `json:"quantity"`
}

type Summary struct {
	InventoryStockRom   []InventoryStockRom   `json:"inventory_stock_rom"`
	InventoryStockJetty []InventoryStockJetty `json:"inventory_stock_jetty"`
}
