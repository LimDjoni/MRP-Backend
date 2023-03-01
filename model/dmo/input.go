package dmo

type VesselAdjustmentInput struct {
	VesselName string  `json:"vessel_name"`
	Quantity   float64 `json:"quantity"`
	Adjustment float64 `json:"adjustment"`
	BlDate     string  `json:"bl_date"`
}

type CreateDmoInput struct {
	Period           string `form:"period" json:"period" validate:"PeriodValidation,required"`
	DocumentDate     string `form:"document_date" json:"document_date" validate:"DateValidation"`
	TransactionBarge []int  `form:"transaction_barge" json:"transaction_barge"`
	GroupingVessel   []int  `form:"grouping_vessel" json:"grouping_vessel"`
	Trader           []int  `form:"trader" json:"trader"`
	EndUser          int    `form:"end_user" json:"end_user" validate:"required"`
	IsDocumentCustom bool   `form:"is_document_custom" json:"is_document_custom"`
}

type UpdateDmoInput struct {
	TransactionBarge []int `form:"transaction_barge" json:"transaction_barge"`
	GroupingVessel   []int `form:"grouping_vessel" json:"grouping_vessel"`
}

type InputUpdateDocumentDmo struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortDmo struct {
	Quantity float64 `json:"quantity"`
	BuyerId  string  `json:"buyer_id"`
	Month    string  `json:"month"`
	Year     string  `json:"year"`
	Field    string  `json:"field"`
	Sort     string  `json:"sort"`
}
