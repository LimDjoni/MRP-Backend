package production

type InputCreateProduction struct {
	ProductionDate string `json:"production_date" validate:"required,DateValidation"`
	Quantity float64 `json:"quantity" validate:"required"`
}

type FilterListProduction struct {
	ProductionDateStart string
	ProductionDateEnd string
	Quantity float64
}
