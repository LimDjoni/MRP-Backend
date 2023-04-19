package electricassignment

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type ElectricAssignment struct {
	gorm.Model
	IdNumber                     string        `json:"id_number" gorm:"UNIQUE"`
	Year                         string        `json:"year"`
	AssignmentLetterLink         string        `json:"assignment_letter_link"`
	RevisionAssignmentLetterLink string        `json:"revision_assignment_letter_link"`
	GrandTotalQuantity           float64       `json:"grand_total_quantity"`
	LetterNumber                 string        `json:"letter_number"`
	RevisionLetterNumber         string        `json:"revision_letter_number"`
	IupopkId                     uint          `json:"iupopk_id"`
	Iupopk                       iupopk.Iupopk `json:"iupopk"`
}
