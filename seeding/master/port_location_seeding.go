package seeding

import (
	"ajebackend/model/master/portlocation"
	"fmt"

	"gorm.io/gorm"
)

func SeedingPortLocation(db *gorm.DB) {

	tx := db.Begin()
	var checkPortLocation []portlocation.PortLocation
	tx.Find(&checkPortLocation)

	if len(checkPortLocation) > 0 {
		return
	}

	var createPortLocation []portlocation.PortLocation
	createPortLocation = append(createPortLocation,
		portlocation.PortLocation{
			Name: "Banten"},
		portlocation.PortLocation{
			Name: "DKI Jakarta"},
		portlocation.PortLocation{
			Name: "Jawa Barat"},
		portlocation.PortLocation{
			Name: "Jawa Tengah"},
		portlocation.PortLocation{
			Name: "Jawa Timur"},
		portlocation.PortLocation{
			Name: "Kalimantan Selatan"},
		portlocation.PortLocation{
			Name: "Sulawesi Tenggara"},
	)

	err := tx.Create(&createPortLocation).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Port Location")
		return
	}

	tx.Commit()
}
