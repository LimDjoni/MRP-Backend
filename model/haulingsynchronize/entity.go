package haulingsynchronize

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type HaulingSynchronize struct {
	gorm.Model
	LastSynchronizeIsp   *string       `json:"last_synchronize_isp" gorm:"DATETIME"`
	LastSynchronizeJetty *string       `json:"last_synchronize_jetty" gorm:"DATETIME"`
	IupopkId             uint          `json:"iupopk_id"`
	Iupopk               iupopk.Iupopk `json:"iupopk"`
	FailedCount          int           `json:"failed_count"`
}
