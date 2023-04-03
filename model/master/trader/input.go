package trader

type InputCreateUpdateTrader struct {
	TraderName  string 	`json:"trader_name" validate:"required"`
	Position    string 	`json:"position" validate:"required"`
	Email		*string `json:"email" validate:"omitempty,email"`
	CompanyId	int	`json:"company_id" validate:"required"`
}
