package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/company"
	"ajebackend/model/master/ports"

	"gorm.io/gorm"
)

type ElectricAssignmentEndUser struct {
	gorm.Model
	ElectricAssignmentId uint                                  `form:"electric_assignment_id" json:"electric_assignment_id"`
	ElectricAssignment   electricassignment.ElectricAssignment `form:"electric_assigment" json:"electric_assigment" gorm:"constraint:OnDelete:CASCADE;"`
	PortId               uint                                  `form:"port_id" json:"port_id"`
	Port                 ports.Port                            `form:"port" json:"port"`
	SupplierId           *uint                                 `form:"supplier_id" json:"supplier_id"`
	Supplier             *company.Company                      `form:"supplier" json:"supplier"`
	AverageCalories      float64                               `form:"average_calories" json:"average_calories"`
	Quantity             float64                               `form:"quantity" json:"quantity"`
	EndUser              string                                `form:"end_user" json:"end_user"`
	LetterNumber         string                                `form:"letter_number" json:"letter_number"`
}
