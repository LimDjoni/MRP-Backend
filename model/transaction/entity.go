package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	DmoId *uint `json:"dmo_id"`
	Dmo *dmo.Dmo `json:"dmo"`
	MinerbaId *uint `json:"minerba_id"`
	Minerba *minerba.Minerba `json:"minerba"`
	IdNumber string `json:"id_number" gorm:"UNIQUE"`
	TransactionType string `json:"transaction_type"`
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
	DpRoyaltyDate *string `json:"dp_royalty_date" gorm:"type:DATE"`
	DpRoyaltyNtpn *string `json:"dp_royalty_ntpn" gorm:"UNIQUE"`
	DpRoyaltyBillingCode *string `json:"dp_royalty_billing_code" gorm:"UNIQUE"`
	DpRoyaltyTotal float64 `json:"dp_royalty_total"`
	PaymentDpRoyaltyCurrency string `json:"payment_dp_royalty_currency"`
	PaymentDpRoyaltyDate *string `json:"payment_dp_royalty_date" gorm:"type:DATE"`
	PaymentDpRoyaltyNtpn *string `json:"payment_dp_royalty_ntpn" gorm:"UNIQUE"`
	PaymentDpRoyaltyBillingCode *string `json:"payment_dp_royalty_billing_code" gorm:"UNIQUE"`
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
	DmoCategory string `json:"dmo_category"`
	SkbDocumentLink string `json:"skb_document_link"`
	SkabDocumentLink string `json:"skab_document_link"`
	BLDocumentLink string `json:"bl_document_link"`
	RoyaltiProvisionDocumentLink string `json:"royalti_provision_document_link"`
	RoyaltiFinalDocumentLink string `json:"royalti_final_document_link"`
	COWDocumentLink string `json:"cow_document_link"`
	COADocumentLink string `json:"coa_document_link"`
	InvoiceAndContractDocumentLink string `json:"invoice_and_contract_document_link"`
	LHVDocumentLink string `json:"lhv_document_link"`
}

type InATrade struct {
	gorm.Model
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
}
