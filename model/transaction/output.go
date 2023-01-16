package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/minerba"
)

type DetailMinerba struct {
	Detail minerba.Minerba `json:"detail"`
	List   []Transaction   `json:"list"`
}

type DetailDmo struct {
	Detail dmo.Dmo       `json:"detail"`
	List   []Transaction `json:"list"`
}

type ErrorResponseUnique struct {
	FailedField string `json:"FailedField"`
	Tag         string `json:"Tag"`
	Value       string `json:"Value"`
}

type Electricity struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type NonElectricity struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type NotClaimable struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type QuantityProduction struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type ReportDetailOutput struct {
	Electricity           Electricity        `json:"electricity"`
	NonElectricity        NonElectricity     `json:"non_electricity"`
	NotClaimable          NotClaimable       `json:"not_claimable"`
	Production            QuantityProduction `json:"production"`
	ElectricityCompany    []string           `json:"electricity_company"`
	NonElectricityCompany []string           `json:"non_electricity_company"`
}

type ReportRecapOutput struct {
	ElectricityTotal                          float64 `json:"electricity_total"`
	NonElectricityTotal                       float64 `json:"non_electricity_total"`
	Total                                     float64 `json:"total"`
	TotalProduction                           float64 `json:"total_production"`
	RateCalories                              string  `json:"rate_calories"`
	ProductionPlan                            float64 `json:"production_plan"`
	ProductionObligation                      float64 `json:"production_obligation"`
	PercentageProductionObligation            float64 `json:"percentage_production_obligation"`
	ProrateProductionPlan                     string  `json:"prorate_production_plan"`
	FulfillmentOfProductionPlan               string  `json:"fulfillment_of_production_plan"`
	FulfillmentOfProductionRealization        string  `json:"fulfillment_of_production_realization"`
	FulfillmentPercentageProductionObligation string  `json:"fulfillment_percentage_production_obligation"`
	Year                                      int     `json:"year"`
}

type ChooseTransactionDmo struct {
	BargeTransaction  []Transaction `json:"barge_transaction"`
	VesselTransaction []Transaction `json:"vessel_transaction"`
}

type DetailGroupingVesselLn struct {
	ListTransactions []Transaction                     `json:"list_transactions"`
	Detail           groupingvesselln.GroupingVesselLn `json:"detail"`
}

type ListDataDnForGroupingVessel struct {
	ListDataDnBargeWithVessel []Transaction `json:"list_data_dn_barge_with_vessel"`
	ListDataDnVessel          []Transaction `json:"list_data_dn_vessel"`
}
