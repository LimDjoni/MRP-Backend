package unit

type RegisterUnitInput struct {
	UnitName         string `json:"unit_name" validate:"required"`
	BrandId          uint   `json:"brand_id" validate:"required"`
	HeavyEquipmentId uint   `json:"heavy_equipment_id" validate:"required"`
	SeriesId         uint   `json:"series_id" validate:"required"`
}

type SortFilterUnit struct {
	Field            string
	Sort             string
	UnitName         string
	BrandId          string
	HeavyEquipmentId string
	SeriesId         string
}
