package groupingvesselln

import "ajebackend/model/insw"

type InputGroupingVesselLn struct {
	ListTransactions          []uint  `json:"list_transactions" validate:"required,min=1"`
	VesselName                string  `json:"vessel_name" validate:"required"`
	Quantity                  float64 `json:"quantity"`
	Adjustment                float64 `json:"adjustment"`
	GrandTotalQuantity        float64 `json:"grand_total_quantity"`
	DescriptionOfDocumentType string  `json:"description_of_document_type"`
	CodeOfDocumentType        string  `json:"code_of_document_type"`
	AjuNumber                 string  `json:"aju_number"`
	PebRegisterNumber         string  `json:"peb_register_number"`
	PebRegisterDate           *string `json:"peb_register_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	DescriptionOfPabeanOffice string  `json:"description_of_pabean_office"`
	CodeOfPabeanOffice        string  `json:"code_of_pabean_office"`
	SeriesPebGoods            string  `json:"series_peb_goods"`
	DescriptionOfGoods        string  `json:"description_of_goods"`
	TarifPosHs                string  `json:"tarif_pos_hs"`
	PebQuantity               float64 `json:"peb_quantity"`
	PebUnit                   string  `json:"peb_unit"`
	ExportValue               float64 `json:"export_value"`
	Currency                  string  `json:"currency"`
	LoadingPort               string  `json:"loading_port"`
	SkaCooNumber              string  `json:"ska_coo_number"`
	SkaCooDate                *string `json:"ska_coo_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	DestinationCountry        string  `json:"destination_country"`
	CodeOfDestinationCountry  string  `json:"code_of_destination_country"`
	LsExportNumber            string  `json:"ls_export_number"`
	LsExportDate              *string `json:"ls_export_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	InsuranceCompanyName      string  `json:"insurance_company_name"`
	PolisNumber               string  `json:"polis_number"`
	NavyCompanyName           string  `json:"navy_company_name"`
	NavyShipName              string  `json:"navy_ship_name"`
	NavyImoNumber             string  `json:"navy_imo_number"`
	Deadweight                float64 `json:"deadweight"`
	CowDate                   *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber                 string  `json:"cow_number"`
	CoaDate                   *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber                 string  `json:"coa_number"`
	QualityTmAr               float64 `json:"quality_tm_ar"`
	QualityImAdb              float64 `json:"quality_im_adb"`
	QualityAshAr              float64 `json:"quality_ash_ar"`
	QualityAshAdb             float64 `json:"quality_ash_adb"`
	QualityVmAdb              float64 `json:"quality_vm_adb"`
	QualityFcAdb              float64 `json:"quality_fc_adb"`
	QualityTsAr               float64 `json:"quality_ts_ar"`
	QualityTsAdb              float64 `json:"quality_ts_adb"`
	QualityCaloriesAr         float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb        float64 `json:"quality_calories_adb"`
	NettQualityCaloriesAr     float64 `json:"nett_quality_calories_ar"`
	IsCoaFinish               bool    `json:"is_coa_finish"`
	BlDate                    *string `json:"bl_date" validate:"omitempty,DateValidation"`
	BlNumber                  string  `json:"bl_number"`
	DmoDestinationPort        string  `json:"dmo_destination_port"`
}

type InputEditGroupingVesselLn struct {
	VesselName                string  `json:"vessel_name" validate:"required"`
	Quantity                  float64 `json:"quantity"`
	Adjustment                float64 `json:"adjustment"`
	GrandTotalQuantity        float64 `json:"grand_total_quantity"`
	DescriptionOfDocumentType string  `json:"description_of_document_type"`
	CodeOfDocumentType        string  `json:"code_of_document_type"`
	AjuNumber                 string  `json:"aju_number"`
	PebRegisterNumber         string  `json:"peb_register_number"`
	PebRegisterDate           *string `json:"peb_register_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	DescriptionOfPabeanOffice string  `json:"description_of_pabean_office"`
	CodeOfPabeanOffice        string  `json:"code_of_pabean_office"`
	SeriesPebGoods            string  `json:"series_peb_goods"`
	DescriptionOfGoods        string  `json:"description_of_goods"`
	TarifPosHs                string  `json:"tarif_pos_hs"`
	PebQuantity               float64 `json:"peb_quantity"`
	PebUnit                   string  `json:"peb_unit"`
	ExportValue               float64 `json:"export_value"`
	Currency                  string  `json:"currency"`
	LoadingPort               string  `json:"loading_port"`
	SkaCooNumber              string  `json:"ska_coo_number"`
	SkaCooDate                *string `json:"ska_coo_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	DestinationCountry        string  `json:"destination_country"`
	CodeOfDestinationCountry  string  `json:"code_of_destination_country"`
	LsExportNumber            string  `json:"ls_export_number"`
	LsExportDate              *string `json:"ls_export_date" gorm:"type:DATE" validate:"omitempty,DateValidation"`
	InsuranceCompanyName      string  `json:"insurance_company_name"`
	PolisNumber               string  `json:"polis_number"`
	NavyCompanyName           string  `json:"navy_company_name"`
	NavyShipName              string  `json:"navy_ship_name"`
	NavyImoNumber             string  `json:"navy_imo_number"`
	Deadweight                float64 `json:"deadweight"`
	CowDate                   *string `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber                 string  `json:"cow_number"`
	CoaDate                   *string `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber                 string  `json:"coa_number"`
	QualityTmAr               float64 `json:"quality_tm_ar"`
	QualityImAdb              float64 `json:"quality_im_adb"`
	QualityAshAr              float64 `json:"quality_ash_ar"`
	QualityAshAdb             float64 `json:"quality_ash_adb"`
	QualityVmAdb              float64 `json:"quality_vm_adb"`
	QualityFcAdb              float64 `json:"quality_fc_adb"`
	QualityTsAr               float64 `json:"quality_ts_ar"`
	QualityTsAdb              float64 `json:"quality_ts_adb"`
	QualityCaloriesAr         float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb        float64 `json:"quality_calories_adb"`
	NettQualityCaloriesAr     float64 `json:"nett_quality_calories_ar"`
	IsCoaFinish               bool    `json:"is_coa_finish"`
	BlDate                    *string `json:"bl_date" validate:"omitempty,DateValidation"`
	BlNumber                  string  `json:"bl_number"`
	DmoDestinationPort        string  `json:"dmo_destination_port"`
}

type SortFilterGroupingVesselLn struct {
	Field      string
	Sort       string
	Quantity   float64
	VesselName string
}

type InputRequestCreateUploadDmoInsw struct {
	Authorization    string           `json:"authorization"`
	Insw             insw.Insw        `json:"insw"`
	GroupingVesselLn GroupingVesselLn `json:"grouping_vessel"`
}
