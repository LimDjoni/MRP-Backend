package jettybalance

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"

	"gorm.io/gorm"
)

type JettyBalance struct {
	gorm.Model
	JettyId      uint          `json:"jetty_id"`
	Jetty        jetty.Jetty   `json:"jetty" gorm:"constraint:OnDelete:CASCADE;"`
	Year         string        `json:"year"`
	StartBalance float64       `json:"start_balance"`
	TotalLoss    float64       `json:"total_loss"`
	IupopkId     uint          `json:"iupopk_id"`
	Iupopk       iupopk.Iupopk `json:"iupopk"`
}
