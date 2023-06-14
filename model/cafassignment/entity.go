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
	GrandTotalQuantity2           float64       `json:"grand_total_quantity2"`
	GrandTotalQuantity3           float64       `json:"grand_total_quantity3"`
	GrandTotalQuantity4           float64       `json:"grand_total_quantity4"`
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
	LetterDate                    string        `json:"letter_date" gorm:"type:DATE"`
	RevisionLetterDate            *string       `json:"revision_letter_date" gorm:"type:DATE"`
	LetterDate2                   *string       `json:"letter_date2" gorm:"type:DATE"`
	RevisionLetterDate2           *string       `json:"revision_letter_date2" gorm:"type:DATE"`
	LetterDate3                   *string       `json:"letter_date3" gorm:"type:DATE"`
	RevisionLetterDate3           *string       `json:"revision_letter_date3" gorm:"type:DATE"`
	LetterDate4                   *string       `json:"letter_date4" gorm:"type:DATE"`
	RevisionLetterDate4           *string       `json:"revision_letter_date4" gorm:"type:DATE"`
	IupopkId                      uint          `json:"iupopk_id"`
	Iupopk                        iupopk.Iupopk `json:"iupopk"`
}
