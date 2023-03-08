package seeding

import (
	"ajebackend/model/master/iupopk"
	"fmt"

	"gorm.io/gorm"
)

func SeedingIupopk(db *gorm.DB) {

	tx := db.Begin()
	var checkIupopk []iupopk.Iupopk
	tx.Find(&checkIupopk)

	if len(checkIupopk) > 0 {
		return
	}

	var createIupopk []iupopk.Iupopk
	createIupopk = append(createIupopk,
		iupopk.Iupopk{
			Name:         "PT Angsana Jaya Energi",
			Address:      "Jl. Sebamban II Dusun III Blok F N0.021 RT. 012 RW.000 Karang Indah Angsana, Kab Tanah Bumbu",
			Province:     "Kalimantan Selatan",
			Email:        "angsanajayaenergi123@gmail.com",
			PhoneNumber:  "",
			FaxNumber:    "",
			DirectorName: "Richard NM Palar",
			Position:     "Direktur",
			Code:         "AJE",
		},
		iupopk.Iupopk{
			Name:         "PT Tantra Mining Development",
			Address:      "Jalan R. Soeprapto, No. 25, Banjarmasin, Kalimantan Selatan",
			Province:     "Kalimantan Selatan",
			Email:        "traffic.operationtmd@gmail.com",
			PhoneNumber:  "",
			FaxNumber:    "",
			DirectorName: "Yansen Andriyan",
			Position:     "Direktur",
			Code:         "TMD",
		},
	)

	err := tx.Create(&createIupopk).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Iupopk")
		return
	}

	tx.Commit()
}
