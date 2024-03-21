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
			Name: "Pembangkit Listrik",
		},
		industrytype.IndustryType{
			Name: "Pabrik Besi",
		},
		industrytype.IndustryType{
			Name: "Pabrik Nikel",
		},
		industrytype.IndustryType{
			Name: "Pabrik Gula",
		},
		industrytype.IndustryType{
			Name: "Pabrik Kertas",
		},
		industrytype.IndustryType{
			Name: "Pabrik Semen",
		},
		industrytype.IndustryType{
			Name: "Pabrik Tekstil",
		},
		industrytype.IndustryType{
			Name: "Perdagangan Besar Bahan Bakar Padat, Cair, dan Gas dan Produk Ybdi",
		},
		industrytype.IndustryType{
			Name: "Trader Batubara",
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
