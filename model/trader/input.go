package trader

type InputTrader struct {
	TraderName  string 	`json:"trader_name" validate:"required"`
	Position    string 	`json:"position" validate:"required"`
	CompanyId	uint	`json:"company_id" validate:"required"`
}
