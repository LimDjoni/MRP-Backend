package seeding

import (
	"ajebackend/model/master/portlocation"
	"ajebackend/model/master/ports"
	"fmt"

	"gorm.io/gorm"
)

func FindLocationId(locations []portlocation.PortLocation, locationName string) uint {
	for _, v := range locations {
		if v.Name == locationName {
			return v.ID
		}
	}
	return 0
}

func SeedingPortsAndLocation(db *gorm.DB) {

	tx := db.Begin()

	var checkPortLocation []portlocation.PortLocation
	tx.Find(&checkPortLocation)

	if len(checkPortLocation) > 0 {
		return
	}

	var createPortLocation []portlocation.PortLocation
	createPortLocation = append(createPortLocation,
		portlocation.PortLocation{
			Name: "Banten"},
		portlocation.PortLocation{
			Name: "DKI Jakarta"},
		portlocation.PortLocation{
			Name: "Jawa Barat"},
		portlocation.PortLocation{
			Name: "Jawa Tengah"},
		portlocation.PortLocation{
			Name: "Jawa Timur"},
		portlocation.PortLocation{
			Name: "Kalimantan Selatan"},
		portlocation.PortLocation{
			Name: "Sulawesi Tenggara"},
		portlocation.PortLocation{
			Name: "Nusa Tenggara Barat"},
		portlocation.PortLocation{
			Name: "Sulawesi Tengah"},
	)

	errPortLocation := tx.Create(&createPortLocation).Error

	if errPortLocation != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Port Location")
		return
	}

	var checkPorts []ports.Port
	tx.Find(&checkPorts)

	if len(checkPorts) > 0 {
		return
	}

	var createPorts []ports.Port
	createPorts = append(createPorts,
		ports.Port{Name: "Jetty Bina Indo Raya", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Kalimantan Selatan")},
		ports.Port{Name: "Jetty PT Deli Niaga Jaya", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Kalimantan Selatan")},
		ports.Port{Name: "Jetty PT Sebamban Terminal Umum", IsLoadingPort: true, IsUnloadingPort: false, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Kalimantan Selatan")},
		ports.Port{Name: "Jetty Walie - Marunda", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "DKI Jakarta")},
		ports.Port{Name: "KCN Marunda", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "DKI Jakarta")},
		ports.Port{Name: "Morosi, Konawe, Kendari", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Sulawesi Tenggara")},
		ports.Port{Name: "Muara Bunati", IsLoadingPort: true, IsUnloadingPort: true, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Kalimantan Selatan")},
		ports.Port{Name: "Muara Sampara", IsLoadingPort: false, IsUnloadingPort: false, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Sulawesi Tenggara")},
		ports.Port{Name: "Muara Satui", IsLoadingPort: true, IsUnloadingPort: true, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Kalimantan Selatan")},
		ports.Port{Name: "Pelabuhan Tanjung Emas", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Tengah")},
		ports.Port{Name: "Pelindo Cilegon", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Banten")},
		ports.Port{Name: "Pelindo Cirebon", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: false, PortLocationId: FindLocationId(createPortLocation, "Jawa Barat")},
		ports.Port{Name: "PLTU Jeranjang", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Nusa Tenggara Barat")},
		ports.Port{Name: "PLTU Paiton Baru", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Timur")},
		ports.Port{Name: "PLTU Paiton PJB", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Timur")},
		ports.Port{Name: "PLTU Rembang", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Tengah")},
		ports.Port{Name: "PLTU Tanjung Awar-Awar", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Timur")},
		ports.Port{Name: "Tanjung Merpati Tanauge", IsLoadingPort: false, IsUnloadingPort: false, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Sulawesi Tengah")},
		ports.Port{Name: "Tanjung Priok", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "DKI Jakarta")},
		ports.Port{Name: "Tersus PT Semen Indonesia", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Timur")},
		ports.Port{Name: "Tersus PT Solusi Bangun Indonesia", IsLoadingPort: false, IsUnloadingPort: true, IsDmoDestinationPort: true, PortLocationId: FindLocationId(createPortLocation, "Jawa Timur")},
	)

	err := tx.Create(&createPorts).Error

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		fmt.Println("Failed Seeding Ports")
		return
	}

	tx.Commit()
}
