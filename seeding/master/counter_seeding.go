package seeding

import (
	"ajebackend/model/counter"
	"ajebackend/model/master/iupopk"
	"fmt"

	"gorm.io/gorm"
)

func SeedingCounter(db *gorm.DB) {
	tx := db.Begin()
	var checkCounter []counter.Counter
	tx.Find(&checkCounter)

	if len(checkCounter) > 0 {
		return
	}

	var checkIupopk []iupopk.Iupopk

	tx.Find(&checkIupopk)

	if len(checkIupopk) == 0 {
		return
	}

	var createCounter []counter.Counter
	for _, v := range checkIupopk {
		createCounter = append(createCounter,
			counter.Counter{
				IupopkId:           v.ID,
				TransactionDn:      1,
				TransactionLn:      1,
				GroupingMvDn:       1,
				GroupingMvLn:       1,
				Sp3medn:            1,
				Sp3meln:            1,
				BaEndUser:          1,
				Dmo:                1,
				Production:         1,
				Insw:               1,
				CoaReport:          1,
				CoaReportLn:        1,
				Rkab:               1,
				ElectricAssignment: 1,
				CafAssignment:      1,
				RoyaltyRecon:       1,
				RoyaltyReport:      1,
			},
		)
	}

	err := tx.Create(&createCounter).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Counter")
		return
	}

	tx.Commit()
}
