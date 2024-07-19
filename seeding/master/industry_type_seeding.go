package seeding

import (
	"ajebackend/model/master/industrytype"
	"fmt"

	"gorm.io/gorm"
)

func SeedingIndustryType(db *gorm.DB) {

	tx := db.Begin()
	var checkIndustryType []industrytype.IndustryType

	tx.Find(&checkIndustryType)

	if len(checkIndustryType) > 0 {
		return
	}

	var createIndustryType []industrytype.IndustryType

	createIndustryType = append(createIndustryType,
		industrytype.IndustryType{
			Name:                   "Pembangkit Listrik",
			Category:               "ELECTRICITY",
			CategoryIndustryTypeId: 1,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Besi",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Nikel",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Gula",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Kertas",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Semen",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 2,
		},
		industrytype.IndustryType{
			Name:                   "Pabrik Tekstil",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Perdagangan Besar Bahan Bakar Padat, Cair, dan Gas dan Produk Ybdi",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
		industrytype.IndustryType{
			Name:                   "Trader Batubara",
			Category:               "NON ELECTRICITY",
			CategoryIndustryTypeId: 3,
		},
	)

	err := tx.Create(&createIndustryType).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Industry Type")
		return
	}

	tx.Commit()
}
