package brand

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	BrandName string `json:"brand_name"`
}
