package seeding

import (
	"ajebackend/model/master/vessel"
	"fmt"

	"gorm.io/gorm"
)

func SeedingVessel(db *gorm.DB) {

	tx := db.Begin()
	var checkVessel []vessel.Vessel

	tx.Find(&checkVessel)

	if len(checkVessel) > 0 {
		return
	}

	var createVessel []vessel.Vessel

	createVessel = append(createVessel,
		vessel.Vessel{
			Name: "MV. Abdul Hamid"},
		vessel.Vessel{
			Name: "MV. Daidan Mustikawati"},
		vessel.Vessel{
			Name: "MV. Daidan Pertiwi"},
		vessel.Vessel{
			Name: "MV. Densa Eagle"},
		vessel.Vessel{
			Name: "MV. DRY Transport"},
		vessel.Vessel{
			Name: "MV. LGH Prosper"},
		vessel.Vessel{
			Name: "MV. Lumoso Aman"},
		vessel.Vessel{
			Name: "MV. Lumoso Harmoni"},
		vessel.Vessel{
			Name: "MV. Lumoso Jaya"},
		vessel.Vessel{
			Name: "MV. MBS Baluran"},
		vessel.Vessel{
			Name: "MV. MDM Bromo"},
	)

	err := tx.Create(&createVessel).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Vessel")
		return
	}

	tx.Commit()
}
