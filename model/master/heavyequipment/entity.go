package heavyequipment

import (
	"mrpbackend/model/master/brand"

	"gorm.io/gorm"
)

type HeavyEquipment struct {
	gorm.Model
	BrandId            uint   `json:"brand_id"`
	HeavyEquipmentName string `json:"heavy_equipment_name"`

	Brand brand.Brand `gorm:"foreignKey:BrandId" json:"brand"` // Add this line
}
