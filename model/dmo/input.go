package dmo

type VesselAdjustmentInput struct {
	VesselName string	`json:"vessel_name"`
	Quantity	float64 `json:"quantity"`
	Adjustment 	float64 `json:"adjustment"`
}

type CreateDmoInput struct {
	Period	string `json:"period" validate:"PeriodValidation,required"`
	TransactionBarge []int `json:"transaction_barge"`
	TransactionVessel []int `json:"transaction_vessel"`
	Trader []int `json:"trader" validate:"required,min=1"`
	EndUser	int `json:"end_user" validate:"required"`
	VesselAdjustment []VesselAdjustmentInput `json:"vessel_adjustment"`
	IsDocumentCustom bool `json:"is_document_custom"`
}

type InputUpdateDocumentDmo struct {
	Data []map[string]interface{} `json:"data"`
}
