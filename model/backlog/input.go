package backlog

type RegisterBackLogInput struct {
	UnitId            uint     `json:"unit_id"`
	HMBreakdown       float64  `json:"hm_breakdown"`
	Problem           string   `json:"problem"`
	Component         string   `json:"component"`
	PartNumber        string   `json:"part_number"`
	PartDescription   string   `json:"part_description"`
	QtyOrder          uint     `json:"qty_order"`
	DateOfInspection  string   `json:"date_of_inspection" gorm:"DATE"`
	PlanReplaceRepair *string  `json:"plan_replace_repair" gorm:"DATE"`
	HMReady           *float64 `json:"hm_ready"`
	PPNumber          *string  `json:"pp_number"`
	PONumber          *string  `json:"po_number"`
	Status            string   `json:"status"`
}

type SortFilterBackLog struct {
	Field              string
	Sort               string
	BrandName          string
	UnitId             string
	HMBreakdown        string
	Problem            string
	Component          string
	PartNumber         string
	PartDescription    string
	QtyOrder           string
	DateOfInspection   string
	PlanReplaceRepair  string
	AgingBacklogByDate string
	HMReady            string
	PPNumber           string
	PONumber           string
	Status             string
}

type SortFilterBackLogSummary struct {
	Field             string
	Sort              string
	BrandName         string
	UnitId            string
	HMBreakdown       string
	Problem           string
	Component         string
	PartNumber        string
	PartDescription   string
	QtyOrder          string
	DateOfInspection  string
	PlanReplaceRepair string
	HMReady           string
	PPNumber          string
	PONumber          string
	Status            string
}

type SortFilterDashboardBacklog struct {
	Year string
}

type DashboardBackLog struct {
	TotalBackLog uint           `json:"total_backlog"`
	Total1       uint           `json:"total_1"`
	Total2       uint           `json:"total_2"`
	Total3       uint           `json:"total_3"`
	Total4       uint           `json:"total_4"`
	Summary      BacklogSummary `json:"backlog_summary"`
	AgingSummary AgingSummary   `json:"aging_summary"`
}

type BacklogSummary struct {
	Pending   uint `json:"pending"`
	Open      uint `json:"open"`
	Closed    uint `json:"closed"`
	Cancelled uint `json:"cancelled"`
	Rejected  uint `json:"rejected"`
}

type rawAging struct {
	Pending0_5   uint
	Open0_5      uint
	Closed0_5    uint
	Cancelled0_5 uint
	Rejected0_5  uint

	Pending6_15   uint
	Open6_15      uint
	Closed6_15    uint
	Cancelled6_15 uint
	Rejected6_15  uint

	Pending16_30   uint
	Open16_30      uint
	Closed16_30    uint
	Cancelled16_30 uint
	Rejected16_30  uint

	Pending30plus   uint
	Open30plus      uint
	Closed30plus    uint
	Cancelled30plus uint
	Rejected30plus  uint
}

type AgingSummary struct {
	AgingTotal1 BacklogSummary `json:"aging_total_1"`
	AgingTotal2 BacklogSummary `json:"aging_total_2"`
	AgingTotal3 BacklogSummary `json:"aging_total_3"`
	AgingTotal4 BacklogSummary `json:"aging_total_4"`
}
