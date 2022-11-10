package minerba

type InputCreateMinerba struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
	ListDataDn []int `json:"list_data_dn" validate:"required,min=1"`
}

type InputUpdateMinerba struct {
	ListDataDn []int `json:"list_data_dn" validate:"required,min=1"`
}

type InputUpdateDocumentMinerba struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortMinerba struct {
	Quantity float64 `json:"quantity"`
	CreatedStart string `json:"created_start"`
	CreatedEnd string `json:"created_end"`
	Field string `json:"field"`
	Sort string `json:"sort"`
}

type CheckMinerbaPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
