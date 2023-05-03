package cafassignmentenduser

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/master/company"
)

type DetailCafAssignment struct {
	Detail                 cafassignment.CafAssignment `json:"detail"`
	ListEndUser            []CafAssignmentEndUser      `json:"list_end_user"`
	ListRealizationEndUser []RealizationEndUser        `json:"list_realization_end_user"`
}

type RealizationEndUser struct {
	ID                         uint            `json:"ID"`
	AverageCalories            float64         `json:"average_calories"`
	RealizationAverageCalories float64         `json:"realization_average_calories"`
	Quantity                   float64         `json:"quantity"`
	RealizationQuantity        float64         `json:"realization_quantity"`
	EndUserId                  uint            `json:"end_user_id"`
	EndUser                    company.Company `json:"end_user"`
	EndUserString              string          `json:"end_user_string"`
}

type Realization struct {
	RealizationAverageCalories float64 `json:"realization_average_calories"`
	RealizationQuantity        float64 `json:"realization_quantity"`
}
