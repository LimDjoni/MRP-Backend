package model

import (
	"github.com/jackc/pgtype"
	"time"
)

// All add gorm.Model to create id, createdAt, deletedAt, updatedAt

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
}

type Trader struct {
	ParentId uint	`json:"parent_id"`
	Parent *Trader `gorm:"association_jointable_foreignkey:parent_id"`
	TraderName string `json:"trader_name"`
	Position string `json:"position"`
	Address string `json:"address"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber string `json:"fax_number"`
}

type TraderDmo struct {
	DmoId uint `json:"dmo_id"`
	Dmo Dmo
	TraderId uint `json:"trader_id"`
	Trader Trader
}

type Dmo struct {
	IdNumber string `json:"id_number"`
	Date time.Time `json:"date"`
	Quantity float64 `json:"quantity"`
	Adjustment float64 `json:"adjustment"`
	GrandTotalQuantity float64 `json:"grand_total_quantity"`
	EndUser pgtype.JSONB `json:"end_user"`
	DataTrader pgtype.JSONB `json:"data_trader"`
	BeritaAcaraDocument string `json:"berita_acara_document"`
}

type Minerba struct {
	IdNumber string `json:"id_number"`
	Date time.Time `json:"date"`
	Periode string `json:"periode"`
	SP3MEDNDocument string `json:"sp3medn_document"`
	RekapDmoDocument string `json:"rekap_dmo_document"`
	RincianDmoDocument string `json:"rincian_dmo_document"`
	SP3MELNDocument string `json:"sp3meln_document"`
	INSWEksporDocument string `json:"insw_ekspor_document"`
}

type MinerbaTransaction struct {
	MinerbaId uint `json:"minerba_id"`
	Minerba Minerba
	TransactionId uint `json:"transaction_id"`
	Transaction Transaction
}

type Transaction struct {
	DmoId uint `json:"dmo_id"`
	IdNumber string `json:"id_number"`
	TransactionType string `json:"transaction_type"`
	Number int `json:"number"`
	ShippingDate time.Time `json:"shipping_date"`
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
	SkbDate time.Time `json:"skb_date"`
	SkbNumber string `json:"skb_number"`
	SkabDate time.Time `json:"skab_date"`
	SkabNumber string `json:"skab_number"`
	BillOfLadingDate time.Time `json:"bill_of_lading_date"`
	BillOfLadingNumber string `json:"bill_of_lading_number"`
	RoyaltyRate float64 `json:"royalty_rate"`
	DpRoyaltyPrice float64 `json:"dp_royalty_price"`
	DpRoyaltyDate time.Time `json:"dp_royalty_date"`
	DpRoyaltyNtpn string `json:"dp_royalty_ntpn"`
	DpRoyaltyBillingCode string `json:"dp_royalty_billing_code"`
	DpRoyaltyTotal float64 `json:"dp_royalty_total"`
	PaymentDpRoyaltyPrice float64 `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyDate time.Time `json:"payment_dp_royalty_date"`
	PaymentDpRoyaltyNtpn string `json:"payment_dp_royalty_ntpn"`
	PaymentDpRoyaltyBillingCode string `json:"payment_dp_royalty_billing_code"`
	PaymentDpRoyaltyTotal float64 `json:"payment_dp_royalty_total"`
	LhvDate time.Time `json:"lhv_date"`
	LhvNumber string `json:"lhv_number"`
	SurveyorName string `json:"surveyor_name"`
	CowDate time.Time `json:"cow_date"`
	CowNumber string `json:"cow_number"`
	CoaDate time.Time `json:"coa_date"`
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
	InvoiceDate time.Time `json:"invoice_date"`
	InvoiceNumber string `json:"invoice_number"`
	InvoicePriceUnit float64 `json:"invoice_price_unit"`
	InvoicePriceTotal float64 `json:"invoice_price_total"`
	DmoReconciliationLetter string `json:"dmo_reconciliation_letter"`
	ContractDate time.Time `json:"contract_date"`
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
	DataLsExportDate time.Time `json:"data_ls_export_date"`
	DataLsExportNumber string `json:"data_ls_export_number"`
	DataSkaCooDate time.Time `json:"data_ska_coo_date"`
	DataSkaCooNumber string `json:"data_ska_coo_number"`
	PebNumber string `json:"peb_number"`
	PebDate time.Time `json:"peb_date"`
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

type History struct {
	DmoId uint `json:"dmo_id"`
	Dmo Dmo
	TransactionId uint `json:"transaction_id"`
	Transaction Transaction
	UserId uint `json:"user_id"`
	User User
	MinerbaId uint `json:"minerba_id"`
	Minerba Minerba
	Status string `json:"status"`
}


