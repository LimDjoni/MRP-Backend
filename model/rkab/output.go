package rkab

type DetailRkab struct {
	ListRkab        []Rkab  `json:"list_rkab"`
	TotalProduction float64 `json:"total_production"`
	TotalSales      float64 `json:"total_sales"`
}

type RkabProductionQuantity struct {
	TotalProduction float64 `json:"total_production"`
}

type RkabSalesQuantity struct {
	TotalSales float64 `json:"total_sales"`
}
