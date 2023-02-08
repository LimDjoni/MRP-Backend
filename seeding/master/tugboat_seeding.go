package seeding

import (
	"ajebackend/model/master/tugboat"
	"fmt"

	"gorm.io/gorm"
)

func SeedingTugboat(db *gorm.DB) {

	tx := db.Begin()
	var checkTugboat []tugboat.Tugboat

	tx.Find(&checkTugboat)

	if len(checkTugboat) > 0 {
		return
	}

	var createTugboat []tugboat.Tugboat

	createTugboat = append(createTugboat,
		tugboat.Tugboat{
			Name: "TB. ARK Dan"},
		tugboat.Tugboat{
			Name: "TB. ARK Levi"},
		tugboat.Tugboat{
			Name: "TB. ARK Naphtali"},
		tugboat.Tugboat{
			Name: "TB. ARK Simeon"},
		tugboat.Tugboat{
			Name: "TB. Ashleigh 01"},
		tugboat.Tugboat{
			Name: "TB. Ashleigh 02"},
		tugboat.Tugboat{
			Name: "TB. Ashleigh 07"},
		tugboat.Tugboat{
			Name: "TB. Atlantic Perkasa"},
		tugboat.Tugboat{
			Name: "TB. Baruna 2"},
		tugboat.Tugboat{
			Name: "TB. Enrico"},
		tugboat.Tugboat{
			Name: "TB. Entebe Megastar 65"},
		tugboat.Tugboat{
			Name: "TB. Entebe Power 10"},
		tugboat.Tugboat{
			Name: "TB. Entebe Power 5"},
		tugboat.Tugboat{
			Name: "TB. Entebe Power 9"},
		tugboat.Tugboat{
			Name: "TB. Entebe Star 29"},
		tugboat.Tugboat{
			Name: "TB. Gold Star"},
		tugboat.Tugboat{
			Name: "TB. Harlina 59"},
		tugboat.Tugboat{
			Name: "TB. IBC Makassar"},
		tugboat.Tugboat{
			Name: "TB. Kelly 02"},
		tugboat.Tugboat{
			Name: "TB. Kelly 03"},
		tugboat.Tugboat{
			Name: "TB. Kelly 05"},
		tugboat.Tugboat{
			Name: "TB. Kelly 06"},
		tugboat.Tugboat{
			Name: "TB. Kelly 07"},
		tugboat.Tugboat{
			Name: "TB. Kietrans 2"},
		tugboat.Tugboat{
			Name: "TB. Lautan Berlian 1"},
		tugboat.Tugboat{
			Name: "TB. Mega Power 16"},
		tugboat.Tugboat{
			Name: "TB. Mimi L 01"},
		tugboat.Tugboat{
			Name: "TB. Mimi L 02"},
		tugboat.Tugboat{
			Name: "TB. Momentum 05"},
		tugboat.Tugboat{
			Name: "TB. Momentum 08"},
		tugboat.Tugboat{
			Name: "TB. Nasya 02"},
		tugboat.Tugboat{
			Name: "TB. Pacific Eight"},
		tugboat.Tugboat{
			Name: "TB. Pacific Eighteen"},
		tugboat.Tugboat{
			Name: "TB. Pacific Eleven"},
		tugboat.Tugboat{
			Name: "TB. Pacific Fifteen"},
		tugboat.Tugboat{
			Name: "TB. Pacific Five"},
		tugboat.Tugboat{
			Name: "TB. Pacific Nine"},
		tugboat.Tugboat{
			Name: "TB. Pacific Nineteen"},
		tugboat.Tugboat{
			Name: "TB. Pacific Seven"},
		tugboat.Tugboat{
			Name: "TB. Pacific Seventeen"},
		tugboat.Tugboat{
			Name: "TB. Pacific Six"},
		tugboat.Tugboat{
			Name: "TB. Pacific Sixteen"},
		tugboat.Tugboat{
			Name: "TB. Pacific Ten"},
		tugboat.Tugboat{
			Name: "TB. Pacific Twelve"},
		tugboat.Tugboat{
			Name: "TB. Pacific Twenty"},
		tugboat.Tugboat{
			Name: "TB. Patria 12"},
		tugboat.Tugboat{
			Name: "TB. Perkasa 2"},
		tugboat.Tugboat{
			Name: "TB. Permata Dolphin"},
		tugboat.Tugboat{
			Name: "TB. Prime 10"},
		tugboat.Tugboat{
			Name: "TB. Prime 20"},
		tugboat.Tugboat{
			Name: "TB. Prime 4"},
		tugboat.Tugboat{
			Name: "TB. Prime 8"},
		tugboat.Tugboat{
			Name: "TB. PSB 05"},
		tugboat.Tugboat{
			Name: "TB. PSB 06"},
		tugboat.Tugboat{
			Name: "TB. PSB 3301"},
		tugboat.Tugboat{
			Name: "TB. PSB 3302"},
		tugboat.Tugboat{
			Name: "TB. PSB 3303"},
		tugboat.Tugboat{
			Name: "TB. Satria Laksana 108"},
		tugboat.Tugboat{
			Name: "TB. Satria Laksana 98"},
		tugboat.Tugboat{
			Name: "TB. Selwyn 3"},
		tugboat.Tugboat{
			Name: "TB. Sereia 55"},
		tugboat.Tugboat{
			Name: "TB. Sunstar"},
		tugboat.Tugboat{
			Name: "TB. Wijaya Trans 28"},
	)

	err := tx.Create(&createTugboat).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Tugboat")
		return
	}

	tx.Commit()
}
