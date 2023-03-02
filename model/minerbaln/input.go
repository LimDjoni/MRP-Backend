package minerbaln

type InputCreateMinerbaLn struct {
	Period     string `json:"period" validate:"PeriodValidation,required"`
	ListDataLn []int  `json:"list_data_ln" validate:"required,min=1"`
}

type InputUpdateMinerbaLn struct {
	ListDataLn []int `json:"list_data_ln" validate:"required,min=1"`
}

type InputUpdateDocumentMinerbaLn struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortMinerbaLn struct {
	Quantity     string `json:"quantity"`
	UpdatedStart string `json:"updated_start"`
	UpdatedEnd   string `json:"updated_end"`
	Field        string `json:"field"`
	Sort         string `json:"sort"`
	Month        string `json:"month"`
	Year         string `json:"year"`
}

type CheckMinerbaLnPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
