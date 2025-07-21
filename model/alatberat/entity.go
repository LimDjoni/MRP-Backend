package alatberat

import (
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/heavyequipment"
	"mrpbackend/model/master/series"

	"gorm.io/gorm"
)

type AlatBerat struct {
	gorm.Model
	BrandId          uint    `json:"brand_id"`
	HeavyEquipmentId uint    `json:"heavy_equipment_id"`
	SeriesId         uint    `json:"series_id"`
	Consumption      float64 `json:"consumption"`
	Tolerance        uint    `json:"tolerance"`

	Brand          brand.Brand                   `gorm:"foreignKey:BrandId" json:"brand"`                    // Add this line
	HeavyEquipment heavyequipment.HeavyEquipment `gorm:"foreignKey:HeavyEquipmentId" json:"heavy_equipment"` // Add this line
	Series         series.Series                 `gorm:"foreignKey:SeriesId" json:"series"`                  // Add this line
}
