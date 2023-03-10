package iupopk

import (
	"gorm.io/gorm"
)

type Iupopk struct {
	gorm.Model
	Name         string  `json:"name" gorm:"UNIQUE"`
	Address      string  `json:"address"`
	Province     string  `json:"province"`
	Email        *string `json:"email" gorm:"UNIQUE"`
	PhoneNumber  *string `json:"phone_number" gorm:"UNIQUE"`
	FaxNumber    *string `json:"fax_number" gorm:"UNIQUE"`
	DirectorName string  `json:"director_name"`
	Position     string  `json:"position"`
	Code         string  `json:"code"`
}
