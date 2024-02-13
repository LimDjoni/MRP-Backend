package production

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"

	"gorm.io/gorm"
)

type Production struct {
	gorm.Model
	ProductionDate string        `json:"production_date" gorm:"type:DATE"`
	Quantity       float64       `json:"quantity"`
	IupopkId       uint          `json:"iupopk_id"`
	Iupopk         iupopk.Iupopk `json:"iupopk"`
	RitaseQuantity int           `json:"ritase_quantity"`
	PitId          *uint         `json:"pit_id"`
	Pit            *pit.Pit      `json:"pit"`
	PitCode        *string       `json:"pit_code"`
	IspId          *uint         `json:"isp_id"`
	Isp            *isp.Isp      `json:"isp"`
	IspCode        *string       `json:"isp_code"`
	JettyId        *uint         `json:"jetty_id"`
	Jetty          *jetty.Jetty  `json:"jetty"`
	JettyCode      *string       `json:"jetty_code"`
}
