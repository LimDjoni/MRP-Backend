package fuelratio

type RegisterFuelRatioInput struct {
	UnitId       uint    `json:"unit_id" validate:"required"`
	OperatorName string  `json:"operator_name" validate:"required"`
	Shift        string  `json:"shift"`
	Tanggal      string  `json:"tanggal"`
	FirstHM      float64 `json:"first_hm"`
	LastHM       float64 `json:"last_hm"`
	TanggalAwal  string  `json:"tanggal_awal" gorm:"DATETIME"`
	TanggalAkhir string  `json:"tanggal_akhir" gorm:"DATETIME"`
	TotalRefill  uint    `json:"total_refill"`
	Status       bool    `json:"status"`
}

type SortFilterFuelRatio struct {
	Field        string
	Sort         string
	UnitId       string
	OperatorName string
	Shift        string
	FirstHM      string
	Status       string
}

type SortFilterFuelRatioSummary struct {
	Field            string
	Sort             string
	UnitID           string `gorm:"column:unit_id"`
	UnitName         string `gorm:"column:unit_name"`
	Shift            string `gorm:"column:shift"`
	TotalRefill      string `gorm:"column:total_refill"`
	Consumption      string `gorm:"column:consumption"`
	Tolerance        string `gorm:"column:tolerance"`
	Tanggal          string `gorm:"column:tanggal"`
	TanggalAwal      string `gorm:"column:tanggal_awal"`
	TanggalAkhir     string `gorm:"column:tanggal_akhir"`
	FirstHM          string
	LastHM           string
	Duration         string  `gorm:"column:duration"`
	BatasBawah       float64 `gorm:"column:batas_bawah"`
	BatasAtas        float64 `gorm:"column:batas_atas"`
	TotalKonsumsiBBM string  `gorm:"column:total_konsumsi_bbm"`
}
