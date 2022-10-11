package trader

import "ajebackend/model/company"

type OutputCompanyDetail struct {
	Company company.Company `json:"company"`
	ListTraders []Trader `json:"list_traders"`
}
