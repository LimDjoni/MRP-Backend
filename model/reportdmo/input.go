package reportdmo

type InputCreateReportDmo struct {
	Period          string `json:"period" validate:"PeriodValidation,required"`
	Transactions    []uint `json:"transactions"`
	GroupingVessels []uint `json:"grouping_vessels"`
}

type InputUpdateReportDmo struct {
	Transactions    []uint `json:"transactions"`
	GroupingVessels []uint `json:"grouping_vessels"`
}

type InputUpdateDocumentReportDmo struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortReportDmo struct {
	Quantity     float64 `json:"quantity"`
	CreatedStart string  `json:"created_start"`
	CreatedEnd   string  `json:"created_end"`
	Field        string  `json:"field"`
	Sort         string  `json:"sort"`
}

type CheckReportDmoPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
