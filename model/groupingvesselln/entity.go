package groupingvesselln

import (
	"ajebackend/model/insw"

	"gorm.io/gorm"
)

type GroupingVesselLn struct {
	gorm.Model
	IdNumber                  *string    `json:"id_number" gorm:"UNIQUE"`
	BlDate                    *string    `json:"bl_date" gorm:"type:DATE"`
	BlNumber                  string     `json:"bl_number"`
	UnloadingPortLocation     string     `json:"unloading_port_location"`
	VesselName                string     `json:"vessel_name"`
	Quantity                  float64    `json:"quantity"`
	Adjustment                float64    `json:"adjustment"`
	GrandTotalQuantity        float64    `json:"grand_total_quantity"`
	DescriptionOfDocumentType string     `json:"description_of_document_type"`
	CodeOfDocumentType        string     `json:"code_of_document_type"`
	AjuNumber                 string     `json:"aju_number"`
	PebRegisterNumber         string     `json:"peb_register_number"`
	PebRegisterDate           *string    `json:"peb_register_date" gorm:"type:DATE"`
	DescriptionOfPabeanOffice string     `json:"description_of_pabean_office"`
	CodeOfPabeanOffice        string     `json:"code_of_pabean_office"`
	SeriesPebGoods            string     `json:"series_peb_goods"`
	DescriptionOfGoods        string     `json:"description_of_goods"`
	TarifPosHs                string     `json:"tarif_pos_hs"`
	PebQuantity               float64    `json:"peb_quantity"`
	PebUnit                   string     `json:"peb_unit"`
	ExportValue               float64    `json:"export_value"`
	Currency                  string     `json:"currency"`
	LoadingPort               string     `json:"loading_port"`
	SkaCooNumber              string     `json:"ska_coo_number"`
	SkaCooDate                *string    `json:"ska_coo_date" gorm:"type:DATE"`
	DestinationCountry        string     `json:"destination_country"`
	CodeOfDestinationCountry  string     `json:"code_of_destination_country"`
	LsExportNumber            string     `json:"ls_export_number"`
	LsExportDate              *string    `json:"ls_export_date" gorm:"type:DATE"`
	InsuranceCompanyName      string     `json:"insurance_company_name"`
	PolisNumber               string     `json:"polis_number"`
	NavyCompanyName           string     `json:"navy_company_name"`
	NavyShipName              string     `json:"navy_ship_name"`
	NavyImoNumber             string     `json:"navy_imo_number"`
	Deadweight                float64    `json:"deadweight"`
	CowDate                   *string    `json:"cow_date" gorm:"type:DATE"`
	CowNumber                 string     `json:"cow_number"`
	CoaDate                   *string    `json:"coa_date" gorm:"type:DATE"`
	CoaNumber                 string     `json:"coa_number"`
	QualityTmAr               float64    `json:"quality_tm_ar"`
	QualityImAdb              float64    `json:"quality_im_adb"`
	QualityAshAr              float64    `json:"quality_ash_ar"`
	QualityAshAdb             float64    `json:"quality_ash_adb"`
	QualityVmAdb              float64    `json:"quality_vm_adb"`
	QualityFcAdb              float64    `json:"quality_fc_adb"`
	QualityTsAr               float64    `json:"quality_ts_ar"`
	QualityTsAdb              float64    `json:"quality_ts_adb"`
	QualityCaloriesAr         float64    `json:"quality_calories_ar"`
	QualityCaloriesAdb        float64    `json:"quality_calories_adb"`
	NettQualityCaloriesAr     float64    `json:"nett_quality_calories_ar"`
	PebDocumentLink           string     `json:"peb_document_link"`
	InsuranceDocumentLink     string     `json:"insurance_document_link"`
	LsExportDocumentLink      string     `json:"ls_export_document_link"`
	NavyDocumentLink          string     `json:"navy_document_link"`
	SkaCooDocumentLink        string     `json:"ska_coo_document_link"`
	CoaCowDocumentLink        string     `json:"coa_cow_document_link"`
	BlMvDocumentLink          string     `json:"bl_mv_document_link"`
	IsCoaFinish               bool       `json:"is_coa_finish"`
	InswId                    *uint      `json:"insw_id"`
	Insw                      *insw.Insw `json:"insw" gorm:"constraint:OnDelete:SET NULL;"`
}
