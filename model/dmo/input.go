package dmo

type VesselAdjustmentInput struct {
	VesselName string	`json:"vessel_name"`
	Quantity	float64 `json:"quantity"`
	Adjustment 	float64 `json:"adjustment"`
}

type CreateDmoInput struct {
	Period	string `form:"period" json:"period" validate:"PeriodValidation,required"`
	TransactionBarge []int `form:"transaction_barge" json:"transaction_barge"`
	TransactionVessel []int `form:"transaction_vessel" json:"transaction_vessel"`
	Trader []int `form:"trader" json:"trader" validate:"required,min=1"`
	EndUser	int `form:"end_user" json:"end_user" validate:"required"`
	VesselAdjustment []VesselAdjustmentInput `form:"vessel_adjustment" json:"vessel_adjustment"`
	IsDocumentCustom bool `form:"is_document_custom" json:"is_document_custom"`
}

type InputUpdateDocumentDmo struct {
	Data []map[string]interface{} `json:"data"`
}

type FilterAndSortDmo struct {
	Quantity float64 `json:"quantity"`
	CreatedStart string `json:"created_start"`
	CreatedEnd string `json:"created_end"`
	Field string  `json:"field"`
	Sort string `json:"sort"`
}
