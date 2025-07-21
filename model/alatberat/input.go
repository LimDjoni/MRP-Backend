package alatberat

type RegisterAlatBeratInput struct {
	BrandId          uint    `json:"brand_id" validate:"required"`
	HeavyEquipmentId uint    `json:"heavy_equipment_id" validate:"required"`
	SeriesId         uint    `json:"series_id" validate:"required"`
	Consumption      float64 `json:"consumption" validate:"required"`
	Tolerance        uint    `json:"tolerance" validate:"required"`
}

type SortFilterAlatBerat struct {
	Field            string
	Sort             string
	BrandId          string
	HeavyEquipmentId string
	SeriesId         string
	Consumption      string
	Tolerance        string
}
