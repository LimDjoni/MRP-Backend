package seeding

import (
	"ajebackend/model/master/insurancecompany"
	"fmt"

	"gorm.io/gorm"
)

func SeedingInsuranceCompany(db *gorm.DB) {

	tx := db.Begin()
	var checkInsuranceCompany []insurancecompany.InsuranceCompany
	tx.Find(&checkInsuranceCompany)

	if len(checkInsuranceCompany) > 0 {
		return
	}

	var createInsuranceCompany []insurancecompany.InsuranceCompany
	createInsuranceCompany = append(createInsuranceCompany,
		insurancecompany.InsuranceCompany{Name: "PT Asuransi Umum Mega"},
	)

	err := tx.Create(&createInsuranceCompany).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Insurance Company")
		return
	}

	tx.Commit()
}
