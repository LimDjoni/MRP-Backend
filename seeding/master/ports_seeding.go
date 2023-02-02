package seeding

import (
	"ajebackend/model/master/ports"
	"fmt"

	"gorm.io/gorm"
)

func SeedingPorts(db *gorm.DB) {

	tx := db.Begin()
	var checkPorts []ports.Port
	tx.Find(&checkPorts)

	if len(checkPorts) > 0 {
		return
	}

	var createPorts []ports.Port
	createPorts = append(createPorts,
		ports.Port{Name: "Jetty Bina Indo Raya", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false},
		ports.Port{Name: "Jetty PT Deli Niaga Jaya", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false},
		ports.Port{Name: "Jetty PT Sebamban Terminal Umum", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false},
		ports.Port{Name: "Jetty Walie - Marunda", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false},
		ports.Port{Name: "KCN Marunda", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Morosi, Konawe, Kendari", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Muara Bunati", IsLoadingPort: true, IsUnloadingPort: true, IsDmoDestinationPort: false},
		ports.Port{Name: "Muara Sampara", IsLoadingPort: false, IsUnloadingPort: false, IsDmoDestinationPort: true},
		ports.Port{Name: "Muara Satui", IsLoadingPort: true, IsUnloadingPort: true, IsDmoDestinationPort: false},
		ports.Port{Name: "Pelabuhan Tanjung Emas", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Pelindo Cilegon", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false},
		ports.Port{Name: "Pelindo Cirebon", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false},
		ports.Port{Name: "PLTU Jeranjang", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "PLTU Paiton Baru", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "PLTU Paiton PJB", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "PLTU Rembang", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "PLTU Tanjung Awar-Awar", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Tanjung Merpati Tanauge", IsLoadingPort: false, IsUnloadingPort: false, IsDmoDestinationPort: true},
		ports.Port{Name: "Tanjung Priok", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Tersus PT Semen Indonesia", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
		ports.Port{Name: "Tersus PT Solusi Bangun Indonesia", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true},
	)

	err := tx.Create(&createPorts).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Ports")
		return
	}

	tx.Commit()
}
