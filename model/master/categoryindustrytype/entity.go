package categoryindustrytype

import "gorm.io/gorm"

type CategoryIndustryType struct {
	gorm.Model
	Name       string `json:"name"`
	SystemName string `json:"system_name"`
	Order      int    `json:"order"`
}
