package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/trader"
)

type DataTransactionInput struct {
	ShippingDate *string `json:"shipping_date" validate:"omitempty,DateValidation"`
	Quantity float64 `json:"quantity"`
	TugboatName string `json:"tugboat_name"`
	BargeName string `json:"barge_name"`
	VesselName string `json:"vessel_name"`
	CustomerName string `json:"customer_name"`
	Seller	string `json:"seller"`
	LoadingPortName string `json:"loading_port_name"`
	LoadingPortLocation string `json:"loading_port_location"`
	UnloadingPortName string `json:"unloading_port_name"`
	UnloadingPortLocation string `json:"unloading_port_location"`
	DmoDestinationPort string `json:"dmo_destination_port"`
	SkbDate *string `json:"skb_date" validate:"omitempty,DateValidation"`
	SkbNumber string `json:"skb_number"`
	SkabDate *string `json:"skab_date" validate:"omitempty,DateValidation"`
	SkabNumber string `json:"skab_number"`
	BillOfLadingDate *string `json:"bill_of_lading_date" validate:"omitempty,DateValidation"`
	BillOfLadingNumber string `json:"bill_of_lading_number"`
	RoyaltyRate float64 `json:"royalty_rate"`
	DpRoyaltyCurrency string `json:"dp_royalty_currency"`
	DpRoyaltyPrice float64 `json:"dp_royalty_price"`
	DpRoyaltyDate *string `json:"dp_royalty_date" validate:"omitempty,DateValidation"`
	DpRoyaltyNtpn *string `json:"dp_royalty_ntpn"`
	DpRoyaltyBillingCode *string `json:"dp_royalty_billing_code"`
	DpRoyaltyTotal float64 `json:"dp_royalty_total"`
	PaymentDpRoyaltyCurrency string `json:"payment_dp_royalty_currency"`
	PaymentDpRoyaltyPrice float64 `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyDate *string `json:"payment_dp_royalty_date" validate:"omitempty,DateValidation"`
	PaymentDpRoyaltyNtpn *string `json:"payment_dp_royalty_ntpn"`
	PaymentDpRoyaltyBillingCode *string `json:"payment_dp_royalty_billing_code"`
	PaymentDpRoyaltyTotal float64 `json:"payment_dp_royalty_total"`
	LhvDate *string `json:"lhv_date" validate:"omitempty,DateValidation"`
	LhvNumber string `json:"lhv_number"`
	SurveyorName string `json:"surveyor_name"`
	CowDate *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber string `json:"cow_number"`
	CoaDate *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber string `json:"coa_number"`
	QualityTmAr float64 `json:"quality_tm_ar"`
	QualityImAdb float64 `json:"quality_im_adb"`
	QualityAshAr float64 `json:"quality_ash_ar"`
	QualityAshAdb float64 `json:"quality_ash_adb"`
	QualityVmAdb float64 `json:"quality_vm_adb"`
	QualityFcAdb float64 `json:"quality_fc_adb"`
	QualityTsAr float64 `json:"quality_ts_ar"`
	QualityTsAdb float64 `json:"quality_ts_adb"`
	QualityCaloriesAr float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb float64 `json:"quality_calories_adb"`
	BargingDistance float64 `json:"barging_distance"`
	SalesSystem string `json:"sales_system"`
	InvoiceDate *string `json:"invoice_date" validate:"omitempty,DateValidation"`
	InvoiceNumber string `json:"invoice_number"`
	InvoicePriceUnit float64 `json:"invoice_price_unit"`
	InvoicePriceTotal float64 `json:"invoice_price_total"`
	DmoReconciliationLetter string `json:"dmo_reconciliation_letter"`
	ContractDate *string `json:"contract_date" validate:"omitempty,DateValidation"`
	ContractNumber string `json:"contract_number"`
	DmoBuyerName string `json:"dmo_buyer_name"`
	DmoIndustryType string `json:"dmo_industry_type"`
	DmoCategory string `json:"dmo_category"`
	IsNotClaim	bool `json:"is_not_claim"`
}

type SortAndFilter struct {
	Field string
	Sort string
	Quantity float64
	TugboatName string
	BargeName string
	VesselName string
	ShippingFrom string
	ShippingTo string
}

type InputRequestCreateExcelMinerba struct {
	Authorization	string `json:"authorization"`
	MinerbaNumber	string `json:"minerba_number"`
	MinerbaPeriod	string `json:"minerba_period"`
	MinerbaId		int	`json:"minerba_id"`
	Transactions	[]Transaction `json:"transactions"`
}

type InputRequestCreateUploadDmo struct {
	Authorization	string `json:"authorization"`
	BastNumber	string `json:"bast_number"`
	DataDmo	dmo.Dmo `json:"data_dmo"`
	DataTransactions []Transaction `json:"data_transactions"`
	Trader	[]trader.Trader `json:"trader"`
	TraderEndUser	trader.Trader `json:"trader_end_user"`
}
