package groupingvesselln

import (
	"ajebackend/model/insw"
	"ajebackend/model/master/country"
	"ajebackend/model/master/currency"
	"ajebackend/model/master/documenttype"
	"ajebackend/model/master/insurancecompany"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/navycompany"
	"ajebackend/model/master/navyship"
	"ajebackend/model/master/pabeanoffice"
	"ajebackend/model/master/portinsw"
	"ajebackend/model/master/unit"
	"ajebackend/model/master/vessel"

	"gorm.io/gorm"
)

type GroupingVesselLn struct {
	gorm.Model
	IdNumber              *string                            `json:"id_number" gorm:"UNIQUE"`
	BlDate                *string                            `json:"bl_date" gorm:"type:DATE"`
	BlNumber              string                             `json:"bl_number"`
	VesselId              *uint                              `json:"vessel_id"`
	Vessel                *vessel.Vessel                     `json:"vessel"`
	Quantity              float64                            `json:"quantity"`
	Adjustment            float64                            `json:"adjustment"`
	GrandTotalQuantity    float64                            `json:"grand_total_quantity"`
	DocumentTypeId        *uint                              `json:"document_type_id"`
	DocumentType          *documenttype.DocumentType         `json:"document_type"`
	AjuNumber             string                             `json:"aju_number"`
	PebRegisterNumber     string                             `json:"peb_register_number"`
	PebRegisterDate       *string                            `json:"peb_register_date" gorm:"type:DATE"`
	PabeanOfficeId        *uint                              `json:"pabean_office_id"`
	PabeanOffice          *pabeanoffice.PabeanOffice         `json:"pabean_office"`
	SeriesPebGoods        string                             `json:"series_peb_goods"`
	DescriptionOfGoods    string                             `json:"description_of_goods"`
	TarifPosHs            string                             `json:"tarif_pos_hs"`
	PebQuantity           float64                            `json:"peb_quantity"`
	PebUnitId             *uint                              `json:"peb_unit_id"`
	PebUnit               *unit.Unit                         `json:"peb_unit"`
	ExportValue           float64                            `json:"export_value"`
	CurrencyId            *uint                              `json:"currency_id"`
	Currency              *currency.Currency                 `json:"currency"`
	LoadingPortId         *uint                              `json:"loading_port_id"`
	LoadingPort           *portinsw.PortInsw                 `json:"loading_port"`
	SkaCooNumber          string                             `json:"ska_coo_number"`
	SkaCooDate            *string                            `json:"ska_coo_date" gorm:"type:DATE"`
	DestinationCountryId  *uint                              `json:"destination_country_id"`
	DestinationCountry    *country.Country                   `json:"destination_country"`
	LsExportNumber        string                             `json:"ls_export_number"`
	LsExportDate          *string                            `json:"ls_export_date" gorm:"type:DATE"`
	InsuranceCompanyId    *uint                              `json:"insurance_company_id"`
	InsuranceCompany      *insurancecompany.InsuranceCompany `json:"insurance_company"`
	PolisNumber           string                             `json:"polis_number"`
	NavyCompanyId         *uint                              `json:"navy_company_id"`
	NavyCompany           *navycompany.NavyCompany           `json:"navy_company"`
	NavyShipId            *uint                              `json:"navy_ship_id"`
	NavyShip              *navyship.NavyShip                 `json:"navy_ship"`
	NavyImoNumber         string                             `json:"navy_imo_number"`
	Deadweight            float64                            `json:"deadweight"`
	CowDate               *string                            `json:"cow_date" gorm:"type:DATE"`
	CowNumber             string                             `json:"cow_number"`
	CoaDate               *string                            `json:"coa_date" gorm:"type:DATE"`
	CoaNumber             string                             `json:"coa_number"`
	QualityTmAr           float64                            `json:"quality_tm_ar"`
	QualityImAdb          float64                            `json:"quality_im_adb"`
	QualityAshAr          float64                            `json:"quality_ash_ar"`
	QualityAshAdb         float64                            `json:"quality_ash_adb"`
	QualityVmAdb          float64                            `json:"quality_vm_adb"`
	QualityFcAdb          float64                            `json:"quality_fc_adb"`
	QualityTsAr           float64                            `json:"quality_ts_ar"`
	QualityTsAdb          float64                            `json:"quality_ts_adb"`
	QualityCaloriesAr     float64                            `json:"quality_calories_ar"`
	QualityCaloriesAdb    float64                            `json:"quality_calories_adb"`
	NettQualityCaloriesAr float64                            `json:"nett_quality_calories_ar"`
	PebDocumentLink       string                             `json:"peb_document_link"`
	InsuranceDocumentLink string                             `json:"insurance_document_link"`
	LsExportDocumentLink  string                             `json:"ls_export_document_link"`
	NavyDocumentLink      string                             `json:"navy_document_link"`
	SkaCooDocumentLink    string                             `json:"ska_coo_document_link"`
	CoaCowDocumentLink    string                             `json:"coa_cow_document_link"`
	BlMvDocumentLink      string                             `json:"bl_mv_document_link"`
	IsCoaFinish           bool                               `json:"is_coa_finish"`
	InswId                *uint                              `json:"insw_id"`
	Insw                  *insw.Insw                         `json:"insw" gorm:"constraint:OnDelete:SET NULL;"`
	IupopkId              uint                               `json:"iupopk_id"`
	Iupopk                iupopk.Iupopk                      `json:"iupopk"`
}
