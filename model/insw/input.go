package insw

type InputCreateInsw struct {
	ListGroupingVesselId []int  `json:"list_grouping_vessel_id" validate:"required,min=1"`
	Month                string `json:"month" validate:"required,LongMonth=January_February_March_April_May_June_July_August_September_October_November_December"`
	Year                 int    `json:"year" validate:"required"`
}
