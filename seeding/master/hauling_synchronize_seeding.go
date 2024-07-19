package seeding

import (
	"ajebackend/model/haulingsynchronize"
	"ajebackend/model/master/iupopk"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedingHaulingSynchronize(db *gorm.DB) {
	tx := db.Begin()
	var checkHaulingSynchronize []haulingsynchronize.HaulingSynchronize
	tx.Find(&checkHaulingSynchronize)

	if len(checkHaulingSynchronize) > 0 {
		return
	}

	var iup iupopk.Iupopk
	tx.Where("code = ?", "TRP").First(&iup)

	var createHaulingsynchronize []haulingsynchronize.HaulingSynchronize
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	createHaulingsynchronize = append(createHaulingsynchronize,
		haulingsynchronize.HaulingSynchronize{
			LastSynchronizeIsp:   &timeNow,
			LastSynchronizeJetty: &timeNow,
			IupopkId:             iup.ID,
			FailedCount:          0,
		})
	err := tx.Create(&createHaulingsynchronize).Error

	if err != nil {
		fmt.Println(err.Error())
		tx.Rollback()
		fmt.Println("Failed Seeding Iupopk")
		return
	}

	tx.Commit()
}
