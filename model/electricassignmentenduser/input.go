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
	EndUserId       uint    `form:"end_user_id" json:"end_user_id"`
}

type UpdateElectricAssignmentInput struct {
	GrandTotalQuantity     float64                     `form:"grand_total_quantity" json:"grand_total_quantity"`
	RevisionLetterNumber   string                      `form:"revision_letter_number" json:"revision_letter_number"`
	ListElectricAssignment []ElectricAssignmentEndUser `form:"list_electric_assignment" json:"list_electric_assignment"`
}
