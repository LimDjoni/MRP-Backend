package industrytype

import (
	"ajebackend/model/master/categoryindustrytype"

	"gorm.io/gorm"
)

type IndustryType struct {
	gorm.Model
	Name                   string                                    `json:"name" gorm:"UNIQUE"`
	Category               string                                    `json:"category"`
	SystemCategory         string                                    `json:"system_category"`
	CategoryIndustryTypeId uint                                      `json:"category_industry_type_id"`
	CategoryIndustryType   categoryindustrytype.CategoryIndustryType `json:"category_industry_type"`
}
