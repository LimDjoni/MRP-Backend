package seeding

import (
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
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
			CompanyName: "PT Angsana Jaya Energi",
			Province:    "Kalimantan Selatan",
			Address:     "Jl. Sebamban II Dusun III Blok F N0.021 RT. 012 RW.000 Karang Indah Angsana, Kab Tanah Bumbu",
		},
		company.Company{
			CompanyName: "PT Deli Niaga Sejahtera",
			Province:    "Jakarta Selatan",
			Address:     "Grand ITC Permata Hijau LT 8 Suite B No 3A",
		},
		company.Company{
			CompanyName: "PT Semen Indonesia, Tbk",
			Province:    "Jawa Timur",
			Address:     "Kantor Pusat Semen Gresik - Pabrik Tuban Desa Sumberarum, Kec. Kerek - Kab. Tuban 62356",
			PhoneNumber: "031 - 3981731 - 33",
			FaxNumber:   "031 - 3972260 ",
		},
		company.Company{
			CompanyName: "PT Inti Mustika Karyatama",
			Province:    "Surabaya",
			Address:     "Jl.Jemur Sari 128 - 130, Jemur Wonosari Jemur Wonosari Wonoporo",
		},
		company.Company{
			CompanyName: "PT Sempurna Indra Pratama",
			Province:    "Surabaya",
			Address:     "Jl. Mawar No. 26 RT. 004 RW. 003 Tegalsari , Kota Surabaya Jawa Timur 60262",
		},
		company.Company{
			CompanyName: "PT Solusi Bangun Indonesia",
			Province:    "Bogor",
			Address:     "Narogong Plant - Jl. Raya Narogong Km. 7 Bogor 16820",
			PhoneNumber: "021 - 8231260",
			FaxNumber:   "021 - 8231254",
		},
		company.Company{
			CompanyName: "PT Tantra Mining Development",
			Province:    "Kalimantan Selatan",
			Address:     "Jalan R. Soeprapto, No. 25, Banjarmasin, Kalimantan Selatan",
		},
		company.Company{
			CompanyName: "PT Daidan Aditama Yaksa",
			Province:    "Jakarta Pusat",
			Address:     "Gedung Town House Lt.5 JL Sungai Gerong No.1&1A RT.010 RW.020 Kebon Melati - Tanah Abang",
		},
		company.Company{
			CompanyName: "PT Anaga Abyudaya Ananta",
			Province:    "Jakarta Pusat",
			Address:     "Gedung Town House Lt.5 JL Sungai Gerong No.1&1A RT.010 RW.020 Kebon Melati - Tanah Abang",
		},
		company.Company{
			CompanyName: "PT Obsidian Stainless Steel",
			Province:    "Jakarta Selatan",
			Address:     "Gedung Bursa Efek Indonesia Tower I Lantai 27 Kel. Senayan, Kec. Kebayoran Baru",
		},
		company.Company{
			CompanyName: "PT Virtue Dragon Nickel Industry",
			Province:    "Jakarta Selatan",
			Address:     "Indonesia Stock Exchange Building 1st Tower 31st Floor Suite 310 Jl Jendral Sudirman Kav. 52-53",
			PhoneNumber: "0592 - 56912015",
			FaxNumber:   "0592 - 2631120",
		},
		company.Company{
			CompanyName: "PT Deli Niaga Jaya",
			Province:    "Kalimantan Selatan",
			Address:     "Jl Desa Satui Barat No.02 RT.005 RW.003 Satui Barat, Satui",
		},
		company.Company{
			CompanyName: "PT Geo Mineral Trading",
			Province:    "Jakarta Utara",
			Address:     "The Suite Tower Lantai 17 Jl. Boulevard Pantai Indah Kapuk No.1 RT. 001 RW.006",
		},
		company.Company{
			CompanyName: "PT Gunbuster Nickel Industry",
			Province:    "Jakarta Selatan",
			Address:     "Jl. Jenderal Sudirman Kav. 52-53, Senayan, Kebayoran Baru Kota Adm. Jakarta Selatan DKI Jakarta, Indonesia",
			PhoneNumber: "021 - 515 - 1530",
		},
		company.Company{
			CompanyName: "PT Bara Indah Sinergi",
			Province:    "Jakarta Utara",
			Address:     "Jl. Pluit Permai Raya No. 128, RT.004 RW.004, Penjaringan Jakarta Utara, DKI Jakarta",
			PhoneNumber: "021 - 260 63894",
			FaxNumber:   "021 - 260 63891",
		},
		company.Company{
			CompanyName: "PT Sinergi Mitra Investama",
			Province:    "Jawa Timur",
			Address:     "Jl. Awikoen Blok A-7, KB. Dalem, Sidokumpul Kec. Gresik, Kab. Gresik, Jawa Timur 61122",
			PhoneNumber: "031 - 3970374",
		},
		company.Company{
			CompanyName: "PT Semen Gresik",
			Province:    "Jawa Timur",
			Address:     "Kantor Pusat Semen Gresik - Pabrik Tuban Desa Sumberarum, Kec. Kerek - Kab. Tuban 62356",
			PhoneNumber: "031 - 3981731 - 33",
			FaxNumber:   "031 - 3972260",
		},
		company.Company{
			CompanyName: "PT Kaldera Energi Nusantara",
			Province:    "Banten",
			Address:     "Wisma Kaldera, Jalan Cabe V No. 52 A Pondok Cabe Ilir, Tangerang Selatan",
			PhoneNumber: "021 - 27599921",
			FaxNumber:   "021 - 7414934",
		},
		company.Company{
			CompanyName: "PT Berkat Raya Optima",
			Province:    "Jakarta Selatan",
			Address:     "Gedung Graha Iskandarsyah, Jl. Iskandarsyah Raya No. 66 C Desa Melawai, Kecamatan Kebayoran Baru",
		},
		company.Company{
			CompanyName: "PT PLN Batubara",
			Province:    "Jakarta Selatan",
			Address:     "Jl. Warung Buncit Raya No. 10 Pancoran, Jakarta Selatan",
			PhoneNumber: "021 - 29122118",
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
	emailDeliNiagaSejahtera := "salim@deli.id"
	emailSolusiBangunIndonesia := "budi.ariapratama@sig.id"
	emailTantraMiningDevelopment := "traffic.operationtmd@gmail.com"
	emailDaidanAditamaYaksa := "ops.traffic@daidanaditama.com"
	emailAnagaAbyudayaAnanta := "ops.traffic@anagaabyudaya.com"
	emailObsidianStainlessSteel := "julia.putri@oss.co.id"
	emailVirtueDragonNickelIndustry := "marsha.rohimone@vdni.co.id"
	emailDeliNiagaJaya := "trafficoperationdnj@gmail.com"
	emailGeoMineralTrading := "trafficoperartiongmt@gmail.com"
	emailGunbusterNickelIndustry := "laurench.gni@gmail.com"
	emailKalderaEnergiNusantara := "info@kalderacoal.co.id"
	emailBerkatRayaOptima := "berkatkaryaoptima77@gmail.com"
	emailPlnBatubara := "adhit@plnbatubara.co.id"

	createTrader = append(createTrader,
		trader.Trader{
			TraderName: "Richard NM Palar",
			Email:      &emailAJE,
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Angsana Jaya Energi"),
		},
		trader.Trader{
			TraderName: "Salim Limanto",
			Email:      &emailDeliNiagaSejahtera,
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Deli Niaga Sejahtera"),
		},
		trader.Trader{
			TraderName: "Fachrur Roji, S.T.",
			Position:   "SM of BU Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Semen Indonesia, Tbk"),
		},
		trader.Trader{
			TraderName: "Iin Iswarini",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Inti Mustika Karyatama"),
		},
		trader.Trader{
			TraderName: "Mariska Natasya S",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Sempurna Indra Pratama"),
		},
		trader.Trader{
			TraderName: "Thomas Haryadi",
			Position:   "Direktur Utama",
			CompanyId:  FindCompanyId(createCompany, "PT Sempurna Indra Pratama"),
		},
		trader.Trader{
			TraderName: "Budi Aria Pratama",
			Position:   "Procurement Manager - Cat. Energy",
			CompanyId:  FindCompanyId(createCompany, "PT Solusi Bangun Indonesia"),
			Email:      &emailSolusiBangunIndonesia,
		},
		trader.Trader{
			TraderName: "Yansen Andriyan",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Tantra Mining Development"),
			Email:      &emailTantraMiningDevelopment,
		},
		trader.Trader{
			TraderName: "Mahdalena",
			Position:   "Traffic Operation",
			CompanyId:  FindCompanyId(createCompany, "PT Daidan Aditama Yaksa"),
			Email:      &emailDaidanAditamaYaksa,
		},
		trader.Trader{
			TraderName: "Krisna Riany",
			Position:   "Traffic Operation",
			CompanyId:  FindCompanyId(createCompany, "PT Daidan Aditama Yaksa"),
			Email:      &emailDaidanAditamaYaksa,
		},
		trader.Trader{
			TraderName: "Mahdalena",
			Position:   "Traffic Operation",
			CompanyId:  FindCompanyId(createCompany, "PT Anaga Abyudaya Ananta"),
			Email:      &emailAnagaAbyudayaAnanta,
		},
		trader.Trader{
			TraderName: "Krisna Riany",
			Position:   "Traffic Operation",
			CompanyId:  FindCompanyId(createCompany, "PT Anaga Abyudaya Ananta"),
			Email:      &emailAnagaAbyudayaAnanta,
		},
		trader.Trader{
			TraderName: "Christya Ayu",
			Position:   "Traffic Operation",
			CompanyId:  FindCompanyId(createCompany, "PT Anaga Abyudaya Ananta"),
			Email:      &emailAnagaAbyudayaAnanta,
		},
		trader.Trader{
			TraderName: "Julia Putri",
			Position:   "Coal Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Obsidian Stainless Steel"),
			Email:      &emailObsidianStainlessSteel,
		},
		trader.Trader{
			TraderName: "Marsha Rohimone",
			Position:   "Coal Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Virtue Dragon Nickel Industry"),
			Email:      &emailVirtueDragonNickelIndustry,
		},
		trader.Trader{
			TraderName: "Anton",
			Position:   "Coal Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Deli Niaga Jaya"),
			Email:      &emailDeliNiagaJaya,
		},
		trader.Trader{
			TraderName: "Lim An Shun",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Geo Mineral Trading"),
			Email:      &emailGeoMineralTrading,
		},
		trader.Trader{
			TraderName: "Laurencya Chandra",
			Position:   "Coal Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Gunbuster Nickel Industry"),
			Email:      &emailGunbusterNickelIndustry,
		},
		trader.Trader{
			TraderName: "Wong John Juadi",
			Position:   "Direktur Utama",
			CompanyId:  FindCompanyId(createCompany, "PT Bara Indah Sinergi"),
		},
		trader.Trader{
			TraderName: "Wahyu Afandi Harun",
			Position:   "Direktur Utama",
			CompanyId:  FindCompanyId(createCompany, "PT Sinergi Mitra Investama"),
		},
		trader.Trader{
			TraderName: "Fachrur Roji, S.T.",
			Position:   "SM of BU Procurement",
			CompanyId:  FindCompanyId(createCompany, "PT Semen Gresik"),
		},
		trader.Trader{
			TraderName: "Alexander Tanuhadi",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Kaldera Energi Nusantara"),
			Email:      &emailKalderaEnergiNusantara,
		},
		trader.Trader{
			TraderName: "Ir. Muhammad Yusrizki Muliawan",
			Position:   "Direktur",
			CompanyId:  FindCompanyId(createCompany, "PT Berkat Raya Optima"),
			Email:      &emailBerkatRayaOptima,
		},
		trader.Trader{
			TraderName: "Adhitya Sapta Adhitama",
			Position:   "Vice President Tambang dan Produksi",
			CompanyId:  FindCompanyId(createCompany, "PT PLN Batubara"),
			Email:      &emailPlnBatubara,
		},
	)

	errTrader := tx.Create(&createTrader).Error

	if errTrader != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Company")
		return
	}

	deleteAjeTraderErr := tx.Where("company_id = ?", FindCompanyId(createCompany, "PT Angsana Jaya Energi")).Delete(&createTrader).Error
	deleteTantraTraderErr := tx.Where("company_id = ?", FindCompanyId(createCompany, "PT Tantra Mining Development")).Delete(&createTrader).Error

	deleteAjeErr := tx.Where("id = ?", FindCompanyId(createCompany, "PT Angsana Jaya Energi")).Delete(&createCompany).Error
	deleteTantraErr := tx.Where("id = ?", FindCompanyId(createCompany, "PT Tantra Mining Development")).Delete(&createCompany).Error

	if deleteAjeTraderErr != nil || deleteTantraTraderErr != nil || deleteAjeErr != nil || deleteTantraErr != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Company")
		return
	}

	tx.Commit()
}
