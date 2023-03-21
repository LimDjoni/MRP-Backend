package seeding

import (
	"ajebackend/model/counter"
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

	emailAJE := "angsanajayaenergi123@gmail.com"
	emailTantra := "traffic.operationtmd@gmail.com"

	var createIupopk []iupopk.Iupopk
	createIupopk = append(createIupopk,
		iupopk.Iupopk{
			Name:         "PT Angsana Jaya Energi",
			Address:      "Jl. Sebamban II Dusun III Blok F N0.021 RT. 012 RW.000 Karang Indah Angsana, Kab Tanah Bumbu",
			Province:     "Kalimantan Selatan",
			Email:        &emailAJE,
			PhoneNumber:  nil,
			FaxNumber:    nil,
			DirectorName: "Richard NM Palar",
			Position:     "Direktur",
			Code:         "AJE",
			Location:     "Tanah Bumbu, Provinsi Kalimantan Selatan",
		},
		iupopk.Iupopk{
			Name:         "PT Tantra Mining Development",
			Address:      "Jalan R. Soeprapto, No. 25, Banjarmasin, Kalimantan Selatan",
			Province:     "Kalimantan Selatan",
			Email:        &emailTantra,
			PhoneNumber:  nil,
			FaxNumber:    nil,
			DirectorName: "Yansen Andriyan",
			Position:     "Direktur",
			Code:         "TMD",
			Location:     "Tanah Bumbu, Provinsi Kalimantan Selatan",
		},
	)

	err := tx.Create(&createIupopk).Error
	if err != nil {
		fmt.Println(err.Error())
		tx.Rollback()
		fmt.Println("Failed Seeding Iupopk")
		return
	}
	var counters []counter.Counter

	for _, v := range createIupopk {
		counters = append(counters, counter.Counter{
			IupopkId:      v.ID,
			TransactionDn: 1,
			TransactionLn: 1,
			GroupingMvDn:  1,
			GroupingMvLn:  1,
			Sp3medn:       1,
			Sp3meln:       1,
			BaEndUser:     1,
			Dmo:           1,
			Production:    1,
			Insw:          1,
		})
	}

	createCounterErr := tx.Create(&counters).Error

	if createCounterErr != nil {
		fmt.Println(createCounterErr.Error())
		tx.Rollback()
		fmt.Println("Failed Seeding Iupopk")
		return
	}
	tx.Commit()
}
