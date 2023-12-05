package seeding

import (
	"ajebackend/model/master/role"
	"fmt"

	"gorm.io/gorm"
)

func SeedingRole(db *gorm.DB) {

	tx := db.Begin()
	var checkRole []role.Role

	tx.Find(&checkRole)

	if len(checkRole) > 0 {
		return
	}

	var createRole []role.Role

	createRole = append(createRole,
		role.Role{
			Name: "superuser",
		},
		role.Role{
			Name: "shipping",
		},
		role.Role{
			Name: "accounting",
		},
		role.Role{
			Name: "admin",
		},
		role.Role{
			Name: "supervisor",
		},
	)

	err := tx.Create(&createRole).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding user role")
		return
	}

	tx.Commit()
}
