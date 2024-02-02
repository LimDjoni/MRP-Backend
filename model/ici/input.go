package ici

type InputCreateUpdateIci struct {
	Date       string  `json:"date"`
	Level      string  `json:"level"`
	Avarage    float64 `json:"average"`
	UnitPrice  float64 `json:"unit_price"`
	Currency   string  `json:"currency"`
	IciLevelId uint    `json:"ici_level_id"`
}
