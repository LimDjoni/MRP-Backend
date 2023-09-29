package electricassignmentenduser

type CreateElectricAssignmentInput struct {
	Year                   string                    `form:"year" json:"year"`
	GrandTotalQuantity     float64                   `form:"grand_total_quantity" json:"grand_total_quantity"`
	LetterNumber           string                    `form:"letter_number" json:"letter_number"`
	LetterDate             string                    `form:"letter_date" json:"letter_date"`
	ListElectricAssignment []ElectricAssignmentInput `form:"list_electric_assignment" json:"list_electric_assignment"`
}

type ElectricAssignmentInput struct {
	PortId          uint    `form:"port_id" json:"port_id"`
	AverageCalories float64 `form:"average_calories" json:"average_calories"`
	Quantity        float64 `form:"quantity" json:"quantity"`
	EndUser         string  `form:"end_user" json:"end_user"`
	LetterNumber    string  `form:"letter_number" json:"letter_number"`
	//SupplierId      *uint   `form:"supplier_id" json:"supplier_id"`
}

type UpdateElectricAssignmentInput struct {
	GrandTotalQuantity     float64                     `form:"grand_total_quantity" json:"grand_total_quantity"`
	GrandTotalQuantity2    float64                     `form:"grand_total_quantity2" json:"grand_total_quantity2"`
	GrandTotalQuantity3    float64                     `form:"grand_total_quantity3" json:"grand_total_quantity3"`
	GrandTotalQuantity4    float64                     `form:"grand_total_quantity4" json:"grand_total_quantity4"`
	LetterNumber           string                      `form:"letter_number" json:"letter_number"`
	LetterNumber2          string                      `form:"letter_number2" json:"letter_number2"`
	LetterNumber3          string                      `form:"letter_number3" json:"letter_number3"`
	LetterNumber4          string                      `form:"letter_number4" json:"letter_number4"`
	RevisionLetterNumber   string                      `form:"revision_letter_number" json:"revision_letter_number"`
	RevisionLetterNumber2  string                      `form:"revision_letter_number2" json:"revision_letter_number2"`
	RevisionLetterNumber3  string                      `form:"revision_letter_number3" json:"revision_letter_number3"`
	RevisionLetterNumber4  string                      `form:"revision_letter_number4" json:"revision_letter_number4"`
	LetterDate             string                      `form:"letter_date" json:"letter_date"`
	LetterDate2            *string                     `form:"letter_date2" json:"letter_date2"`
	LetterDate3            *string                     `form:"letter_date3" json:"letter_date3"`
	LetterDate4            *string                     `form:"letter_date4" json:"letter_date4"`
	RevisionLetterDate     *string                     `form:"revision_letter_date" json:"revision_letter_date"`
	RevisionLetterDate2    *string                     `form:"revision_letter_date2" json:"revision_letter_date2"`
	RevisionLetterDate3    *string                     `form:"revision_letter_date3" json:"revision_letter_date3"`
	RevisionLetterDate4    *string                     `form:"revision_letter_date4" json:"revision_letter_date4"`
	ListElectricAssignment []ElectricAssignmentEndUser `form:"list_electric_assignment" json:"list_electric_assignment"`
}
