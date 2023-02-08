package seeding

import (
	"ajebackend/model/master/unit"
	"fmt"

	"gorm.io/gorm"
)

func SeedingUnit(db *gorm.DB) {

	tx := db.Begin()
	var checkUnit []unit.Unit
	tx.Find(&checkUnit)

	if len(checkUnit) > 0 {
		return
	}

	var createUnit []unit.Unit
	createUnit = append(createUnit,
		unit.Unit{Name: "Metric ton (1000 kg)", Code: "TNE"},
	)

	err := tx.Create(&createUnit).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Unit")
		return
	}

	tx.Commit()
}
