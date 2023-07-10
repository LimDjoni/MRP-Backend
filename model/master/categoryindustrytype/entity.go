package categoryindustrytype

import "gorm.io/gorm"

type CategoryIndustryType struct {
	gorm.Model
	Name string `json:"name"`
}
