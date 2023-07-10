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
			Name: "Kelistrikan",
		},
		categoryindustrytype.CategoryIndustryType{
			Name: "Semen",
		},
		categoryindustrytype.CategoryIndustryType{
			Name: "Smelter",
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
