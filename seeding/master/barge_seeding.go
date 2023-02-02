package seeding

import (
	"ajebackend/model/master/barge"
	"fmt"

	"gorm.io/gorm"
)

func SeedingBarge(db *gorm.DB) {

	tx := db.Begin()
	var checkBarge []barge.Barge

	tx.Find(&checkBarge)

	if len(checkBarge) > 0 {
		return
	}

	var createBarge []barge.Barge

	createBarge = append(createBarge,
		barge.Barge{
			Name: "BG. Alphecca"},
		barge.Barge{
			Name: "BG. ARK Amethyst"},
		barge.Barge{
			Name: "BG. ARK Citrine"},
		barge.Barge{
			Name: "BG. ARK Emerald"},
		barge.Barge{
			Name: "BG. ARK Peridot"},
		barge.Barge{
			Name: "BG. Asherah 18"},
		barge.Barge{
			Name: "BG. Bahtera 3005"},
		barge.Barge{
			Name: "BG. BKT 301"},
		barge.Barge{
			Name: "BG. BPL 2"},
		barge.Barge{
			Name: "BG. Chaterina 01"},
		barge.Barge{
			Name: "BG. Chaterina 02"},
		barge.Barge{
			Name: "BG. Edeline"},
		barge.Barge{
			Name: "BG. Finacia 57"},
		barge.Barge{
			Name: "BG. Finacia 58"},
		barge.Barge{
			Name: "BG. Finacia 63"},
		barge.Barge{
			Name: "BG. Finacia 78"},
		barge.Barge{
			Name: "BG. Finacia 86"},
		barge.Barge{
			Name: "BG. Finacia 87"},
		barge.Barge{
			Name: "BG. Finacia 96"},
		barge.Barge{
			Name: "BG. GHI 01"},
		barge.Barge{
			Name: "BG. GHI 02"},
		barge.Barge{
			Name: "BG. GHI 05"},
		barge.Barge{
			Name: "BG. GHI 06"},
		barge.Barge{
			Name: "BG. GHI 07"},
		barge.Barge{
			Name: "BG. Indo Marina"},
		barge.Barge{
			Name: "BG. Jingxi 01"},
		barge.Barge{
			Name: "BG. Jingxi 02"},
		barge.Barge{
			Name: "BG. Jingxi 05"},
		barge.Barge{
			Name: "BG. Jingxi 07"},
		barge.Barge{
			Name: "BG. Lius Star"},
		barge.Barge{
			Name: "BG. Momentun 3008"},
		barge.Barge{
			Name: "BG. Momentun 3009"},
		barge.Barge{
			Name: "BG. Pacific 3006"},
		barge.Barge{
			Name: "BG. Pacific 3007"},
		barge.Barge{
			Name: "BG. Pacific 3008"},
		barge.Barge{
			Name: "BG. Pacific 3009"},
		barge.Barge{
			Name: "BG. Pacific 3010"},
		barge.Barge{
			Name: "BG. Pacific 3011"},
		barge.Barge{
			Name: "BG. Pacific 3012"},
		barge.Barge{
			Name: "BG. Pacific 3015"},
		barge.Barge{
			Name: "BG. Pacific 3016"},
		barge.Barge{
			Name: "BG. Pacific 3017"},
		barge.Barge{
			Name: "BG. Pacific 3018"},
		barge.Barge{
			Name: "BG. Pacific 3019"},
		barge.Barge{
			Name: "BG. Pacific 3020"},
		barge.Barge{
			Name: "BG. Pacific 3302"},
		barge.Barge{
			Name: "BG. Parta Jaya 3005"},
		barge.Barge{
			Name: "BG. PSB 01"},
		barge.Barge{
			Name: "BG. PSB 02"},
		barge.Barge{
			Name: "BG. PSB 03"},
		barge.Barge{
			Name: "BG. PSB 3005"},
		barge.Barge{
			Name: "BG. PSB 3006"},
		barge.Barge{
			Name: "BG. PSPM 2"},
		barge.Barge{
			Name: "BG. Rezeki Lautan 1"},
		barge.Barge{
			Name: "BG. Satria Laut 3058"},
		barge.Barge{
			Name: "BG. Satria Laut 3098"},
		barge.Barge{
			Name: "BG. Soekawati 2705"},
		barge.Barge{
			Name: "BG. Soekawati 2711"},
		barge.Barge{
			Name: "BG. Star Marine 3002"},
		barge.Barge{
			Name: "BG. Sumanggala"},
		barge.Barge{
			Name: "BG. Support 15"},
		barge.Barge{
			Name: "BG. Support 16"},
		barge.Barge{
			Name: "BG. Support 4"},
		barge.Barge{
			Name: "BG. Support 5"},
		barge.Barge{
			Name: "BG. Wijaya Trans 228"},
	)

	err := tx.Create(&createBarge).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Tugboat")
		return
	}

	tx.Commit()
}
