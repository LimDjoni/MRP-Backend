package fuelratio

type RegisterFuelRatioInput struct {
	UnitId      uint   `json:"unit_id" validate:"required"`
	EmployeeId  uint   `json:"employee_id" validate:"required"`
	Shift       string `json:"shift" validate:"required"`
	FirstHM     string `json:"first_hm" gorm:"DATETIME"`
	LastHM      string `json:"last_hm" gorm:"DATETIME"`
	TotalRefill uint   `json:"total_refill"`
	Status      bool   `json:"status"`
}

type SortFilterFuelRatio struct {
	Field      string
	Sort       string
	UnitId     string
	EmployeeId string
	Shift      string
	FirstHM    string
	Status     string
}

type SortFilterFuelRatioSummary struct {
	Field       string
	Sort        string
	UnitID      string `gorm:"column:unit_id"`
	UnitName    string `gorm:"column:unit_name"`
	Shift       string `gorm:"column:shift"`
	TotalRefill string `gorm:"column:total_refill"`
	Consumption string `gorm:"column:consumption"`
	Tolerance   string `gorm:"column:tolerance"`
	FirstHM     string
	LastHM      string
	Duration    string  `gorm:"column:duration"`
	BatasBawah  float64 `gorm:"column:batas_bawah"`
	BatasAtas   float64 `gorm:"column:batas_atas"`
}
