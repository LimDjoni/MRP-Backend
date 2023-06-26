package rkab

type DetailRkab struct {
	ListRkab        []Rkab  `json:"list_rkab"`
	TotalProduction float64 `json:"total_production"`
}

type RkabProductionQuantity struct {
	TotalProduction float64 `json:"total_production"`
}
