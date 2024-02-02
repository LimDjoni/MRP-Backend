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
		icilevel.IciLevel{Name: "ICI LEVEL 1"},
		icilevel.IciLevel{Name: "ICI LEVEL 2"},
		icilevel.IciLevel{Name: "ICI LEVEL 3"},
		icilevel.IciLevel{Name: "ICI LEVEL 4"},
		icilevel.IciLevel{Name: "ICI LEVEL 5"},
	)

	err := tx.Create(&createIciLevel).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding ICI Level")
		return
	}

	tx.Commit()
}
