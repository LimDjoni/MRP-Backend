package cafassignmentenduser

type CreateCafAssignmentInput struct {
	Year               string               `form:"year" json:"year"`
	GrandTotalQuantity float64              `form:"grand_total_quantity" json:"grand_total_quantity"`
	LetterNumber       string               `form:"letter_number" json:"letter_number"`
	ListCafAssignment  []CafAssignmentInput `form:"list_caf_assignment" json:"list_caf_assignment"`
}

type CafAssignmentInput struct {
	AverageCalories float64 `form:"average_calories" json:"average_calories"`
	Quantity        float64 `form:"quantity" json:"quantity"`
	EndUserId       uint    `form:"end_user_id" json:"end_user_id"`
	LetterNumber    string  `form:"letter_number" json:"letter_number"`
}

type UpdateCafAssignmentInput struct {
	GrandTotalQuantity    float64                `form:"grand_total_quantity" json:"grand_total_quantity"`
	GrandTotalQuantity2   float64                `form:"grand_total_quantity2" json:"grand_total_quantity2"`
	GrandTotalQuantity3   float64                `form:"grand_total_quantity3" json:"grand_total_quantity3"`
	GrandTotalQuantity4   float64                `form:"grand_total_quantity4" json:"grand_total_quantity4"`
	LetterNumber          string                 `form:"letter_number" json:"letter_number"`
	LetterNumber2         string                 `form:"letter_number2" json:"letter_number2"`
	LetterNumber3         string                 `form:"letter_number3" json:"letter_number3"`
	LetterNumber4         string                 `form:"letter_number4" json:"letter_number4"`
	RevisionLetterNumber  string                 `form:"revision_letter_number" json:"revision_letter_number"`
	RevisionLetterNumber2 string                 `form:"revision_letter_number2" json:"revision_letter_number2"`
	RevisionLetterNumber3 string                 `form:"revision_letter_number3" json:"revision_letter_number3"`
	RevisionLetterNumber4 string                 `form:"revision_letter_number4" json:"revision_letter_number4"`
	ListCafAssignment     []CafAssignmentEndUser `form:"list_caf_assignment" json:"list_caf_assignment"`
}
