package rkab

type DetailRkab struct {
	ListRkab         []Rkab  `json:"list_rkab"`
	TotalProduction  float64 `json:"total_production"`
	TotalSales       float64 `json:"total_sales"`
	TotalProduction2 float64 `json:"total_production_2"`
	TotalSales2      float64 `json:"total_sales_2"`
	TotalProduction3 float64 `json:"total_production_3"`
	TotalSales3      float64 `json:"total_sales_3"`
}

type RkabProductionQuantity struct {
	TotalProduction float64 `json:"total_production"`
}

type RkabSalesQuantity struct {
	TotalSales float64 `json:"total_sales"`
}
