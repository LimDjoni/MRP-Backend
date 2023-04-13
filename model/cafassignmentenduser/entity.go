package cafassignmentenduser

import (
	"ajebackend/model/cafassignment"

	"gorm.io/gorm"
)

type CafAssignmentEndUser struct {
	gorm.Model
	CafAssignmentId uint                        `json:"caf_assignment_id"`
	CafAssignment   cafassignment.CafAssignment `json:"caf_assignment" gorm:"constraint:OnDelete:CASCADE;"`
	AverageCalories float64                     `json:"average_calories"`
	Quantity        float64                     `json:"quantity"`
	EndUser         string                      `json:"end_user"`
}
