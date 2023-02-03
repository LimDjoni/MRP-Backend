package reportdmo

type InputCreateReportDmo struct {
	Period          string `json:"period" validate:"PeriodValidation,required"`
	Transactions    []uint `json:"transactions"`
	GroupingVessels []uint `json:"grouping_vessels"`
}

type InputUpdateDocumentReportDmo struct {
	Data []map[string]interface{} `json:"data"`
}
