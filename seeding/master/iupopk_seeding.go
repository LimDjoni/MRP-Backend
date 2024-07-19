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

	emailTriop := "triop@gmail.com"
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
		iupopk.Iupopk{
			Name:         "PT Tri Oetama Persada",
			Address:      "Jl. Triop",
			Province:     "Kalimantan Tengah",
			Email:        &emailTriop,
			PhoneNumber:  nil,
			FaxNumber:    nil,
			DirectorName: "Triop",
			Position:     "Direktur",
			Code:         "TRP",
			Location:     "Tanah Bumbu, Provinsi Kalimantan Tengah",
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

	formatAje := "YYYY/BAST/CODE/MM/COUNTER"
	formatTmd := "BAST/CODE/YYYY/MM/COUNTER"

	for _, v := range createIupopk {
		var formatBast string

		if v.Name == "PT Angsana Jaya Energi" {
			formatBast = formatAje
		}

		if v.Name == "PT Tantra Mining Development" {
			formatBast = formatTmd
		}

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
			BastFormat:    formatBast,
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
