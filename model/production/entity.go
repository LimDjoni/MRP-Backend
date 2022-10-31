package production

import (
	"gorm.io/gorm"
)

type Production struct {
	gorm.Model
	ProductionDate string `json:"production_date" gorm:"type:DATE"`
	Quantity float64 `json:"quantity"`
}
