package cafassignment

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type CafAssignment struct {
	gorm.Model
	IdNumber                      string        `json:"id_number" gorm:"UNIQUE"`
	Year                          string        `json:"year"`
	AssignmentLetterLink          string        `json:"assignment_letter_link"`
	RevisionAssignmentLetterLink  string        `json:"revision_assignment_letter_link"`
	LetterNumber                  string        `json:"letter_number"`
	RevisionLetterNumber          string        `json:"revision_letter_number"`
	GrandTotalQuantity            float64       `json:"grand_total_quantity"`
	AssignmentLetterLink2         string        `json:"assignment_letter_link2"`
	RevisionAssignmentLetterLink2 string        `json:"revision_assignment_letter_link2"`
	LetterNumber2                 string        `json:"letter_number2"`
	RevisionLetterNumber2         string        `json:"revision_letter_number2"`
	AssignmentLetterLink3         string        `json:"assignment_letter_link3"`
	RevisionAssignmentLetterLink3 string        `json:"revision_assignment_letter_link3"`
	LetterNumber3                 string        `json:"letter_number3"`
	RevisionLetterNumber3         string        `json:"revision_letter_number3"`
	AssignmentLetterLink4         string        `json:"assignment_letter_link4"`
	RevisionAssignmentLetterLink4 string        `json:"revision_assignment_letter_link4"`
	LetterNumber4                 string        `json:"letter_number4"`
	RevisionLetterNumber4         string        `json:"revision_letter_number4"`
	IupopkId                      uint          `json:"iupopk_id"`
	Iupopk                        iupopk.Iupopk `json:"iupopk"`
}
