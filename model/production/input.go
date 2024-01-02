package production

type InputCreateProduction struct {
	ProductionDate string  `json:"production_date" validate:"required,DateValidation"`
	Quantity       float64 `json:"quantity" validate:"required"`
	RitaseQuantity int     `json:"ritase_quantity"`
	PitId          *uint   `json:"pit_id"`
	JettyId        *uint   `json:"jetty_id"`
}

type FilterListProduction struct {
	ProductionDateStart string
	ProductionDateEnd   string
	Quantity            string
	PitId               string
	JettyId             string
	Field               string
	Sort                string
}
