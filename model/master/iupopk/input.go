package iupopk

type InputIupopk struct {
	Name         string  `json:"name" validate:"required"`
	Address      string  `json:"address" validate:"required"`
	Province     string  `json:"province" validate:"required"`
	Email        *string `json:"email" validate:"omitempty,email"`
	PhoneNumber  *string `json:"phone_number"`
	FaxNumber    *string `json:"fax_number"`
	DirectorName string  `json:"director_name" validate:"required"`
	Position     string  `json:"position" validate:"required"`
	Code         string  `json:"code" validate:"required"`
}
