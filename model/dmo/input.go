package dmo

import "ajebackend/model/trader"

type VesselAdjustmentInput struct {
	VesselName string	`json:"vessel_name"`
	Quantity	float64 `json:"quantity"`
	Adjustment 	float64 `json:"adjustment"`
}

type CreateDmoInput struct {
	TransactionBarge []int `json:"transaction_barge"`
	TransactionVessel []int `json:"transaction_vessel"`
	Trader []trader.Trader `json:"trader"`
	EndUser	string `json:"end_user"`
	VesselAdjustment []VesselAdjustmentInput `json:"vessel_adjustment"`
}
