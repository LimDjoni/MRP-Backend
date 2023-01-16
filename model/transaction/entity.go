package transaction

import (
	"ajebackend/model/destination"
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	DmoId                          *uint                              `json:"dmo_id"`
	Dmo                            *dmo.Dmo                           `json:"dmo" gorm:"constraint:OnDelete:SET NULL;"`
	MinerbaId                      *uint                              `json:"minerba_id"`
	Minerba                        *minerba.Minerba                   `json:"minerba" gorm:"constraint:OnDelete:SET NULL;"`
	IdNumber                       *string                            `json:"id_number" gorm:"UNIQUE"`
	TransactionType                string                             `json:"transaction_type"`
	ShippingDate                   *string                            `json:"shipping_date" gorm:"type:DATE"`
	Quantity                       float64                            `json:"quantity"`
	TugboatName                    string                             `json:"tugboat_name"`
	BargeName                      string                             `json:"barge_name"`
	VesselName                     string                             `json:"vessel_name"`
	Seller                         string                             `json:"seller"`
	CustomerName                   string                             `json:"customer_name"`
	LoadingPortName                string                             `json:"loading_port_name"`
	LoadingPortLocation            string                             `json:"loading_port_location"`
	UnloadingPortName              string                             `json:"unloading_port_name"`
	UnloadingPortLocation          string                             `json:"unloading_port_location"`
	DmoDestinationPort             string                             `json:"dmo_destination_port"`
	SkbDate                        *string                            `json:"skb_date" gorm:"type:DATE"`
	SkbNumber                      string                             `json:"skb_number"`
	SkabDate                       *string                            `json:"skab_date" gorm:"type:DATE"`
	SkabNumber                     string                             `json:"skab_number"`
	BillOfLadingDate               *string                            `json:"bill_of_lading_date" gorm:"type:DATE"`
	BillOfLadingNumber             string                             `json:"bill_of_lading_number"`
	RoyaltyRate                    float64                            `json:"royalty_rate"`
	DpRoyaltyPrice                 float64                            `json:"dp_royalty_price"`
	DpRoyaltyCurrency              string                             `json:"dp_royalty_currency"`
	DpRoyaltyDate                  *string                            `json:"dp_royalty_date" gorm:"type:DATE"`
	DpRoyaltyNtpn                  *string                            `json:"dp_royalty_ntpn" gorm:"UNIQUE"`
	DpRoyaltyBillingCode           *string                            `json:"dp_royalty_billing_code" gorm:"UNIQUE"`
	DpRoyaltyTotal                 float64                            `json:"dp_royalty_total"`
	PaymentDpRoyaltyPrice          float64                            `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyCurrency       string                             `json:"payment_dp_royalty_currency"`
	PaymentDpRoyaltyDate           *string                            `json:"payment_dp_royalty_date" gorm:"type:DATE"`
	PaymentDpRoyaltyNtpn           *string                            `json:"payment_dp_royalty_ntpn" gorm:"UNIQUE"`
	PaymentDpRoyaltyBillingCode    *string                            `json:"payment_dp_royalty_billing_code" gorm:"UNIQUE"`
	PaymentDpRoyaltyTotal          float64                            `json:"payment_dp_royalty_total"`
	LhvDate                        *string                            `json:"lhv_date" gorm:"type:DATE"`
	LhvNumber                      string                             `json:"lhv_number"`
	SurveyorName                   string                             `json:"surveyor_name"`
	CowDate                        *string                            `json:"cow_date" gorm:"type:DATE"`
	CowNumber                      string                             `json:"cow_number"`
	CoaDate                        *string                            `json:"coa_date" gorm:"type:DATE"`
	CoaNumber                      string                             `json:"coa_number"`
	QualityTmAr                    float64                            `json:"quality_tm_ar"`
	QualityImAdb                   float64                            `json:"quality_im_adb"`
	QualityAshAr                   float64                            `json:"quality_ash_ar"`
	QualityAshAdb                  float64                            `json:"quality_ash_adb"`
	QualityVmAdb                   float64                            `json:"quality_vm_adb"`
	QualityFcAdb                   float64                            `json:"quality_fc_adb"`
	QualityTsAr                    float64                            `json:"quality_ts_ar"`
	QualityTsAdb                   float64                            `json:"quality_ts_adb"`
	QualityCaloriesAr              float64                            `json:"quality_calories_ar"`
	QualityCaloriesAdb             float64                            `json:"quality_calories_adb"`
	BargingDistance                float64                            `json:"barging_distance"`
	SalesSystem                    string                             `json:"sales_system"`
	InvoiceDate                    *string                            `json:"invoice_date" gorm:"type:DATE"`
	InvoiceNumber                  string                             `json:"invoice_number"`
	InvoicePriceUnit               float64                            `json:"invoice_price_unit"`
	InvoicePriceTotal              float64                            `json:"invoice_price_total"`
	DmoReconciliationLetter        string                             `json:"dmo_reconciliation_letter"`
	ContractDate                   *string                            `json:"contract_date" gorm:"type:DATE"`
	ContractNumber                 string                             `json:"contract_number"`
	DmoBuyerName                   string                             `json:"dmo_buyer_name"`
	DmoIndustryType                string                             `json:"dmo_industry_type"`
	DmoCategory                    string                             `json:"dmo_category"`
	SkbDocumentLink                string                             `json:"skb_document_link"`
	SkabDocumentLink               string                             `json:"skab_document_link"`
	BLDocumentLink                 string                             `json:"bl_document_link"`
	RoyaltiProvisionDocumentLink   string                             `json:"royalti_provision_document_link"`
	RoyaltiFinalDocumentLink       string                             `json:"royalti_final_document_link"`
	COWDocumentLink                string                             `json:"cow_document_link"`
	COADocumentLink                string                             `json:"coa_document_link"`
	InvoiceAndContractDocumentLink string                             `json:"invoice_and_contract_document_link"`
	LHVDocumentLink                string                             `json:"lhv_document_link"`
	IsNotClaim                     bool                               `json:"is_not_claim"`
	IsMigration                    bool                               `json:"is_migration"`
	IsFinanceCheck                 bool                               `json:"is_finance_check"`
	IsCoaFinish                    bool                               `json:"is_coa_finish"`
	IsRoyaltyFinalFinish           bool                               `json:"is_royalty_final_finish"`
	MinerbaLnId                    *uint                              `json:"minerba_ln_id"`
	MinerbaLn                      *minerbaln.MinerbaLn               `json:"minerba_ln" gorm:"constraint:OnDelete:SET NULL;"`
	GroupingVesselLnId             *uint                              `json:"grouping_vessel_ln_id"`
	GroupingVesselLn               *groupingvesselln.GroupingVesselLn `json:"grouping_vessel_ln" gorm:"constraint:OnDelete:SET NULL;"`
	DestinationId                  *uint                              `json:"destination_id"`
	Destination                    *destination.Destination           `json:"destination"`
	GroupingVesselDnId             *uint                              `json:"grouping_vessel_dn_id"`
	GroupingVesselDn               *groupingvesseldn.GroupingVesselDn `json:"grouping_vessel_dn" gorm:"constraint:OnDelete:SET NULL;"`
}
