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
	LetterNumber          string                 `form:"letter_number" json:"letter_number"`
	RevisionLetterNumber  string                 `form:"revision_letter_number" json:"revision_letter_number"`
	ListCafAssignment     []CafAssignmentEndUser `form:"list_caf_assignment" json:"list_caf_assignment"`
	LetterNumber2         string                 `json:"letter_number2"`
	RevisionLetterNumber2 string                 `json:"revision_letter_number2"`
	LetterNumber3         string                 `json:"letter_number3"`
	RevisionLetterNumber3 string                 `json:"revision_letter_number3"`
	LetterNumber4         string                 `json:"letter_number4"`
	RevisionLetterNumber4 string                 `json:"revision_letter_number4"`
}
