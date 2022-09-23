package seeding

import (
	"ajebackend/model/trader"
	"fmt"
	"gorm.io/gorm"
)

func SeedingTraderData(db *gorm.DB) {

	var checkTrader []trader.Trader

	db.Find(&checkTrader)

	if len(checkTrader) > 0 {
		return
	}
	var traderCompanyName []trader.Trader

	traderCompanyName = append(traderCompanyName,
		trader.Trader{CompanyName: "PT. DELI NIAGA SEJAHTERA"},
		trader.Trader{CompanyName: "PT. DELI NIAGA JAYA"},
		trader.Trader{CompanyName: "PT. GEO MINERAL TRADING"},
		trader.Trader{CompanyName: "PT. INTI MUSTIKA KARYATAMA"},
		trader.Trader{CompanyName: "PT. BERKAT RAYA OPTIMA"},
		trader.Trader{CompanyName: "PT. ANAGA ABYUDAYA ANANTA"},
		trader.Trader{CompanyName: "PT. DAIDAN ADITAMA YAKSA"},
		trader.Trader{CompanyName: "PT. VIRTUE DRAGON NICKEL INDUSTRY"},
		trader.Trader{CompanyName: "PT. SEMPURNA INDRA PRATAMA"},
		trader.Trader{CompanyName: "PT. MINERATAMA  PRIMA ABADI"},
		trader.Trader{CompanyName: "PT. VIREMA IMPEX"},
		trader.Trader{CompanyName: "PT. BARA INDAH SINERGI"},
		trader.Trader{CompanyName: "PT. MITRA BARA ABADI BANDUNG",})

	err := db.Create(&traderCompanyName).Error

	if err != nil {
		fmt.Println("Failed Seeding Trader")
		return
	}
}
