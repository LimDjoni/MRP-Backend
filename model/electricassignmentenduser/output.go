package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/ports"
)

type DetailElectricAssignment struct {
	Detail                 electricassignment.ElectricAssignment `json:"detail"`
	ListEndUser            []ElectricAssignmentEndUser           `json:"list_end_user"`
	ListRealizationEndUser []RealizationEndUser                  `json:"list_realization_end_user"`
}

type RealizationEndUser struct {
	ID                         uint       `json:"ID"`
	PortId                     uint       `json:"port_id"`
	Port                       ports.Port `json:"port"`
	Supplier                   string     `json:"supplier"`
	AverageCalories            float64    `json:"average_calories"`
	RealizationAverageCalories float64    `json:"realization_average_calories"`
	Quantity                   float64    `json:"quantity"`
	RealizationQuantity        float64    `json:"realization_quantity"`
	EndUser                    string     `json:"end_user"`
	LetterNumber               string     `json:"letter_number"`
}

type Realization struct {
	RealizationAverageCalories float64 `json:"realization_average_calories"`
	RealizationQuantity        float64 `json:"realization_quantity"`
}
