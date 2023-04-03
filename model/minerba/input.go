package minerba

type InputCreateMinerba struct {
	Period     string `json:"period" validate:"PeriodValidation,required"`
	ListDataDn []int  `json:"list_data_dn" validate:"required,min=1"`
}

type InputUpdateMinerba struct {
	ListDataDn []int `json:"list_data_dn" validate:"required,min=1"`
}

type InputUpdateDocumentMinerba struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortMinerba struct {
	Quantity     string `json:"quantity"`
	UpdatedStart string `json:"updated_start"`
	UpdatedEnd   string `json:"updated_end"`
	Field        string `json:"field"`
	Sort         string `json:"sort"`
	Month        string `json:"month"`
	Year         string `json:"year"`
}

type CheckMinerbaPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
