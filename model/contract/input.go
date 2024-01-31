package contract

type InputCreateUpdateContract struct {
	ContractDate   string  `form:"contract_date" json:"contract_date" gorm:"DATETIME"`
	ContractNumber string  `json:"contract_number" gorm:"UNIQUE"`
	CustomerId     uint    `form:"customer_id" json:"customer_id"`
	Quantity       float64 `form:"quantity" json:"quantity"`
	Validity       string  `form:"validity" json:"validity" gorm:"DATETIME"`
}

type FilterAndSortContract struct {
	ContractDateStart string `json:"contract_date_start" gorm:"DATETIME"`
	ContractDateEnd   string `json:"contract_date_end" gorm:"DATETIME"`
	ContractNumber    string `json:"contract_number" gorm:"UNIQUE"`
	CustomerId        string `json:"customer_id"`
	Quantity          string `json:"quantity"`
	ValidityStart     string `json:"validity_start" gorm:"DATETIME"`
	ValidityEnd       string `json:"validity_end" gorm:"DATETIME"`
	Field             string `json:"field"`
	Sort              string `json:"sort"`
}
