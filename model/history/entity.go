package history

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/coareport"
	"ajebackend/model/coareportln"
	"ajebackend/model/dmo"
	"ajebackend/model/electricassignment"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/jettybalance"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/rkab"
	"ajebackend/model/royaltyrecon"
	"ajebackend/model/royaltyreport"
	"ajebackend/model/transaction"
	"ajebackend/model/user"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	DmoId                *uint                                  `json:"dmo_id"`
	Dmo                  *dmo.Dmo                               `json:"dmo" gorm:"constraint:OnDelete:CASCADE;"`
	TransactionId        *uint                                  `json:"transaction_id"`
	Transaction          *transaction.Transaction               `json:"transaction" gorm:"constraint:OnDelete:CASCADE;"`
	UserId               uint                                   `json:"user_id"`
	User                 user.User                              `json:"user"`
	MinerbaId            *uint                                  `json:"minerba_id"`
	Minerba              *minerba.Minerba                       `json:"minerba" gorm:"constraint:OnDelete:CASCADE;"`
	ProductionId         *uint                                  `json:"production_id"`
	Production           production.Production                  `json:"production" gorm:"constraint:OnDelete:CASCADE;"`
	Status               string                                 `json:"status"`
	GroupingVesselDnId   *uint                                  `json:"grouping_vessel_dn_id"`
	GroupingVesselDn     groupingvesseldn.GroupingVesselDn      `json:"grouping_vessel_dn" gorm:"constraint:OnDelete:CASCADE;"`
	GroupingVesselLnId   *uint                                  `json:"grouping_vessel_ln_id"`
	GroupingVesselLn     groupingvesselln.GroupingVesselLn      `json:"grouping_vessel_ln" gorm:"constraint:OnDelete:CASCADE;"`
	MinerbaLnId          *uint                                  `json:"minerba_ln_id"`
	MinerbaLn            *minerbaln.MinerbaLn                   `json:"minerba_ln" gorm:"constraint:OnDelete:CASCADE;"`
	InswId               *uint                                  `json:"insw_id"`
	Insw                 *insw.Insw                             `json:"insw" gorm:"constraint:OnDelete:CASCADE;"`
	ReportDmoId          *uint                                  `json:"report_dmo_id"`
	ReportDmo            *reportdmo.ReportDmo                   `json:"report_dmo" gorm:"constraint:OnDelete:CASCADE;"`
	BeforeData           datatypes.JSON                         `json:"before_data"`
	AfterData            datatypes.JSON                         `json:"after_data"`
	IupopkId             *uint                                  `json:"iupopk_id"`
	Iupopk               *iupopk.Iupopk                         `json:"iupopk" gorm:"constraint:OnDelete:CASCADE;"`
	CoaReportId          *uint                                  `json:"coa_report_id"`
	CoaReport            *coareport.CoaReport                   `json:"coa_report" gorm:"constraint:OnDelete:CASCADE;"`
	CoaReportLnId        *uint                                  `json:"coa_report_ln_id"`
	CoaReportLn          *coareportln.CoaReportLn               `json:"coa_report_ln" gorm:"constraint:OnDelete:CASCADE;"`
	RkabId               *uint                                  `json:"rkab_id"`
	Rkab                 *rkab.Rkab                             `json:"rkab" gorm:"constraint:OnDelete:CASCADE;"`
	ElectricAssignmentId *uint                                  `json:"electric_assignment_id"`
	ElectricAssignment   *electricassignment.ElectricAssignment `json:"electrict_assignment" gorm:"constraint:OnDelete:CASCADE;"`
	CafAssignmentId      *uint                                  `json:"caf_assignment_id"`
	CafAssignment        *cafassignment.CafAssignment           `json:"caf_assignment" gorm:"constraint:OnDelete:CASCADE;"`
	RoyaltyReconId       *uint                                  `json:"royalty_recon_id"`
	RoyaltyRecon         *royaltyrecon.RoyaltyRecon             `json:"royalty_recon" gorm:"constraint:OnDelete:CASCADE;"`
	RoyaltyReportId      *uint                                  `json:"royalty_report_id"`
	RoyaltyReport        *royaltyreport.RoyaltyReport           `json:"royalty_report" gorm:"constraint:OnDelete:CASCADE;"`
	JettyBalanceId       *uint                                  `json:"jetty_balance_id"`
	JettyBalance         *jettybalance.JettyBalance             `json:"jetty_balance" gorm:"constraint:OnDelete:CASCADE;"`
}
