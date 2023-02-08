package groupingvesseldn

import (
	"ajebackend/model/master/company"
	"ajebackend/model/master/country"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/ports"
	"ajebackend/model/master/vessel"
	"ajebackend/model/reportdmo"

	"gorm.io/gorm"
)

type GroupingVesselDn struct {
	gorm.Model
	IdNumber             *string                  `json:"id_number" gorm:"UNIQUE"`
	BlDate               *string                  `json:"bl_date" gorm:"type:DATE"`
	BlNumber             string                   `json:"bl_number"`
	SalesSystem          string                   `json:"sales_system"`
	VesselId             *uint                    `json:"vessel_id"`
	Vessel               *vessel.Vessel           `json:"vessel"`
	Quantity             float64                  `json:"quantity"`
	Adjustment           float64                  `json:"adjustment"`
	GrandTotalQuantity   float64                  `json:"grand_total_quantity"`
	CowDate              *string                  `json:"cow_date" gorm:"type:DATE"`
	CowNumber            *string                  `json:"cow_number"`
	CoaDate              *string                  `json:"coa_date" gorm:"type:DATE"`
	CoaNumber            *string                  `json:"coa_number"`
	SkabDate             *string                  `json:"skab_date" gorm:"type:DATE"`
	SkabNumber           *string                  `json:"skab_number"`
	QualityTmAr          *float64                 `json:"quality_tm_ar"`
	QualityImAdb         *float64                 `json:"quality_im_adb"`
	QualityAshAr         *float64                 `json:"quality_ash_ar"`
	QualityAshAdb        *float64                 `json:"quality_ash_adb"`
	QualityVmAdb         *float64                 `json:"quality_vm_adb"`
	QualityFcAdb         *float64                 `json:"quality_fc_adb"`
	QualityTsAr          *float64                 `json:"quality_ts_ar"`
	QualityTsAdb         *float64                 `json:"quality_ts_adb"`
	QualityCaloriesAr    *float64                 `json:"quality_calories_ar"`
	QualityCaloriesAdb   *float64                 `json:"quality_calories_adb"`
	DestinationId        *uint                    `json:"destination_id"`
	Destination          *destination.Destination `json:"destination"`
	DestinationCountryId *uint                    `json:"destination_country_id"`
	DestinationCountry   *country.Country         `json:"destination_country"`
	DmoDestinationPortId *uint                    `json:"dmo_destination_port_id"`
	DmoDestinationPort   *ports.Port              `json:"dmo_destination_port"`
	BuyerName            string                   `json:"buyer_name"`
	BuyerId              *uint                    `json:"buyer_id"`
	Buyer                *company.Company         `json:"buyer"`
	BlMvDocumentLink     *string                  `json:"bl_mv_document_link"`
	CoaCowDocumentLink   *string                  `json:"coa_cow_document_link"`
	SkabDocumentLink     *string                  `json:"skab_document_link"`
	IsCoaFinish          bool                     `json:"is_coa_finish"`
	ReportDmoId          *uint                    `json:"report_dmo_id"`
	ReportDmo            *reportdmo.ReportDmo     `json:"report_dmo" gorm:"constraint:OnDelete:SET NULL;"`
}
