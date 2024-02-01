package ici

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Ici struct {
	gorm.Model
	Date      string        `json:"date"`
	Level     string        `json:"level"`
	Avarage   float64       `json:"average"`
	UnitPrice float64       `json:"unit_price"`
	Currency  string        `json:"currency"`
	IupopkId  uint          `json:"iupopk_id"`
	Iupopk    iupopk.Iupopk `json:"iupopk"`
}
