package seeding

import (
	"ajebackend/model/master/destination"
	"fmt"

	"gorm.io/gorm"
)

func SeedingDestination(db *gorm.DB) {

	tx := db.Begin()
	var checkDestination []destination.Destination

	tx.Find(&checkDestination)

	if len(checkDestination) > 0 {
		return
	}

	var createDestination []destination.Destination

	createDestination = append(createDestination,
		destination.Destination{
			Name: "Domestic",
		},
		destination.Destination{
			Name: "Export",
		},
	)

	err := tx.Create(&createDestination).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Company")
		return
	}

	tx.Commit()
}
