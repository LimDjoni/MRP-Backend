package groupingvesseldn

import (
	"ajebackend/model/destination"

	"gorm.io/gorm"
)

type GroupingVesselDn struct {
	gorm.Model
	IdNumber              *string                 `json:"id_number" gorm:"UNIQUE"`
	BlDate                *string                 `json:"bl_date" gorm:"type:DATE"`
	BlNumber              string                  `json:"bl_number"`
	SalesSystem           string                  `json:"sales_system"`
	VesselName            string                  `json:"vessel_name"`
	Quantity              float64                 `json:"quantity"`
	Adjustment            float64                 `json:"adjustment"`
	GrandTotalQuantity    float64                 `json:"grand_total_quantity"`
	CowDate               *string                 `json:"cow_date" gorm:"type:DATE"`
	CowNumber             *string                 `json:"cow_number"`
	CoaDate               *string                 `json:"coa_date" gorm:"type:DATE"`
	CoaNumber             *string                 `json:"coa_number"`
	QualityTmAr           *float64                `json:"quality_tm_ar"`
	QualityImAdb          *float64                `json:"quality_im_adb"`
	QualityAshAr          *float64                `json:"quality_ash_ar"`
	QualityAshAdb         *float64                `json:"quality_ash_adb"`
	QualityVmAdb          *float64                `json:"quality_vm_adb"`
	QualityFcAdb          *float64                `json:"quality_fc_adb"`
	QualityTsAr           *float64                `json:"quality_ts_ar"`
	QualityTsAdb          *float64                `json:"quality_ts_adb"`
	QualityCaloriesAr     *float64                `json:"quality_calories_ar"`
	QualityCaloriesAdb    *float64                `json:"quality_calories_adb"`
	DestinationId         uint                    `json:"destination_id"`
	Destination           destination.Destination `json:"destination"`
	DestinationCountry    string                  `json:"destination_country"`
	UnloadingPortLocation string                  `json:"unloading_port_location"`
	CoaCowDocumentLink    *string                 `json:"coa_cow_document_link"`
	BlMvDocumentLink      *string                 `json:"bl_mv_document_link"`
	IsCoaFinish           bool                    `json:"is_coa_finish"`
}
