package cafassignmentenduser

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/master/company"
)

type DetailCafAssignment struct {
	Detail          cafassignment.CafAssignment `json:"detail"`
	ListRealization []ListRealization           `json:"list_realization"`
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
	LetterNumber               string          `json:"letter_number"`
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
