package series

import (
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/heavyequipment"

	"gorm.io/gorm"
)

type Series struct {
	gorm.Model
	SeriesName       string `json:"series_name"`
	BrandId          uint   `json:"brand_id"`
	HeavyEquipmentId uint   `json:"heavy_equipment_id"`

	Brand          brand.Brand                   `gorm:"foreignKey:BrandId" json:"brand"`                    // Add this line
	HeavyEquipment heavyequipment.HeavyEquipment `gorm:"foreignKey:HeavyEquipmentId" json:"heavy_equipment"` // Add this line
}
