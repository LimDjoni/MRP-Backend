package seeding

import (
	"ajebackend/model/master/userrole"
	"fmt"

	"gorm.io/gorm"
)

func SeedingUserRole(db *gorm.DB) {

	tx := db.Begin()
	var checkUserRole []userrole.UserRole

	tx.Find(&checkUserRole)

	if len(checkUserRole) > 0 {
		return
	}

	var createUserRole []userrole.UserRole

	createUserRole = append(createUserRole,
		userrole.UserRole{
			Name: "superuser",
		},
		userrole.UserRole{
			Name: "shipping",
		},
		userrole.UserRole{
			Name: "accounting",
		},
		userrole.UserRole{
			Name: "admin",
		},
		userrole.UserRole{
			Name: "supervisor",
		},
	)

	err := tx.Create(&createUserRole).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding user role")
		return
	}

	tx.Commit()
}
