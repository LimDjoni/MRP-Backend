package seeding

import (
	"ajebackend/model/company"
	"ajebackend/model/trader"
	"fmt"
	"gorm.io/gorm"
)

func FindCompanyId(companies []company.Company, companyName string) uint {
	for _, v := range companies {
		if v.CompanyName == companyName {
			return v.ID
		}
	}
	return 0
}

func SeedingTraderAndCompanyData(db *gorm.DB) {

	tx := db.Begin()
	var checkCompany []company.Company

	tx.Find(&checkCompany)

	if len(checkCompany) > 0 {
		return
	}

	var createCompany []company.Company

	createCompany = append(createCompany,
			company.Company{
				CompanyName: "PT. ANGSANA JAYA ENERGI",
				Province: "KALIMANTAN SELATAN",
				Address: "Jl. Sebamban II Dusun III Blok F N0.021 RT. 012 RW.000 Karang Indah Angsana, Kab Tanah Bumbu",
			},
			company.Company{
				CompanyName: "PT DELI NIAGA SEJAHTERA",
				Province: "JAKARTA SELATAN",
				Address: "Grand ITC Permata Hijau LT 8 Suite B No 3A",
			},
			company.Company{
				CompanyName: "PT. SEMEN INDONESIA, TBK",
				Province: "JAWA TIMUR",
				Address: "Kantor Pusat Semen Gresik - Pabrik Tuban Desa Sumberarum, Kec. Kerek - Kab. Tuban 62356",
				PhoneNumber: "031 - 3981731 - 33",
				FaxNumber: "031 - 3972260 ",
			},
			company.Company{
				CompanyName: "PT. INTI MUSTIKA KARYATAMA",
				Province: "SURABAYA",
				Address: "Jl.Jemur Sari 128 - 130, Jemur Wonosari Jemur Wonosari Wonoporo",
			},
			company.Company{
				CompanyName: "PT. SEMPURNA INDRA PRATAMA",
				Province: "SURABAYA",
				Address: "Jl. Mawar No. 26 RT. 004 RW. 003 Tegalsari , Kota Surabaya Jawa Timur 60262",
			},
			company.Company{
				CompanyName: "PT. SOLUSI BANGUN INDONESIA",
				Province: "BOGOR",
				Address: "Narogong Plant - Jl. Raya Narogong Km. 7 Bogor 16820",
				PhoneNumber: "021 - 8231260",
				FaxNumber: "021 - 8231254",
			},
		)

	err := tx.Create(&createCompany).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Company")
		return
	}

	var createTrader []trader.Trader

	emailAJE := "angsanajayaenergi123@gmail.com"
	emailDeli := "salim@deli.id "
	createTrader = append(createTrader,
			trader.Trader{
				TraderName: "RICHARD NM PALAR",
				Email: &emailAJE,
				Position: "DIREKTUR",
				CompanyId: FindCompanyId(createCompany, "PT. ANGSANA JAYA ENERGI"),
			},
			trader.Trader{
				TraderName: "SALIM LIMANTO",
				Email: &emailDeli,
				Position: "DIREKTUR",
				CompanyId: FindCompanyId(createCompany, "PT DELI NIAGA SEJAHTERA"),
			},
			trader.Trader{
				TraderName: "FACHRUR ROJI, ST",
				Position: "SM OF BU PROCUREMENT",
				CompanyId: FindCompanyId(createCompany, "PT. SEMEN INDONESIA, TBK"),
			},
			trader.Trader{
				TraderName: "IIN ISWARINI",
				Position: "DIREKTUR",
				CompanyId: FindCompanyId(createCompany, "PT. INTI MUSTIKA KARYATAMA"),
			},
			trader.Trader{
				TraderName: "MARISKA NATASYA S",
				Position: "DIREKTUR",
				CompanyId: FindCompanyId(createCompany, "PT. SEMPURNA INDRA PRATAMA"),
			},
			trader.Trader{
				TraderName: "BUDI ARIA PRATAMA",
				Position: "PROCUREMENT MANAGER CAT. ENERGY",
				CompanyId: FindCompanyId(createCompany, "PT. SOLUSI BANGUN INDONESIA"),
			},
		)

	errTrader := tx.Create(&createTrader).Error

	if errTrader != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Company")
		return
	}

	tx.Commit()
}
