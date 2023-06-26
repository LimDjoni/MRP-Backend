package cafassignmentenduser

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/master/company"

	"gorm.io/gorm"
)

type CafAssignmentEndUser struct {
	gorm.Model
	CafAssignmentId uint                        `json:"caf_assignment_id"`
	CafAssignment   cafassignment.CafAssignment `json:"caf_assignment" gorm:"constraint:OnDelete:CASCADE;"`
	AverageCalories float64                     `json:"average_calories"`
	Quantity        float64                     `json:"quantity"`
	EndUserString   string                      `json:"end_user_string"`
	EndUserId       uint                        `json:"end_user_id"`
	EndUser         company.Company             `json:"end_user"`
	LetterNumber    string                      `json:"letter_number"`
}
