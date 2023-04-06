package logs

import (
	"ajebackend/model/coareport"
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/rkab"
	"ajebackend/model/transaction"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Logs struct {
	gorm.Model
	Input              datatypes.JSON                    `json:"input"`
	Message            datatypes.JSON                    `json:"message"`
	TransactionId      *uint                             `json:"transaction_id"`
	Transaction        *transaction.Transaction          `json:"transaction" gorm:"constraint:OnDelete:CASCADE;"`
	MinerbaId          *uint                             `json:"minerba_id"`
	Minerba            *minerba.Minerba                  `json:"minerba" gorm:"constraint:OnDelete:CASCADE;"`
	DmoId              *uint                             `json:"dmo_id"`
	Dmo                *dmo.Dmo                          `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	ProductionId       *uint                             `json:"production_id"`
	Production         production.Production             `json:"production" gorm:"constraint:OnDelete:CASCADE;"`
	GroupingVesselDnId *uint                             `json:"grouping_vessel_dn_id"`
	GroupingVesselDn   groupingvesseldn.GroupingVesselDn `json:"grouping_vessel_dn" gorm:"constraint:OnDelete:CASCADE;"`
	GroupingVesselLnId *uint                             `json:"grouping_vessel_ln_id"`
	GroupingVesselLn   groupingvesselln.GroupingVesselLn `json:"grouping_vessel_ln" gorm:"constraint:OnDelete:CASCADE;"`
	MinerbaLnId        *uint                             `json:"minerba_ln_id"`
	MinerbaLn          *minerbaln.MinerbaLn              `json:"minerba_ln" gorm:"constraint:OnDelete:CASCADE;"`
	InswId             *uint                             `json:"insw_id"`
	Insw               *insw.Insw                        `json:"insw" gorm:"constraint:OnDelete:CASCADE;"`
	ReportDmoId        *uint                             `json:"report_dmo_id"`
	ReportDmo          *reportdmo.ReportDmo              `json:"report_dmo"`
	CoaReportId        *uint                             `json:"coa_report_id"`
	CoaReport          *coareport.CoaReport              `json:"coa_report" gorm:"constraint:OnDelete:CASCADE;"`
	RkabId             *uint                             `json:"rkab_id"`
	Rkab               *rkab.Rkab                        `json:"rkab" gorm:"constraint:OnDelete:CASCADE;"`
}
