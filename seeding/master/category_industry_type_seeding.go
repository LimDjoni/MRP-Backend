package seeding

import (
	"ajebackend/model/master/categoryindustrytype"
	"fmt"

	"gorm.io/gorm"
)

func SeedingCategoryIndustryType(db *gorm.DB) {
	tx := db.Begin()
	var checkCategory []categoryindustrytype.CategoryIndustryType

	tx.Find(&checkCategory)

	if len(checkCategory) > 0 {
		return
	}

	var createCategory []categoryindustrytype.CategoryIndustryType

	createCategory = append(createCategory,
		categoryindustrytype.CategoryIndustryType{
			Name:       "Kelistrikan",
			SystemName: "electricity",
			Order:      1,
		},
		categoryindustrytype.CategoryIndustryType{
			Name:       "Semen",
			SystemName: "cement",
			Order:      2,
		},
		categoryindustrytype.CategoryIndustryType{
			Name:       "Smelter",
			SystemName: "non_electricity",
			Order:      3,
		},
	)

	err := tx.Create(&createCategory).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Category")
		return
	}

	tx.Commit()
}
