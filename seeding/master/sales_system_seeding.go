package seeding

import (
	"ajebackend/model/master/salessystem"
	"fmt"

	"gorm.io/gorm"
)

func SeedingSalesSystem(db *gorm.DB) {

	tx := db.Begin()
	var checkSalesSystem []salessystem.SalesSystem

	tx.Find(&checkSalesSystem)

	if len(checkSalesSystem) > 0 {
		return
	}

	var createSalesSystem []salessystem.SalesSystem

	createSalesSystem = append(createSalesSystem,
		salessystem.SalesSystem{
			Name: "FOB Barge",
		},
		salessystem.SalesSystem{
			Name: "FOB Vessel",
		},
		salessystem.SalesSystem{
			Name: "CIF Barge",
		},
		salessystem.SalesSystem{
			Name: "CIF Vessel",
		},
		salessystem.SalesSystem{
			Name: "CNF Barge",
		},
		salessystem.SalesSystem{
			Name: "CNF Vessel",
		},
	)

	err := tx.Create(&createSalesSystem).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Sales System")
		return
	}

	tx.Commit()
}
