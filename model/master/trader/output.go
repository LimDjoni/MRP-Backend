package trader

import "ajebackend/model/master/company"

type OutputCompanyDetail struct {
	Company     company.Company `json:"company"`
	ListTraders []Trader        `json:"list_traders"`
}
