package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/ports"

	"gorm.io/gorm"
)

type ElectricAssignmentEndUser struct {
	gorm.Model
	ElectricAssignmentId uint                                  `json:"electric_assignment_id"`
	ElectricAssignment   electricassignment.ElectricAssignment `json:"electric_assigment" gorm:"constraint:OnDelete:CASCADE;"`
	PortId               uint                                  `json:"port_id"`
	Port                 ports.Port                            `json:"port"`
	Supplier             string                                `json:"supplier"`
	AverageCalories      float64                               `json:"average_calories"`
	Quantity             float64                               `json:"quantity"`
	EndUser              string                                `json:"end_user"`
	LetterNumber         string                                `json:"letter_number"`
}
