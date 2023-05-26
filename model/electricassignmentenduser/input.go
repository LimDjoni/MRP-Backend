package electricassignmentenduser

type CreateElectricAssignmentInput struct {
	Year                   string                    `form:"year" json:"year"`
	GrandTotalQuantity     float64                   `form:"grand_total_quantity" json:"grand_total_quantity"`
	LetterNumber           string                    `form:"letter_number" json:"letter_number"`
	ListElectricAssignment []ElectricAssignmentInput `form:"list_electric_assignment" json:"list_electric_assignment"`
}

type ElectricAssignmentInput struct {
	PortId          uint    `form:"port_id" json:"port_id"`
	Supplier        string  `form:"supplier" json:"supplier"`
	AverageCalories float64 `form:"average_calories" json:"average_calories"`
	Quantity        float64 `form:"quantity" json:"quantity"`
	EndUser         string  `form:"end_user" json:"end_user"`
	LetterNumber    string  `form:"letter_number" json:"letter_number"`
}

type UpdateElectricAssignmentInput struct {
	GrandTotalQuantity     float64                     `form:"grand_total_quantity" json:"grand_total_quantity"`
	LetterNumber           string                      `form:"letter_number" json:"letter_number"`
	RevisionLetterNumber   string                      `form:"revision_letter_number" json:"revision_letter_number"`
	ListElectricAssignment []ElectricAssignmentEndUser `form:"list_electric_assignment" json:"list_electric_assignment"`
	LetterNumber2          string                      `json:"letter_number2"`
	RevisionLetterNumber2  string                      `json:"revision_letter_number2"`
	LetterNumber3          string                      `json:"letter_number3"`
	RevisionLetterNumber3  string                      `json:"revision_letter_number3"`
	LetterNumber4          string                      `json:"letter_number4"`
	RevisionLetterNumber4  string                      `json:"revision_letter_number4"`
}
