package insw

type InputCreateInsw struct {
	Month string `json:"month" validate:"required,ShortMonth=Jan_Feb_Mar_Apr_May_Jun_Jul_Aug_Sep_Oct_Nov_Dec"`
	Year  int    `json:"year" validate:"required"`
}

type InputUpdateDocumentInsw struct {
	Data []map[string]interface{} `json:"data"`
}

type SortFilterInsw struct {
	SortMonth string
	SortYear  string
	Month     string
	Year      string
}
