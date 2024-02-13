package haulingsynchronize

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type HaulingSynchronize struct {
	gorm.Model
	LastSynchronizeIsp         *string       `json:"last_synchronize_isp" gorm:"DATETIME"`
	LastSynchronizeJetty       *string       `json:"last_synchronize_jetty" gorm:"DATETIME"`
	LastSynchronizeMasterToIsp *string       `json:"last_synchronize_master_to_isp" gorm:"DATETIME"`
	IupopkId                   uint          `json:"iupopk_id"`
	Iupopk                     iupopk.Iupopk `json:"iupopk"`
}
