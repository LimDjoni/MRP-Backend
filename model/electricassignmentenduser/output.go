package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/company"
	"ajebackend/model/master/ports"
)

type DetailElectricAssignment struct {
	Detail                  electricassignment.ElectricAssignment `json:"detail"`
	ListRealization         []ListRealization                     `json:"list_realization"`
	ListRealizationSupplier []RealizationSupplier                 `json:"list_realization_supplier"`
}

type RealizationSupplier struct {
	PortId                     uint             `json:"port_id"`
	Port                       ports.Port       `json:"port"`
	RealizationAverageCalories float64          `json:"realization_average_calories"`
	RealizationQuantity        float64          `json:"realization_quantity"`
	SupplierId                 *uint            `json:"supplier_id"`
	Supplier                   *company.Company `json:"supplier"`
}

type RealizationEndUser struct {
	ID                         uint       `json:"ID"`
	PortId                     uint       `json:"port_id"`
	Port                       ports.Port `json:"port"`
	AverageCalories            float64    `json:"average_calories"`
	RealizationAverageCalories float64    `json:"realization_average_calories"`
	Quantity                   float64    `json:"quantity"`
	RealizationQuantity        float64    `json:"realization_quantity"`
	EndUser                    string     `json:"end_user"`
	LetterNumber               string     `json:"letter_number"`
	// SupplierId                 *uint            `json:"supplier_id"`
	// Supplier                   *company.Company `json:"supplier"`
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
