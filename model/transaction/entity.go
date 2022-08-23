package transaction

import (
	"ajebackend/model/dmo"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	DmoId *uint `json:"dmo_id"`
	Dmo *dmo.Dmo `json:"dmo"`
	IdNumber string `json:"id_number" gorm:"UNIQUE"`
	TransactionType string `json:"transaction_type"`
	Number int `json:"number"`
	ShippingDate *string `json:"shipping_date" gorm:"type:DATE"`
	Quantity float64 `json:"quantity"`
	ShipName string `json:"ship_name"`
	BargeName string `json:"barge_name"`
	VesselName string `json:"vessel_name"`
	Seller string `json:"seller"`
	CustomerName string `json:"customer_name"`
	LoadingPortName string `json:"loading_port_name"`
	LoadingPortLocation string `json:"loading_port_location"`
	UnloadingPortName string `json:"unloading_port_name"`
	UnloadingPortLocation string `json:"unloading_port_location"`
	DmoDestinationPort string `json:"dmo_destination_port"`
	SkbDate *string `json:"skb_date" gorm:"type:DATE"`
	SkbNumber string `json:"skb_number"`
	SkabDate *string `json:"skab_date" gorm:"type:DATE"`
	SkabNumber string `json:"skab_number"`
	BillOfLadingDate *string `json:"bill_of_lading_date" gorm:"type:DATE"`
	BillOfLadingNumber string `json:"bill_of_lading_number"`
	RoyaltyRate float64 `json:"royalty_rate"`
	DpRoyaltyCurrency string `json:"dp_royalty_currency"`
	DpRoyaltyPrice float64 `json:"dp_royalty_price"`
	DpRoyaltyDate *string `json:"dp_royalty_date" gorm:"type:DATE"`
	DpRoyaltyNtpn string `json:"dp_royalty_ntpn"`
	DpRoyaltyBillingCode string `json:"dp_royalty_billing_code"`
	DpRoyaltyTotal float64 `json:"dp_royalty_total"`
	PaymentDpRoyaltyCurrency string `json:"payment_dp_royalty_currency"`
	PaymentDpRoyaltyPrice float64 `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyDate *string `json:"payment_dp_royalty_date" gorm:"type:DATE"`
	PaymentDpRoyaltyNtpn string `json:"payment_dp_royalty_ntpn"`
	PaymentDpRoyaltyBillingCode string `json:"payment_dp_royalty_billing_code"`
	PaymentDpRoyaltyTotal float64 `json:"payment_dp_royalty_total"`
	LhvDate *string `json:"lhv_date" gorm:"type:DATE"`
	LhvNumber string `json:"lhv_number"`
	SurveyorName string `json:"surveyor_name"`
	CowDate *string `json:"cow_date" gorm:"type:DATE"`
	CowNumber string `json:"cow_number"`
	CoaDate *string `json:"coa_date" gorm:"type:DATE"`
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
	InvoiceDate *string `json:"invoice_date" gorm:"type:DATE"`
	InvoiceNumber string `json:"invoice_number"`
	InvoicePriceUnit float64 `json:"invoice_price_unit"`
	InvoicePriceTotal float64 `json:"invoice_price_total"`
	DmoReconciliationLetter string `json:"dmo_reconciliation_letter"`
	ContractDate *string `json:"contract_date" gorm:"type:DATE"`
	ContractNumber string `json:"contract_number"`
	DmoBuyerName string `json:"dmo_buyer_name"`
	DmoIndustryType string `json:"dmo_industry_type"`
	DmoStatusReconciliationLetter string `json:"dmo_status_reconciliation_letter"`
	InATradeNumber string `json:"in_a_trade_number"`
	DescriptionOfGood string `json:"description_of_good"`
	TarifPosHs string `json:"tarif_pos_hs"`
	Volume float64 `json:"volume"`
	Unit string `json:"unit"`
	Value float64 `json:"value"`
	Currency string `json:"currency"`
	InATradeLoadingPort string `json:"in_a_trade_loading_port"`
	DestinationCountry string `json:"destination_country"`
	DataLsExportDate *string `json:"data_ls_export_date" gorm:"type:DATE"`
	DataLsExportNumber string `json:"data_ls_export_number"`
	DataSkaCooDate *string `json:"data_ska_coo_date" gorm:"type:DATE"`
	DataSkaCooNumber string `json:"data_ska_coo_number"`
	PebNumber string `json:"peb_number"`
	PebDate *string `json:"peb_date" gorm:"type:DATE"`
	AjuNumber string `json:"aju_number"`
	DwtValue float64 `json:"dwt_value"`
	InsuranceCompanyName string `json:"insurance_company_name"`
	InsurancePolisNumber string `json:"insurance_polis_number"`
	NavyShipName string `json:"navy_ship_name"`
	NavyCompanyName string `json:"navy_company_name"`
	NavyImoNumber string `json:"navy_imo_number"`
	SkbDocument string `json:"skb_document"`
	SkabDocument string `json:"skab_document"`
	BLDocument string `json:"bl_document"`
	RoyaltiProvisionDocument string `json:"royalti_provision_document"`
	RoyaltiFinalDocument string `json:"royalti_final_document"`
	COWDocument string `json:"cow_document"`
	COADocument string `json:"coa_document"`
	InvoiceAndContractDocument string `json:"invoice_and_contract_document"`
	LHVDocument string `json:"lhv_document"`
}
