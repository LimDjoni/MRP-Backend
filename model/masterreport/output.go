package masterreport

import (
	"ajebackend/model/master/barge"
	"ajebackend/model/master/company"
	"ajebackend/model/master/country"
	"ajebackend/model/master/currency"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/ports"
	"ajebackend/model/master/salessystem"
	"ajebackend/model/master/surveyor"
	"ajebackend/model/master/tugboat"
	"ajebackend/model/master/vessel"
	"ajebackend/model/rkab"
)

// Report Recap DMO Output

type ReportDmoOutput struct {
	RecapElectricity    RecapElectricity    `json:"recap_electricity"`
	RecapCement         RecapCement         `json:"recap_cement"`
	RecapNonElectricity RecapNonElectricity `json:"recap_non_electricity"`
	NotClaimable        NotClaimable        `json:"not_claimable"`
	Production          QuantityProduction  `json:"production"`
	Rkabs               []rkab.Rkab         `json:"rkabs"`
}

type RecapElectricity struct {
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

type RecapCement struct {
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

type RecapNonElectricity struct {
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

// Report Realization Output

type RealizationOutput struct {
	Electric    RealizationElectric    `json:"electric"`
	Cement      RealizationCement      `json:"cement"`
	NonElectric RealizationNonElectric `json:"non_electric"`
}

type RealizationElectric struct {
	January   []RealizationTransaction `json:"january"`
	February  []RealizationTransaction `json:"february"`
	March     []RealizationTransaction `json:"march"`
	April     []RealizationTransaction `json:"april"`
	May       []RealizationTransaction `json:"may"`
	June      []RealizationTransaction `json:"june"`
	July      []RealizationTransaction `json:"july"`
	August    []RealizationTransaction `json:"august"`
	September []RealizationTransaction `json:"september"`
	October   []RealizationTransaction `json:"october"`
	November  []RealizationTransaction `json:"november"`
	December  []RealizationTransaction `json:"december"`
}

type RealizationCement struct {
	January   []RealizationTransaction `json:"january"`
	February  []RealizationTransaction `json:"february"`
	March     []RealizationTransaction `json:"march"`
	April     []RealizationTransaction `json:"april"`
	May       []RealizationTransaction `json:"may"`
	June      []RealizationTransaction `json:"june"`
	July      []RealizationTransaction `json:"july"`
	August    []RealizationTransaction `json:"august"`
	September []RealizationTransaction `json:"september"`
	October   []RealizationTransaction `json:"october"`
	November  []RealizationTransaction `json:"november"`
	December  []RealizationTransaction `json:"december"`
}

type RealizationNonElectric struct {
	January   []RealizationTransaction `json:"january"`
	February  []RealizationTransaction `json:"february"`
	March     []RealizationTransaction `json:"march"`
	April     []RealizationTransaction `json:"april"`
	May       []RealizationTransaction `json:"may"`
	June      []RealizationTransaction `json:"june"`
	July      []RealizationTransaction `json:"july"`
	August    []RealizationTransaction `json:"august"`
	September []RealizationTransaction `json:"september"`
	October   []RealizationTransaction `json:"october"`
	November  []RealizationTransaction `json:"november"`
	December  []RealizationTransaction `json:"december"`
}

type RealizationTransaction struct {
	ShippingDate      string           `json:"shipping_date"`
	Trader            *company.Company `json:"trader"`
	EndUser           *company.Company `json:"end_user"`
	QualityCaloriesAr float64          `json:"quality_calories_ar"`
	Quantity          float64          `json:"quantity"`
	IsBastOk          bool             `json:"is_bast_ok"`
}

// Report Detail Output

type SaleDetail struct {
	Rkabs               []rkab.Rkab         `json:"rkabs"`
	Electricity         Electricity         `json:"electricity"`
	Cement              Cement              `json:"cement"`
	NonElectricity      NonElectricity      `json:"non_electricity"`
	RecapElectricity    RecapElectricity    `json:"recap_electricity"`
	RecapCement         RecapCement         `json:"recap_cement"`
	RecapNonElectricity RecapNonElectricity `json:"recap_non_electricity"`
	// try new flexible
	DataDetailIndustry map[string]map[string]map[string]map[string]float64 `json:"data_detail_industry"`
	DataRecapIndustry  map[string]map[string]float64                       `json:"data_recap_industry"`
	//
	NotClaimable       NotClaimable       `json:"not_claimable"`
	Production         QuantityProduction `json:"production"`
	Domestic           Domestic           `json:"domestic"`
	Export             Export             `json:"export"`
	ElectricAssignment ElectricAssignment `json:"electric_assignment"`
	CafAssignment      CafAssignment      `json:"caf_assignment"`
	// try new flexible
	Company map[string]map[string][]string `json:"company"`
	// try new flexible
	CompanyElectricity    map[string][]string `json:"company_electricity"`
	CompanyCement         map[string][]string `json:"company_cement"`
	CompanyNonElectricity map[string][]string `json:"company_non_electricity"`
	ProductionJetty       ProductionJetty     `json:"production_jetty"`
	SalesJetty            SalesJetty          `json:"sales_jetty"`
	JettyList             []string            `json:"jetty_list"`
	JettyBalanceLoss      []JettyBalanceLoss  `json:"jetty_balance_loss"`
	LossJetty             LossJetty           `json:"lost_jetty"`
}

type Electricity struct {
	January   map[string]map[string]float64 `json:"january"`
	February  map[string]map[string]float64 `json:"february"`
	March     map[string]map[string]float64 `json:"march"`
	April     map[string]map[string]float64 `json:"april"`
	May       map[string]map[string]float64 `json:"may"`
	June      map[string]map[string]float64 `json:"june"`
	July      map[string]map[string]float64 `json:"july"`
	August    map[string]map[string]float64 `json:"august"`
	September map[string]map[string]float64 `json:"september"`
	October   map[string]map[string]float64 `json:"october"`
	November  map[string]map[string]float64 `json:"november"`
	December  map[string]map[string]float64 `json:"december"`
	Total     float64                       `json:"total"`
}

type Cement struct {
	January   map[string]map[string]float64 `json:"january"`
	February  map[string]map[string]float64 `json:"february"`
	March     map[string]map[string]float64 `json:"march"`
	April     map[string]map[string]float64 `json:"april"`
	May       map[string]map[string]float64 `json:"may"`
	June      map[string]map[string]float64 `json:"june"`
	July      map[string]map[string]float64 `json:"july"`
	August    map[string]map[string]float64 `json:"august"`
	September map[string]map[string]float64 `json:"september"`
	October   map[string]map[string]float64 `json:"october"`
	November  map[string]map[string]float64 `json:"november"`
	December  map[string]map[string]float64 `json:"december"`
	Total     float64                       `json:"total"`
}
type NonElectricity struct {
	January   map[string]map[string]float64 `json:"january"`
	February  map[string]map[string]float64 `json:"february"`
	March     map[string]map[string]float64 `json:"march"`
	April     map[string]map[string]float64 `json:"april"`
	May       map[string]map[string]float64 `json:"may"`
	June      map[string]map[string]float64 `json:"june"`
	July      map[string]map[string]float64 `json:"july"`
	August    map[string]map[string]float64 `json:"august"`
	September map[string]map[string]float64 `json:"september"`
	October   map[string]map[string]float64 `json:"october"`
	November  map[string]map[string]float64 `json:"november"`
	December  map[string]map[string]float64 `json:"december"`
	Total     float64                       `json:"total"`
}

type Domestic struct {
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

type Export struct {
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

type ElectricAssignment struct {
	Quantity            float64 `json:"quantity"`
	RealizationQuantity float64 `json:"realization_quantity"`
}

type CafAssignment struct {
	Quantity            float64 `json:"quantity"`
	RealizationQuantity float64 `json:"realization_quantity"`
}

type TransactionReport struct {
	TransactionType             string                   `json:"transaction_type"`
	ShippingDate                *string                  `json:"shipping_date" gorm:"type:DATE"`
	Quantity                    float64                  `json:"quantity"`
	QuantityUnloading           float64                  `json:"quantity_unloading"`
	TugboatId                   *uint                    `json:"tugboat_id"`
	Tugboat                     *tugboat.Tugboat         `json:"tugboat"`
	BargeId                     *uint                    `json:"barge_id"`
	Barge                       *barge.Barge             `json:"barge"`
	VesselId                    *uint                    `json:"vessel_id"`
	Vessel                      *vessel.Vessel           `json:"vessel"`
	SellerId                    *uint                    `json:"seller_id"`
	Seller                      *iupopk.Iupopk           `json:"seller"`
	CustomerId                  *uint                    `json:"customer_id"`
	Customer                    *company.Company         `json:"customer"`
	LoadingPortId               *uint                    `json:"loading_port_id"`
	LoadingPort                 *jetty.Jetty             `json:"loading_port"`
	UnloadingPortId             *uint                    `json:"unloading_port_id"`
	UnloadingPort               *ports.Port              `json:"unloading_port"`
	DmoDestinationPortId        *uint                    `json:"dmo_destination_port_id"`
	DmoDestinationPort          *ports.Port              `json:"dmo_destination_port"`
	SkbDate                     *string                  `json:"skb_date" gorm:"type:DATE"`
	SkbNumber                   string                   `json:"skb_number"`
	SkabDate                    *string                  `json:"skab_date" gorm:"type:DATE"`
	SkabNumber                  string                   `json:"skab_number"`
	BillOfLadingDate            *string                  `json:"bill_of_lading_date" gorm:"type:DATE"`
	BillOfLadingNumber          string                   `json:"bill_of_lading_number"`
	RoyaltyRate                 float64                  `json:"royalty_rate"`
	DpRoyaltyPrice              float64                  `json:"dp_royalty_price"`
	DpRoyaltyCurrencyId         *uint                    `json:"dp_royalty_currency_id"`
	DpRoyaltyCurrency           *currency.Currency       `json:"dp_royalty_currency"`
	DpRoyaltyDate               *string                  `json:"dp_royalty_date" gorm:"type:DATE"`
	DpRoyaltyNtpn               *string                  `json:"dp_royalty_ntpn" gorm:"UNIQUE"`
	DpRoyaltyBillingCode        *string                  `json:"dp_royalty_billing_code" gorm:"UNIQUE"`
	DpRoyaltyTotal              float64                  `json:"dp_royalty_total"`
	PaymentDpRoyaltyPrice       float64                  `json:"payment_dp_royalty_price"`
	PaymentDpRoyaltyCurrencyId  *uint                    `json:"payment_dp_royalty_currency_id"`
	PaymentDpRoyaltyCurrency    *currency.Currency       `json:"payment_dp_royalty_currency"`
	PaymentDpRoyaltyDate        *string                  `json:"payment_dp_royalty_date" gorm:"type:DATE"`
	PaymentDpRoyaltyNtpn        *string                  `json:"payment_dp_royalty_ntpn" gorm:"UNIQUE"`
	PaymentDpRoyaltyBillingCode *string                  `json:"payment_dp_royalty_billing_code" gorm:"UNIQUE"`
	PaymentDpRoyaltyTotal       float64                  `json:"payment_dp_royalty_total"`
	LhvDate                     *string                  `json:"lhv_date" gorm:"type:DATE"`
	LhvNumber                   string                   `json:"lhv_number"`
	SurveyorId                  *uint                    `json:"surveyor_id"`
	Surveyor                    *surveyor.Surveyor       `json:"surveyor"`
	CowDate                     *string                  `json:"cow_date" gorm:"type:DATE"`
	CowNumber                   string                   `json:"cow_number"`
	CoaDate                     *string                  `json:"coa_date" gorm:"type:DATE"`
	CoaNumber                   string                   `json:"coa_number"`
	QualityTmAr                 float64                  `json:"quality_tm_ar"`
	QualityImAdb                float64                  `json:"quality_im_adb"`
	QualityAshAr                float64                  `json:"quality_ash_ar"`
	QualityAshAdb               float64                  `json:"quality_ash_adb"`
	QualityVmAdb                float64                  `json:"quality_vm_adb"`
	QualityFcAdb                float64                  `json:"quality_fc_adb"`
	QualityTsAr                 float64                  `json:"quality_ts_ar"`
	QualityTsAdb                float64                  `json:"quality_ts_adb"`
	QualityCaloriesAr           float64                  `json:"quality_calories_ar"`
	QualityCaloriesAdb          float64                  `json:"quality_calories_adb"`
	BargingDistance             float64                  `json:"barging_distance"`
	SalesSystemId               *uint                    `json:"sales_system_id"`
	SalesSystem                 *salessystem.SalesSystem `json:"sales_system"`
	InvoiceDate                 *string                  `json:"invoice_date" gorm:"type:DATE"`
	InvoiceNumber               string                   `json:"invoice_number"`
	InvoicePriceUnit            float64                  `json:"invoice_price_unit"`
	InvoicePriceTotal           float64                  `json:"invoice_price_total"`
	ContractDate                *string                  `json:"contract_date" gorm:"type:DATE"`
	ContractNumber              string                   `json:"contract_number"`
	DmoBuyerId                  *uint                    `json:"dmo_buyer_id"`
	DmoBuyer                    *company.Company         `json:"dmo_buyer"`
	DestinationCountryId        *uint                    `json:"destination_country_id"`
	DestinationCountry          *country.Country         `json:"destination_country"`
	DestinationId               *uint                    `json:"destination_id"`
	Destination                 *destination.Destination `json:"destination"`
	DmoDestinationPortLnName    string                   `json:"dmo_destination_port_ln_name"`
}

type ProductionJetty struct {
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
}

type SalesJetty struct {
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
}

type LossJetty struct {
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
}
type JettyBalanceLoss struct {
	ID           uint        `json:"id"`
	JettyId      uint        `json:"jetty_id"`
	Jetty        jetty.Jetty `json:"jetty"`
	StartBalance float64     `json:"start_balance"`
	TotalLoss    float64     `json:"total_loss"`
}
