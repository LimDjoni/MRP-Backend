package seeding

import (
	"ajebackend/model/master/portinsw"
	"fmt"

	"gorm.io/gorm"
)

func SeedingPortInsw(db *gorm.DB) {
	tx := db.Begin()
	var checkPortInsw []portinsw.PortInsw
	tx.Find(&checkPortInsw)

	if len(checkPortInsw) > 0 {
		return
	}

	var createPortInsw []portinsw.PortInsw
	createPortInsw = append(createPortInsw,
		portinsw.PortInsw{Name: "Satui", Code: "IDSTU"},
	)

	err := tx.Create(&createPortInsw).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Port Insw")
		return
	}

	tx.Commit()
}
