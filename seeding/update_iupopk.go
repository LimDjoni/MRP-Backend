package seeding

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"fmt"

	"gorm.io/gorm"
)

func UpdateIupopk(db *gorm.DB) {
	var iup iupopk.Iupopk

	iupFindErr := db.Where("name = 'PT Angsana Jaya Energi'").First(&iup).Error

	if iupFindErr != nil {
		fmt.Println("Iupopk not found while updating")
		return
	}

	db.Unscoped().Model(&minerba.Minerba{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&dmo.Dmo{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&groupingvesseldn.GroupingVesselDn{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&groupingvesselln.GroupingVesselLn{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&insw.Insw{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&minerbaln.MinerbaLn{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&production.Production{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	db.Unscoped().Model(&reportdmo.ReportDmo{}).Where("iupopk_id IS NULL").Update("iupopk_id", iup.ID)

	return
}
