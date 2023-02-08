package seeding

import (
	"ajebackend/model/master/surveyor"
	"fmt"

	"gorm.io/gorm"
)

func SeedingSurveyor(db *gorm.DB) {

	tx := db.Begin()
	var checkSurveyor []surveyor.Surveyor
	tx.Find(&checkSurveyor)

	if len(checkSurveyor) > 0 {
		return
	}

	var createSurveyor []surveyor.Surveyor
	createSurveyor = append(createSurveyor,
		surveyor.Surveyor{
			Name: "PT Alfred H Knight Testing Indonesia"},
		surveyor.Surveyor{
			Name: "PT Anindya"},
		surveyor.Surveyor{
			Name: "PT Asiatrust Technovima Qualiti"},
		surveyor.Surveyor{
			Name: "PT Carsurin"},
		surveyor.Surveyor{
			Name: "PT CICC"},
		surveyor.Surveyor{
			Name: "PT Geoservices"},
		surveyor.Surveyor{
			Name: "PT Superintending Company of Indonesia"},
		surveyor.Surveyor{
			Name: "PT Surveyor Carbon Consulting Indonesia"},
		surveyor.Surveyor{
			Name: "PT Tribhakti Inspektama"},
	)

	err := tx.Create(&createSurveyor).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Surveyor")
		return
	}

	tx.Commit()
}
