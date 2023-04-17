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
}

type UpdateCafAssignmentInput struct {
	GrandTotalQuantity   float64                `form:"grand_total_quantity" json:"grand_total_quantity"`
	RevisionLetterNumber string                 `form:"revision_letter_number" json:"revision_letter_number"`
	ListCafAssignment    []CafAssignmentEndUser `form:"list_caf_assignment" json:"list_caf_assignment"`
}
