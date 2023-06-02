package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/company"
	"ajebackend/model/master/ports"
)

type DetailElectricAssignment struct {
	Detail          electricassignment.ElectricAssignment `json:"detail"`
	ListRealization []ListRealization                     `json:"list_realization"`
}

type RealizationEndUser struct {
	ID                         uint             `json:"ID"`
	PortId                     uint             `json:"port_id"`
	Port                       ports.Port       `json:"port"`
	SupplierId                 *uint            `json:"supplier_id"`
	Supplier                   *company.Company `json:"supplier"`
	AverageCalories            float64          `json:"average_calories"`
	RealizationAverageCalories float64          `json:"realization_average_calories"`
	Quantity                   float64          `json:"quantity"`
	RealizationQuantity        float64          `json:"realization_quantity"`
	EndUser                    string           `json:"end_user"`
	LetterNumber               string           `json:"letter_number"`
}

type ListRealization struct {
	Order                  int                  `json:"order"`
	LetterNumber           string               `json:"letter_number"`
	ListRealizationEndUser []RealizationEndUser `json:"list_realization_end_user"`
}

type Realization struct {
	RealizationAverageCalories float64 `json:"realization_average_calories"`
	RealizationQuantity        float64 `json:"realization_quantity"`
}
