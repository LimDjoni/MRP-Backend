package production

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Production struct {
	gorm.Model
	ProductionDate string        `json:"production_date" gorm:"type:DATE"`
	Quantity       float64       `json:"quantity"`
	IupopkId       uint          `json:"iupopk_id"`
	Iupopk         iupopk.Iupopk `json:"iupopk"`
}
