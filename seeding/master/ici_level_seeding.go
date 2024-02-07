package seeding

import (
	"ajebackend/model/icilevel"
	"fmt"

	"gorm.io/gorm"
)

func SeedingIciLevel(db *gorm.DB) {
	tx := db.Begin()
	var checkIciLevel []icilevel.IciLevel
	tx.Find(&checkIciLevel)

	if len(checkIciLevel) > 0 {
		return
	}

	var createIciLevel []icilevel.IciLevel
	createIciLevel = append(
		createIciLevel,
		icilevel.IciLevel{Name: "ICI LEVEL 1", Description: "GAR 6500 Weekly"},
		icilevel.IciLevel{Name: "ICI LEVEL 2", Description: "GAR 5800 Weekly"},
		icilevel.IciLevel{Name: "ICI LEVEL 3", Description: "GAR 5000 Weekly"},
		icilevel.IciLevel{Name: "ICI LEVEL 4", Description: "GAR 4200 Weekly"},
		icilevel.IciLevel{Name: "ICI LEVEL 5", Description: "GAR 3400 Weekly"},
	)

	err := tx.Create(&createIciLevel).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding ICI Level")
		return
	}

	tx.Commit()
}
