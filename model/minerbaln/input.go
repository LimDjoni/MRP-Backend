package minerbaln

type InputCreateMinerbaLn struct {
	Period     string `json:"period" validate:"PeriodValidation,required"`
	ListDataLn []int  `json:"list_data_ln" validate:"required,min=1"`
}

type InputUpdateMinerbaLn struct {
	ListDataLn []int `json:"list_data_ln" validate:"required,min=1"`
}

type InputUpdateDocumentMinerbaLn struct {
	Data map[string]interface{} `json:"data"`
}

type FilterAndSortMinerbaLn struct {
	Quantity     float64 `json:"quantity"`
	CreatedStart string  `json:"created_start"`
	CreatedEnd   string  `json:"created_end"`
	Field        string  `json:"field"`
	Sort         string  `json:"sort"`
}

type CheckMinerbaLnPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
