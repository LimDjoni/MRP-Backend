package adjuststock

type RegisterAdjustStockInput struct {
	Date  string  `json:"date"`
	Stock float64 `json:"stock"`
}

type SortFilterAdjustStock struct {
	Field string
	Sort  string
	Date  string
	Stock string
}
