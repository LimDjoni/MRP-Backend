package insw

type InputCreateInsw struct {
	Month string `json:"month" validate:"required,LongMonth=January_February_March_April_May_June_July_August_September_October_November_December"`
	Year  int    `json:"year" validate:"required"`
}

type InputUpdateDocumentInsw struct {
	Data []map[string]interface{} `json:"data"`
}
