package groupingvesselln

import "ajebackend/model/insw"

type InputGroupingVesselLn struct {
	ListTransactions      []uint  `json:"list_transactions" validate:"required,min=1"`
	VesselId              *uint   `json:"vessel_id" validate:"required"`
	Quantity              float64 `json:"quantity"`
	Adjustment            float64 `json:"adjustment"`
	GrandTotalQuantity    float64 `json:"grand_total_quantity"`
	DocumentTypeId        *uint   `json:"document_type_id"`
	AjuNumber             string  `json:"aju_number"`
	PebRegisterNumber     string  `json:"peb_register_number"`
	PebRegisterDate       *string `json:"peb_register_date" validate:"omitempty,DateValidation"`
	PabeanOfficeId        *uint   `json:"pabean_office_id"`
	SeriesPebGoods        string  `json:"series_peb_goods"`
	DescriptionOfGoods    string  `json:"description_of_goods"`
	TarifPosHs            string  `json:"tarif_pos_hs"`
	PebQuantity           float64 `json:"peb_quantity"`
	PebUnitId             *uint   `json:"peb_unit_id"`
	ExportValue           float64 `json:"export_value"`
	CurrencyId            *uint   `json:"currency_id"`
	LoadingPortId         *uint   `json:"loading_port_id"`
	SkaCooNumber          string  `json:"ska_coo_number"`
	SkaCooDate            *string `json:"ska_coo_date" validate:"omitempty,DateValidation"`
	DestinationCountryId  *uint   `json:"destination_country_id"`
	LsExportNumber        string  `json:"ls_export_number"`
	LsExportDate          *string `json:"ls_export_date" validate:"omitempty,DateValidation"`
	InsuranceCompanyId    *uint   `json:"insurance_company_id"`
	PolisNumber           string  `json:"polis_number"`
	NavyCompanyName       string  `json:"navy_company_name"`
	NavyShipName          string  `json:"navy_ship_name"`
	NavyImoNumber         string  `json:"navy_imo_number"`
	Deadweight            float64 `json:"deadweight"`
	CowDate               *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber             string  `json:"cow_number"`
	CoaDate               *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber             string  `json:"coa_number"`
	QualityTmAr           float64 `json:"quality_tm_ar"`
	QualityImAdb          float64 `json:"quality_im_adb"`
	QualityAshAr          float64 `json:"quality_ash_ar"`
	QualityAshAdb         float64 `json:"quality_ash_adb"`
	QualityVmAdb          float64 `json:"quality_vm_adb"`
	QualityFcAdb          float64 `json:"quality_fc_adb"`
	QualityTsAr           float64 `json:"quality_ts_ar"`
	QualityTsAdb          float64 `json:"quality_ts_adb"`
	QualityCaloriesAr     float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb    float64 `json:"quality_calories_adb"`
	NettQualityCaloriesAr float64 `json:"nett_quality_calories_ar"`
	IsCoaFinish           bool    `json:"is_coa_finish"`
	BlDate                *string `json:"bl_date" validate:"DateValidation"`
	BlNumber              string  `json:"bl_number"`
}

type InputEditGroupingVesselLn struct {
	VesselId              *uint   `json:"vessel_id" validate:"required"`
	Quantity              float64 `json:"quantity"`
	Adjustment            float64 `json:"adjustment"`
	GrandTotalQuantity    float64 `json:"grand_total_quantity"`
	DocumentTypeId        *uint   `json:"document_type_id"`
	AjuNumber             string  `json:"aju_number"`
	PebRegisterNumber     string  `json:"peb_register_number"`
	PebRegisterDate       *string `json:"peb_register_date" validate:"omitempty,DateValidation"`
	PabeanOfficeId        *uint   `json:"pabean_office_id"`
	SeriesPebGoods        string  `json:"series_peb_goods"`
	DescriptionOfGoods    string  `json:"description_of_goods"`
	TarifPosHs            string  `json:"tarif_pos_hs"`
	PebQuantity           float64 `json:"peb_quantity"`
	PebUnitId             *uint   `json:"peb_unit_id"`
	ExportValue           float64 `json:"export_value"`
	CurrencyId            *uint   `json:"currency_id"`
	LoadingPortId         *uint   `json:"loading_port_id"`
	SkaCooNumber          string  `json:"ska_coo_number"`
	SkaCooDate            *string `json:"ska_coo_date" validate:"omitempty,DateValidation"`
	DestinationCountryId  *uint   `json:"destination_country_id"`
	LsExportNumber        string  `json:"ls_export_number"`
	LsExportDate          *string `json:"ls_export_date" validate:"omitempty,DateValidation"`
	InsuranceCompanyId    *uint   `json:"insurance_company_id"`
	PolisNumber           string  `json:"polis_number"`
	NavyCompanyName       string  `json:"navy_company_name"`
	NavyShipName          string  `json:"navy_ship_name"`
	NavyImoNumber         string  `json:"navy_imo_number"`
	Deadweight            float64 `json:"deadweight"`
	CowDate               *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber             string  `json:"cow_number"`
	CoaDate               *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber             string  `json:"coa_number"`
	QualityTmAr           float64 `json:"quality_tm_ar"`
	QualityImAdb          float64 `json:"quality_im_adb"`
	QualityAshAr          float64 `json:"quality_ash_ar"`
	QualityAshAdb         float64 `json:"quality_ash_adb"`
	QualityVmAdb          float64 `json:"quality_vm_adb"`
	QualityFcAdb          float64 `json:"quality_fc_adb"`
	QualityTsAr           float64 `json:"quality_ts_ar"`
	QualityTsAdb          float64 `json:"quality_ts_adb"`
	QualityCaloriesAr     float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb    float64 `json:"quality_calories_adb"`
	NettQualityCaloriesAr float64 `json:"nett_quality_calories_ar"`
	IsCoaFinish           bool    `json:"is_coa_finish"`
	BlDate                *string `json:"bl_date" validate:"DateValidation"`
	BlNumber              string  `json:"bl_number"`
}

type SortFilterGroupingVesselLn struct {
	Field      string
	Sort       string
	Quantity   float64
	VesselName string
}

type InputRequestCreateUploadInsw struct {
	Authorization    string             `json:"authorization"`
	Insw             insw.Insw          `json:"insw"`
	GroupingVesselLn []GroupingVesselLn `json:"grouping_vessel"`
}
