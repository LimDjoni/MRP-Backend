package groupingvesseldn

type InputGroupingVesselDn struct {
	ListTransactions     []uint   `json:"list_transactions" validate:"required,min=1"`
	BlNumber             string   `json:"bl_number"`
	VesselId             *uint    `json:"vessel_id" validate:"required"`
	Quantity             float64  `json:"quantity"`
	Adjustment           float64  `json:"adjustment"`
	GrandTotalQuantity   float64  `json:"grand_total_quantity"`
	BlDate               *string  `json:"bl_date" validate:"DateValidation"`
	SalesSystem          string   `json:"sales_system" validate:"required,SalesSystemValidation=Barge_Vessel"`
	DestinationId        *uint    `json:"destination_id"`
	DestinationCountryId *uint    `json:"destination_country_id"`
	DmoDestinationPortId *uint    `json:"dmo_destination_port_id"`
	BuyerId              *uint    `json:"buyer_id"`
	CowDate              *string  `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber            string   `json:"cow_number"`
	CoaDate              *string  `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber            string   `json:"coa_number"`
	SkabDate             *string  `json:"skab_date" gorm:"type:DATE"`
	SkabNumber           string   `json:"skab_number"`
	QualityTmAr          *float64 `json:"quality_tm_ar"`
	QualityImAdb         *float64 `json:"quality_im_adb"`
	QualityAshAr         *float64 `json:"quality_ash_ar"`
	QualityAshAdb        *float64 `json:"quality_ash_adb"`
	QualityVmAdb         *float64 `json:"quality_vm_adb"`
	QualityFcAdb         *float64 `json:"quality_fc_adb"`
	QualityTsAr          *float64 `json:"quality_ts_ar"`
	QualityTsAdb         *float64 `json:"quality_ts_adb"`
	QualityCaloriesAr    *float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb   *float64 `json:"quality_calories_adb"`
	IsCoaFinish          bool     `json:"is_coa_finish"`
}

type InputEditGroupingVesselDn struct {
	BlNumber             string   `json:"bl_number"`
	VesselId             *uint    `json:"vessel_id" validate:"required"`
	Quantity             float64  `json:"quantity"`
	Adjustment           float64  `json:"adjustment"`
	GrandTotalQuantity   float64  `json:"grand_total_quantity"`
	BlDate               *string  `json:"bl_date" validate:"DateValidation"`
	SalesSystem          string   `json:"sales_system" validate:"required,SalesSystemValidation=Barge_Vessel"`
	DestinationId        *uint    `json:"destination_id"`
	DestinationCountryId *uint    `json:"destination_country_id"`
	DmoDestinationPortId *uint    `json:"dmo_destination_port_id"`
	BuyerId              *uint    `json:"buyer_id"`
	CowDate              *string  `json:"cow_date" validate:"omitempty,DateValidation"`
	CowNumber            string   `json:"cow_number"`
	CoaDate              *string  `json:"coa_date" validate:"omitempty,DateValidation"`
	CoaNumber            string   `json:"coa_number"`
	SkabDate             *string  `json:"skab_date" gorm:"type:DATE"`
	SkabNumber           string   `json:"skab_number"`
	QualityTmAr          *float64 `json:"quality_tm_ar"`
	QualityImAdb         *float64 `json:"quality_im_adb"`
	QualityAshAr         *float64 `json:"quality_ash_ar"`
	QualityAshAdb        *float64 `json:"quality_ash_adb"`
	QualityVmAdb         *float64 `json:"quality_vm_adb"`
	QualityFcAdb         *float64 `json:"quality_fc_adb"`
	QualityTsAr          *float64 `json:"quality_ts_ar"`
	QualityTsAdb         *float64 `json:"quality_ts_adb"`
	QualityCaloriesAr    *float64 `json:"quality_calories_ar"`
	QualityCaloriesAdb   *float64 `json:"quality_calories_adb"`
	IsCoaFinish          bool     `json:"is_coa_finish"`
}

type SortFilterGroupingVesselDn struct {
	Field       string
	Sort        string
	Quantity    string
	VesselId    string
	BlDateStart string
	BlDateEnd   string
}
