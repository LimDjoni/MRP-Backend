package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/master/trader"
	"ajebackend/model/minerba"
	"ajebackend/model/reportdmo"
)

type DataTransactionInput struct {
	ShippingDate                *string `json:"shipping_date" validate:"DateValidation"`
	Quantity                    float64 `json:"quantity"`
	TugboatId                   *uint   `json:"tugboat_id"`
	BargeId                     *uint   `json:"barge_id"`
	VesselName                  string  `json:"vessel_name"`
	VesselId                    *uint   `json:"vessel_id"`
	CustomerId                  *uint   `json:"customer_id"`
	LoadingPortId               *uint   `json:"loading_port_id"`
	UnloadingPortId             *uint   `json:"unloading_port_id"`
	DmoDestinationPortId        *uint   `json:"dmo_destination_port_id"`
	DmoDestinationPortLnName    string  `json:"dmo_destination_port_ln_name"`
	SkbDate                     *string `json:"skb_date" validate:"omitempty,DateValidation"`
	SkbNumber                   string  `json:"skb_number"`
	SkabDate                    *string `json:"skab_date" validate:"omitempty,DateValidation"`
	SkabNumber                  string  `json:"skab_number"`
	BillOfLadingDate            *string `json:"bill_of_lading_date" validate:"omitempty,DateValidation"`
	BillOfLadingNumber          string  `json:"bill_of_lading_number"`
	RoyaltyRate                 float64 `json:"royalty_rate"`
	DpRoyaltyPrice              float64 `json:"dp_royalty_price"`
	DpRoyaltyDate               *string `json:"dp_royalty_date" validate:"omitempty,DateValidation"`
	DpRoyaltyNtpn               *string `json:"dp_royalty_ntpn"`
	DpRoyaltyBillingCode        *string `json:"dp_royalty_billing_code"`
	DpRoyaltyTotal              float64 `json:"dp_royalty_total"`
	PaymentDpRoyaltyPrice       float64 `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyDate        *string `json:"payment_dp_royalty_date" validate:"omitempty,DateValidation"`
	PaymentDpRoyaltyNtpn        *string `json:"payment_dp_royalty_ntpn"`
	PaymentDpRoyaltyBillingCode *string `json:"payment_dp_royalty_billing_code"`
	PaymentDpRoyaltyTotal       float64 `json:"payment_dp_royalty_total"`
	LhvDate                     *string `json:"lhv_date" validate:"omitempty,DateValidation"`
	LhvNumber                   string  `json:"lhv_number"`
	SurveyorId                  *uint   `json:"surveyor_id"`
	CowDate                     *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber                   string  `json:"cow_number"`
	CoaDate                     *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber                   string  `json:"coa_number"`
	QualityTmAr                 float64 `json:"quality_tm_ar"`
	QualityImAdb                float64 `json:"quality_im_adb"`
	QualityAshAr                float64 `json:"quality_ash_ar"`
	QualityAshAdb               float64 `json:"quality_ash_adb"`
	QualityVmAdb                float64 `json:"quality_vm_adb"`
	QualityFcAdb                float64 `json:"quality_fc_adb"`
	QualityTsAr                 float64 `json:"quality_ts_ar"`
	QualityTsAdb                float64 `json:"quality_ts_adb"`
	QualityCaloriesAr           float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb          float64 `json:"quality_calories_adb"`
	BargingDistance             float64 `json:"barging_distance"`
	SalesSystemId               *uint   `json:"sales_system_id"`
	InvoiceDate                 *string `json:"invoice_date" validate:"omitempty,DateValidation"`
	InvoiceNumber               string  `json:"invoice_number"`
	InvoicePriceUnit            float64 `json:"invoice_price_unit"`
	InvoicePriceTotal           float64 `json:"invoice_price_total"`
	ContractDate                *string `json:"contract_date" validate:"omitempty,DateValidation"`
	ContractNumber              string  `json:"contract_number"`
	DmoBuyerId                  *uint   `json:"dmo_buyer_id"`
	IsNotClaim                  bool    `json:"is_not_claim"`
	IsFinanceCheck              bool    `json:"is_finance_check"`
	IsCoaFinish                 bool    `json:"is_coa_finish"`
	IsRoyaltyFinalFinish        bool    `json:"is_royalty_final_finish"`
	DestinationId               *uint   `json:"destination_id"`
	DestinationCountryId        *uint   `json:"destination_country_id"`
}

type SortAndFilter struct {
	Field              string
	Sort               string
	Quantity           string
	TugboatId          string
	BargeId            string
	VesselId           string
	ShippingStart      string
	ShippingEnd        string
	VerificationFilter string
}

type InputRequestCreateExcelMinerba struct {
	Authorization string          `json:"authorization"`
	MinerbaNumber string          `json:"minerba_number"`
	MinerbaPeriod string          `json:"minerba_period"`
	MinerbaId     int             `json:"minerba_id"`
	Transactions  []Transaction   `json:"transactions"`
	Minerba       minerba.Minerba `json:"minerba"`
}

type InputRequestCreateUploadDmo struct {
	Authorization                 string                              `json:"authorization"`
	BastNumber                    string                              `json:"bast_number"`
	DataDmo                       dmo.Dmo                             `json:"data_dmo"`
	Trader                        []trader.Trader                     `json:"trader"`
	TraderEndUser                 trader.Trader                       `json:"trader_end_user"`
	ListTransactionBarge          []Transaction                       `json:"list_transaction_barge"`
	ListTransactionGroupingVessel []Transaction                       `json:"list_transaction_grouping_vessel"`
	ListGroupingVessel            []groupingvesseldn.GroupingVesselDn `json:"list_grouping_vessel"`
}

type InputRequestGetReport struct {
	ProductionPlan                 float64 `json:"production_plan" validate:"required"`
	PercentageProductionObligation float64 `json:"percentage_production_obligation" validate:"required"`
	Year                           int     `json:"year"`
}

type InputRequestCreateReportDmo struct {
	Authorization   string                              `json:"authorization"`
	ReportDmo       reportdmo.ReportDmo                 `json:"report_dmo"`
	Transactions    []Transaction                       `json:"transactions"`
	GroupingVessels []groupingvesseldn.GroupingVesselDn `json:"grouping_vessels"`
}
