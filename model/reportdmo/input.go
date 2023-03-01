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
	UpdatedStart string  `json:"updated_start"`
	UpdatedEnd   string  `json:"updated_end"`
	Month        string  `json:"month"`
	Year         string  `json:"year"`
	Field        string  `json:"field"`
	Sort         string  `json:"sort"`
}

type CheckReportDmoPeriod struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
}
