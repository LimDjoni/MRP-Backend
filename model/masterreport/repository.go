package masterreport

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/electricassignment"
	"ajebackend/model/electricassignmentenduser"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/master/categoryindustrytype"
	"ajebackend/model/master/jetty"
	"ajebackend/model/pitloss"
	"ajebackend/model/production"
	"ajebackend/model/rkab"
	"ajebackend/model/transaction"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func helperString(listString []string, dataString string) bool {
	for _, v := range listString {
		if v == dataString {
			return true
		}
	}
	return false
}

type Repository interface {
	RecapDmo(year string, iupopkId int) (ReportDmoOutput, error)
	RealizationReport(year string, iupopkId int) (RealizationOutput, error)
	SaleDetailReport(year string, iupopkId int) (SaleDetail, error)
	GetTransactionReport(iupopkId int, input TransactionReportInput, typeTransaction string) ([]TransactionReport, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) RecapDmo(year string, iupopkId int) (ReportDmoOutput, error) {

	var reportDmoOuput ReportDmoOutput

	var listProduction []production.Production

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	// Production Query
	queryFilterProduction := fmt.Sprintf("production_date >= '%s' AND production_date <= '%s' AND iupopk_id = %v", startFilter, endFilter, iupopkId)

	errFindProduction := r.db.Where(queryFilterProduction).Order("id ASC").Find(&listProduction).Error

	if errFindProduction != nil {
		return reportDmoOuput, errFindProduction
	}

	for _, v := range listProduction {
		date, _ := time.Parse("2006-01-02T00:00:00Z", v.ProductionDate)
		_, month, _ := date.Date()
		switch int(month) {
		case 1:
			reportDmoOuput.Production.January += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 2:
			reportDmoOuput.Production.February += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 3:
			reportDmoOuput.Production.March += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 4:
			reportDmoOuput.Production.April += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 5:
			reportDmoOuput.Production.May += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 6:
			reportDmoOuput.Production.June += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 7:
			reportDmoOuput.Production.July += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 8:
			reportDmoOuput.Production.August += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 9:
			reportDmoOuput.Production.September += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 10:
			reportDmoOuput.Production.October += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 11:
			reportDmoOuput.Production.November += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		case 12:
			reportDmoOuput.Production.December += v.Quantity
			reportDmoOuput.Production.Total += v.Quantity
		}
	}

	// Rkabs Query
	var rkabs []rkab.Rkab

	queryRkab := fmt.Sprintf("year = '%s' AND iupopk_id = %v", year, iupopkId)

	errFindRkab := r.db.Where(queryRkab).Order("id ASC").Find(&rkabs).Error

	if errFindRkab != nil {
		return reportDmoOuput, errFindRkab
	}

	reportDmoOuput.Rkabs = rkabs

	// Transaction Query
	var listTransactions []transaction.Transaction

	queryFilter := fmt.Sprintf("transactions.seller_id = %v AND transactions.transaction_type = 'DN' AND dmos.period LIKE '%%%v' AND transactions.dmo_id IS NOT NULL AND (grouping_vessel_dns.sales_system != 'Vessel' OR grouping_vessel_dn_id IS NULL)", iupopkId, year)

	errFind := r.db.Preload("ReportDmo").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Table("transactions").Select("transactions.*").Joins("left join dmos on dmos.id = transactions.dmo_id left join grouping_vessel_dns on grouping_vessel_dns.id = transactions.grouping_vessel_dn_id").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

	if errFind != nil {
		return reportDmoOuput, errFind
	}

	for _, v := range listTransactions {

		if v.ReportDmoId == nil {
			continue
		}

		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		if v.IsNotClaim == false {
			if v.GroupingVesselDnId != nil && v.SalesSystem != nil && strings.Contains(v.SalesSystem.Name, "Vessel") {
				continue
			}

			switch periodSplit[0] {
			case "Jan":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.January += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.January += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.January += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Feb":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.February += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.February += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.February += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Mar":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.March += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.March += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.March += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Apr":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.April += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.April += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.April += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "May":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.May += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.May += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.May += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Jun":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.June += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.June += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.June += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Jul":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.July += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.July += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.July += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Aug":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.August += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.August += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.August += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Sep":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.September += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.September += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.September += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Oct":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.October += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.October += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.October += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Nov":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.November += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.November += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.November += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			case "Dec":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.December += v.Quantity
					reportDmoOuput.RecapElectricity.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.December += v.Quantity
					reportDmoOuput.RecapCement.Total += v.Quantity
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.December += v.Quantity
					reportDmoOuput.RecapNonElectricity.Total += v.Quantity
				}
			}
		} else {
			switch periodSplit[0] {
			case "Jan":
				reportDmoOuput.NotClaimable.January += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Feb":
				reportDmoOuput.NotClaimable.February += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Mar":
				reportDmoOuput.NotClaimable.March += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Apr":
				reportDmoOuput.NotClaimable.April += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "May":
				reportDmoOuput.NotClaimable.May += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Jun":
				reportDmoOuput.NotClaimable.June += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Jul":
				reportDmoOuput.NotClaimable.July += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Aug":
				reportDmoOuput.NotClaimable.August += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Sep":
				reportDmoOuput.NotClaimable.September += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Oct":
				reportDmoOuput.NotClaimable.October += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Nov":
				reportDmoOuput.NotClaimable.November += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			case "Dec":
				reportDmoOuput.NotClaimable.December += v.Quantity
				reportDmoOuput.NotClaimable.Total += v.Quantity
			}
		}
	}

	var groupingVessels []groupingvesseldn.GroupingVesselDn

	queryFilterGrouping := fmt.Sprintf("grouping_vessel_dns.iupopk_id = %v AND report_dmos.period LIKE '%%%v' AND grouping_vessel_dns.sales_system = 'Vessel' AND report_dmo_id IS NOT NULL", iupopkId, year)

	errFindGrouping := r.db.Preload("ReportDmo").Preload("Buyer.IndustryType.CategoryIndustryType").Table("grouping_vessel_dns").Select("grouping_vessel_dns.*").Joins("left join report_dmos on report_dmos.id = grouping_vessel_dns.report_dmo_id").Where(queryFilterGrouping).Find(&groupingVessels).Error

	if errFindGrouping != nil {
		return reportDmoOuput, errFindGrouping
	}

	for _, v := range groupingVessels {
		if v.ReportDmoId == nil {
			continue
		}
		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		switch periodSplit[0] {
		case "Jan":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.January += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.January += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.January += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Feb":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.February += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.February += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.February += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Mar":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.March += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.March += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.March += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Apr":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.April += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.April += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.April += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "May":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.May += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.May += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.May += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Jun":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.June += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.June += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.June += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Jul":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.July += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.July += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.July += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Aug":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.August += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.August += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.August += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Sep":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.September += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.September += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.September += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Oct":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.October += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.October += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.October += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Nov":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.November += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.November += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.November += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		case "Dec":
			if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
				reportDmoOuput.RecapElectricity.December += v.GrandTotalQuantity
				reportDmoOuput.RecapElectricity.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
				reportDmoOuput.RecapCement.December += v.GrandTotalQuantity
				reportDmoOuput.RecapCement.Total += v.GrandTotalQuantity
			} else if v.Buyer != nil && v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
				reportDmoOuput.RecapNonElectricity.December += v.GrandTotalQuantity
				reportDmoOuput.RecapNonElectricity.Total += v.GrandTotalQuantity
			}
		}
	}

	return reportDmoOuput, nil
}

func (r *repository) RealizationReport(year string, iupopkId int) (RealizationOutput, error) {
	var realizationOutput RealizationOutput

	var listTransactions []transaction.Transaction

	queryFilter := fmt.Sprintf("transactions.seller_id = %v AND transactions.transaction_type = 'DN' AND dmos.period LIKE '%%%v' AND dmo_id IS NOT NULL", iupopkId, year)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Table("transactions").Select("transactions.*").Joins("left join dmos on dmos.id = transactions.dmo_id").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

	if errFind != nil {
		return realizationOutput, errFind
	}

	for _, v := range listTransactions {
		if v.ReportDmoId == nil || (v.GroupingVesselDnId != nil && v.SalesSystem != nil && strings.Contains(v.SalesSystem.Name, "Vessel")) {
			continue
		}

		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		var transactionTemp RealizationTransaction

		switch periodSplit[0] {
		case "Jan":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.January = append(realizationOutput.Electric.January, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.January = append(realizationOutput.Cement.January, transactionTemp)
					} else {
						realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
			}
		case "Feb":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.February = append(realizationOutput.Electric.February, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.February = append(realizationOutput.Cement.February, transactionTemp)
					} else {
						realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
			}
		case "Mar":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.March = append(realizationOutput.Electric.March, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.March = append(realizationOutput.Cement.March, transactionTemp)
					} else {
						realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
			}
		case "Apr":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.April = append(realizationOutput.Electric.April, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.April = append(realizationOutput.Cement.April, transactionTemp)
					} else {
						realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
			}
		case "May":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.May = append(realizationOutput.Electric.May, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.May = append(realizationOutput.Cement.May, transactionTemp)
					} else {
						realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
			}
		case "Jun":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.June = append(realizationOutput.Electric.June, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.June = append(realizationOutput.Cement.June, transactionTemp)
					} else {
						realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
			}
		case "Jul":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.July = append(realizationOutput.Electric.July, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.July = append(realizationOutput.Cement.July, transactionTemp)
					} else {
						realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
			}
		case "Aug":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.August = append(realizationOutput.Electric.August, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.August = append(realizationOutput.Cement.August, transactionTemp)
					} else {
						realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
			}
		case "Sep":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.September = append(realizationOutput.Electric.September, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.September = append(realizationOutput.Cement.September, transactionTemp)
					} else {
						realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
			}
		case "Oct":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.October = append(realizationOutput.Electric.October, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.October = append(realizationOutput.Cement.October, transactionTemp)
					} else {
						realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
			}
		case "Nov":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.November = append(realizationOutput.Electric.November, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.November = append(realizationOutput.Cement.November, transactionTemp)
					} else {
						realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
			}
		case "Dec":
			var shippingDate string
			shippingDateSplit := strings.Split(*v.ShippingDate, "T")
			shippingDate = shippingDateSplit[0]
			transactionTemp.ShippingDate = shippingDate
			if v.Customer != nil {
				transactionTemp.Trader = v.Customer
			}

			if v.DmoBuyer != nil {
				transactionTemp.EndUser = v.DmoBuyer
			}
			transactionTemp.Quantity = v.Quantity
			transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
			if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.DmoBuyer != nil {
				if v.DmoBuyer.IndustryType != nil {
					if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.December = append(realizationOutput.Electric.December, transactionTemp)
					} else if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.December = append(realizationOutput.Cement.December, transactionTemp)
					} else {
						realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
			}
		}
	}

	var groupingVessels []groupingvesseldn.GroupingVesselDn

	queryFilterGrouping := fmt.Sprintf("grouping_vessel_dns.iupopk_id = %v AND report_dmos.period LIKE '%%%v' AND grouping_vessel_dns.sales_system = 'Vessel'", iupopkId, year)

	errFindGrouping := r.db.Preload("ReportDmo").Preload("Buyer.IndustryType.CategoryIndustryType").Table("grouping_vessel_dns").Select("grouping_vessel_dns.*").Joins("left join report_dmos on report_dmos.id = grouping_vessel_dns.report_dmo_id").Where(queryFilterGrouping).Find(&groupingVessels).Error

	if errFindGrouping != nil {
		return realizationOutput, errFindGrouping
	}

	for _, v := range groupingVessels {
		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		var transactionTemp RealizationTransaction

		var tempTransaction transaction.Transaction

		errFind := r.db.Preload("Dmo").Preload("Customer.IndustryType.CategoryIndustryType").Where("grouping_vessel_dn_id = ? and seller_id = ?", v.ID, iupopkId).First(&tempTransaction).Error

		if errFind != nil {
			return realizationOutput, errFind
		}

		switch periodSplit[0] {
		case "Jan":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.January = append(realizationOutput.Electric.January, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.January = append(realizationOutput.Cement.January, transactionTemp)
					} else {
						realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
			}
		case "Feb":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.February = append(realizationOutput.Electric.February, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.February = append(realizationOutput.Cement.February, transactionTemp)
					} else {
						realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
			}
		case "Mar":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.March = append(realizationOutput.Electric.March, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.March = append(realizationOutput.Cement.March, transactionTemp)
					} else {
						realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
			}
		case "Apr":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.April = append(realizationOutput.Electric.April, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.April = append(realizationOutput.Cement.April, transactionTemp)
					} else {
						realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
			}
		case "May":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.May = append(realizationOutput.Electric.May, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.May = append(realizationOutput.Cement.May, transactionTemp)
					} else {
						realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
			}
		case "Jun":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.June = append(realizationOutput.Electric.June, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.June = append(realizationOutput.Cement.June, transactionTemp)
					} else {
						realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
			}
		case "Jul":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.July = append(realizationOutput.Electric.July, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.July = append(realizationOutput.Cement.July, transactionTemp)
					} else {
						realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
			}
		case "Aug":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.August = append(realizationOutput.Electric.August, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.August = append(realizationOutput.Cement.August, transactionTemp)
					} else {
						realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
			}
		case "Sep":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.September = append(realizationOutput.Electric.September, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.September = append(realizationOutput.Cement.September, transactionTemp)
					} else {
						realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
			}
		case "Oct":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.October = append(realizationOutput.Electric.October, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.October = append(realizationOutput.Cement.October, transactionTemp)
					} else {
						realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
			}
		case "Nov":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.November = append(realizationOutput.Electric.November, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.November = append(realizationOutput.Cement.November, transactionTemp)
					} else {
						realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
			}
		case "Dec":
			var blDate string
			blDateSplit := strings.Split(*v.BlDate, "T")
			blDate = blDateSplit[0]
			transactionTemp.ShippingDate = blDate
			if tempTransaction.Customer != nil {
				transactionTemp.Trader = tempTransaction.Customer
			}

			if v.Buyer != nil {
				transactionTemp.EndUser = v.Buyer
			}
			transactionTemp.Quantity = v.GrandTotalQuantity
			if v.QualityCaloriesAr != nil {
				transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
			}
			if tempTransaction.Dmo != nil && tempTransaction.Dmo.IsBastDocumentSigned {
				transactionTemp.IsBastOk = true
			} else {
				transactionTemp.IsBastOk = false
			}

			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						realizationOutput.Electric.December = append(realizationOutput.Electric.December, transactionTemp)
					} else if v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
						realizationOutput.Cement.December = append(realizationOutput.Cement.December, transactionTemp)
					} else {
						realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.December = append(realizationOutput.NonElectric.December, transactionTemp)
			}
		}
	}

	return realizationOutput, nil
}

func (r *repository) SaleDetailReport(year string, iupopkId int) (SaleDetail, error) {
	var saleDetail SaleDetail

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	companyElectricity := make(map[string][]string)
	companyCement := make(map[string][]string)
	companyNonElectricity := make(map[string][]string)

	var jetty []jetty.Jetty

	errFindJetty := r.db.Where("iupopk_id = ?", iupopkId).Order("name asc").Find(&jetty).Error

	if errFindJetty != nil {
		return saleDetail, errFindJetty
	}

	// Production Query
	var listProduction []production.Production

	queryFilterProduction := fmt.Sprintf("production_date >= '%s' AND production_date <= '%s' AND iupopk_id = %v", startFilter, endFilter, iupopkId)

	errFindProduction := r.db.Preload(clause.Associations).Where(queryFilterProduction).Order("id ASC").Find(&listProduction).Error

	if errFindProduction != nil {
		return saleDetail, errFindProduction
	}

	var salesJetty SalesJetty

	salesJetty.January = make(map[string]float64)
	salesJetty.February = make(map[string]float64)
	salesJetty.March = make(map[string]float64)
	salesJetty.April = make(map[string]float64)
	salesJetty.May = make(map[string]float64)
	salesJetty.June = make(map[string]float64)
	salesJetty.July = make(map[string]float64)
	salesJetty.August = make(map[string]float64)
	salesJetty.September = make(map[string]float64)
	salesJetty.October = make(map[string]float64)
	salesJetty.November = make(map[string]float64)
	salesJetty.December = make(map[string]float64)

	saleDetail.LossJetty.January = make(map[string]float64)
	saleDetail.LossJetty.February = make(map[string]float64)
	saleDetail.LossJetty.March = make(map[string]float64)
	saleDetail.LossJetty.April = make(map[string]float64)
	saleDetail.LossJetty.May = make(map[string]float64)
	saleDetail.LossJetty.June = make(map[string]float64)
	saleDetail.LossJetty.July = make(map[string]float64)
	saleDetail.LossJetty.August = make(map[string]float64)
	saleDetail.LossJetty.September = make(map[string]float64)
	saleDetail.LossJetty.October = make(map[string]float64)
	saleDetail.LossJetty.November = make(map[string]float64)
	saleDetail.LossJetty.December = make(map[string]float64)

	var productionJetty ProductionJetty

	productionJetty.January = make(map[string]float64)
	productionJetty.February = make(map[string]float64)
	productionJetty.March = make(map[string]float64)
	productionJetty.April = make(map[string]float64)
	productionJetty.May = make(map[string]float64)
	productionJetty.June = make(map[string]float64)
	productionJetty.July = make(map[string]float64)
	productionJetty.August = make(map[string]float64)
	productionJetty.September = make(map[string]float64)
	productionJetty.October = make(map[string]float64)
	productionJetty.November = make(map[string]float64)
	productionJetty.December = make(map[string]float64)

	for _, v := range listProduction {

		date, _ := time.Parse("2006-01-02T00:00:00Z", v.ProductionDate)
		_, month, _ := date.Date()
		switch int(month) {
		case 1:
			saleDetail.Production.January += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.January[v.Jetty.Name]; ok {
					productionJetty.January[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.January[v.Jetty.Name] = v.Quantity
				}

			}

		case 2:
			saleDetail.Production.February += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.February[v.Jetty.Name]; ok {
					productionJetty.February[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.February[v.Jetty.Name] = v.Quantity
				}

			}

		case 3:
			saleDetail.Production.March += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.March[v.Jetty.Name]; ok {
					productionJetty.March[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.March[v.Jetty.Name] = v.Quantity
				}

			}

		case 4:
			saleDetail.Production.April += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.April[v.Jetty.Name]; ok {
					productionJetty.April[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.April[v.Jetty.Name] = v.Quantity
				}

			}

		case 5:
			saleDetail.Production.May += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.May[v.Jetty.Name]; ok {
					productionJetty.May[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.May[v.Jetty.Name] = v.Quantity
				}

			}
		case 6:
			saleDetail.Production.June += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.June[v.Jetty.Name]; ok {
					productionJetty.June[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.June[v.Jetty.Name] = v.Quantity
				}

			}

		case 7:
			saleDetail.Production.July += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.July[v.Jetty.Name]; ok {
					productionJetty.July[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.July[v.Jetty.Name] = v.Quantity
				}

			}

		case 8:
			saleDetail.Production.August += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.August[v.Jetty.Name]; ok {
					productionJetty.August[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.August[v.Jetty.Name] = v.Quantity
				}

			}

		case 9:
			saleDetail.Production.September += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.September[v.Jetty.Name]; ok {
					productionJetty.September[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.September[v.Jetty.Name] = v.Quantity
				}

			}

		case 10:
			saleDetail.Production.October += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.October[v.Jetty.Name]; ok {
					productionJetty.October[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.October[v.Jetty.Name] = v.Quantity
				}

			}

		case 11:
			saleDetail.Production.November += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.November[v.Jetty.Name]; ok {
					productionJetty.November[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.November[v.Jetty.Name] = v.Quantity
				}

			}

		case 12:
			saleDetail.Production.December += v.Quantity
			saleDetail.Production.Total += v.Quantity

			if v.Jetty != nil {
				if !helperString(saleDetail.JettyList, v.Jetty.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.Jetty.Name)
				}

				if _, ok := productionJetty.December[v.Jetty.Name]; ok {
					productionJetty.December[v.Jetty.Name] += v.Quantity
				} else {
					productionJetty.December[v.Jetty.Name] = v.Quantity
				}

			}

		}
	}

	saleDetail.ProductionJetty = productionJetty

	// Query Transaction
	var listCurrentTransactions []transaction.Transaction

	queryCurrentFilter := fmt.Sprintf("seller_id = %v AND shipping_date >= '%s' AND shipping_date <= '%s'", iupopkId, startFilter, endFilter)

	errCurrentFind := r.db.Preload("LoadingPort").Where(queryCurrentFilter).Order("shipping_date ASC").Find(&listCurrentTransactions).Error

	if errCurrentFind != nil {
		return saleDetail, errCurrentFind
	}

	var listTransactions []transaction.Transaction

	queryFilter := fmt.Sprintf("transactions.seller_id = %v AND transactions.transaction_type = 'DN' AND dmos.period LIKE '%%%v' AND dmo_id IS NOT NULL", iupopkId, year)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Table("transactions").Select("transactions.*").Joins("left join dmos on dmos.id = transactions.dmo_id").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

	if errFind != nil {
		return saleDetail, errFind
	}

	var electricAssignment electricassignment.ElectricAssignment
	var electricAssignmentEndUser []electricassignmentenduser.ElectricAssignmentEndUser

	r.db.Preload(clause.Associations).Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&electricAssignment)

	if electricAssignment.ID != 0 {
		errFindElectricAssignmentEndUser := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ?", electricAssignment.ID).Find(&electricAssignmentEndUser).Error

		if errFindElectricAssignmentEndUser != nil {
			return saleDetail, errFindElectricAssignmentEndUser
		}
	}

	var cafAssignment cafassignment.CafAssignment
	var cafAssignmentEndUser []cafassignmentenduser.CafAssignmentEndUser

	r.db.Preload(clause.Associations).Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&cafAssignment)

	if cafAssignment.ID != 0 {
		errFindCafAssignmentEndUser := r.db.Preload(clause.Associations).Where("caf_assignment_id = ?", cafAssignment.ID).Find(&cafAssignmentEndUser).Error

		if errFindCafAssignmentEndUser != nil {
			return saleDetail, errFindCafAssignmentEndUser
		}
	}

	// Rkabs Query
	var rkabs []rkab.Rkab

	queryRkab := fmt.Sprintf("(year = '%s' OR year2 = '%s' OR year3 = '%s') AND iupopk_id = %v", year, year, year, iupopkId)

	errFindRkab := r.db.Where(queryRkab).Order("id DESC").Find(&rkabs).Error

	if errFindRkab != nil {
		return saleDetail, errFindRkab
	}

	saleDetail.DataDetailIndustry = make(map[string]map[string]map[string]map[string]float64)
	saleDetail.DataRecapIndustry = make(map[string]map[string]float64)
	saleDetail.Company = make(map[string]map[string][]string)

	// START FLEXIBLE
	var categoryIndustryType []categoryindustrytype.CategoryIndustryType

	errFindCategory := r.db.Find(&categoryIndustryType).Error

	if errFindCategory != nil {
		return saleDetail, errFindCategory
	}

	for _, v := range categoryIndustryType {
		saleDetail.DataDetailIndustry[v.SystemName] = make(map[string]map[string]map[string]float64)
		saleDetail.DataRecapIndustry[v.SystemName] = make(map[string]float64)
		saleDetail.Company[v.SystemName] = make(map[string][]string)

		saleDetail.DataDetailIndustry[v.SystemName]["january"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["february"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["march"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["april"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["may"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["june"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["july"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["august"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["september"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["october"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["november"] = make(map[string]map[string]float64)
		saleDetail.DataDetailIndustry[v.SystemName]["december"] = make(map[string]map[string]float64)

		saleDetail.DataRecapIndustry[v.SystemName]["january"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["february"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["march"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["april"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["may"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["june"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["july"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["august"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["september"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["october"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["november"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["december"] = 0
		saleDetail.DataRecapIndustry[v.SystemName]["total"] = 0

	}

	saleDetail.Rkabs = rkabs

	saleDetail.Electricity.January = make(map[string]map[string]float64)
	saleDetail.Electricity.February = make(map[string]map[string]float64)
	saleDetail.Electricity.March = make(map[string]map[string]float64)
	saleDetail.Electricity.April = make(map[string]map[string]float64)
	saleDetail.Electricity.May = make(map[string]map[string]float64)
	saleDetail.Electricity.June = make(map[string]map[string]float64)
	saleDetail.Electricity.July = make(map[string]map[string]float64)
	saleDetail.Electricity.August = make(map[string]map[string]float64)
	saleDetail.Electricity.September = make(map[string]map[string]float64)
	saleDetail.Electricity.October = make(map[string]map[string]float64)
	saleDetail.Electricity.November = make(map[string]map[string]float64)
	saleDetail.Electricity.December = make(map[string]map[string]float64)

	saleDetail.Cement.January = make(map[string]map[string]float64)
	saleDetail.Cement.February = make(map[string]map[string]float64)
	saleDetail.Cement.March = make(map[string]map[string]float64)
	saleDetail.Cement.April = make(map[string]map[string]float64)
	saleDetail.Cement.May = make(map[string]map[string]float64)
	saleDetail.Cement.June = make(map[string]map[string]float64)
	saleDetail.Cement.July = make(map[string]map[string]float64)
	saleDetail.Cement.August = make(map[string]map[string]float64)
	saleDetail.Cement.September = make(map[string]map[string]float64)
	saleDetail.Cement.October = make(map[string]map[string]float64)
	saleDetail.Cement.November = make(map[string]map[string]float64)
	saleDetail.Cement.December = make(map[string]map[string]float64)

	saleDetail.NonElectricity.January = make(map[string]map[string]float64)
	saleDetail.NonElectricity.February = make(map[string]map[string]float64)
	saleDetail.NonElectricity.March = make(map[string]map[string]float64)
	saleDetail.NonElectricity.April = make(map[string]map[string]float64)
	saleDetail.NonElectricity.May = make(map[string]map[string]float64)
	saleDetail.NonElectricity.June = make(map[string]map[string]float64)
	saleDetail.NonElectricity.July = make(map[string]map[string]float64)
	saleDetail.NonElectricity.August = make(map[string]map[string]float64)
	saleDetail.NonElectricity.September = make(map[string]map[string]float64)
	saleDetail.NonElectricity.October = make(map[string]map[string]float64)
	saleDetail.NonElectricity.November = make(map[string]map[string]float64)
	saleDetail.NonElectricity.December = make(map[string]map[string]float64)

	saleDetail.ElectricAssignment.Quantity = electricAssignment.GrandTotalQuantity + electricAssignment.GrandTotalQuantity2 + electricAssignment.GrandTotalQuantity3 + electricAssignment.GrandTotalQuantity4

	saleDetail.CafAssignment.Quantity = cafAssignment.GrandTotalQuantity + cafAssignment.GrandTotalQuantity2 + cafAssignment.GrandTotalQuantity3 + cafAssignment.GrandTotalQuantity4

	for _, v := range listCurrentTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

		switch int(month) {
		case 1:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.January += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.January += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.January[v.LoadingPort.Name]

				if ok {
					salesJetty.January[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.January[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.January["-"]

				if ok {
					salesJetty.January["-"] += v.Quantity
				} else {
					salesJetty.January["-"] = v.Quantity
				}
			}

		case 2:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.February += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.February += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.February[v.LoadingPort.Name]

				if ok {
					salesJetty.February[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.February[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.February["-"]

				if ok {
					salesJetty.February["-"] += v.Quantity
				} else {
					salesJetty.February["-"] = v.Quantity
				}
			}

		case 3:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.March += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.March += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.March[v.LoadingPort.Name]

				if ok {
					salesJetty.March[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.March[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.March["-"]

				if ok {
					salesJetty.March["-"] += v.Quantity
				} else {
					salesJetty.March["-"] = v.Quantity
				}
			}

		case 4:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.April += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.April += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.April[v.LoadingPort.Name]

				if ok {
					salesJetty.April[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.April[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.April["-"]

				if ok {
					salesJetty.April["-"] += v.Quantity
				} else {
					salesJetty.April["-"] = v.Quantity
				}
			}

		case 5:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.May += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.May += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.May[v.LoadingPort.Name]

				if ok {
					salesJetty.May[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.May[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.May["-"]

				if ok {
					salesJetty.May["-"] += v.Quantity
				} else {
					salesJetty.May["-"] = v.Quantity
				}
			}

		case 6:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.June += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.June += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.June[v.LoadingPort.Name]

				if ok {
					salesJetty.June[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.June[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.June["-"]

				if ok {
					salesJetty.June["-"] += v.Quantity
				} else {
					salesJetty.June["-"] = v.Quantity
				}
			}

		case 7:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.July += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.July += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}
			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.July[v.LoadingPort.Name]

				if ok {
					salesJetty.July[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.July[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.July["-"]

				if ok {
					salesJetty.July["-"] += v.Quantity
				} else {
					salesJetty.July["-"] = v.Quantity
				}
			}

		case 8:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.August += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.August += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.August[v.LoadingPort.Name]

				if ok {
					salesJetty.August[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.August[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.August["-"]

				if ok {
					salesJetty.August["-"] += v.Quantity
				} else {
					salesJetty.August["-"] = v.Quantity
				}
			}

		case 9:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.September += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.September += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.September[v.LoadingPort.Name]

				if ok {
					salesJetty.September[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.September[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.September["-"]

				if ok {
					salesJetty.September["-"] += v.Quantity
				} else {
					salesJetty.September["-"] = v.Quantity
				}
			}

		case 10:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.October += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.October += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.October[v.LoadingPort.Name]

				if ok {
					salesJetty.October[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.October[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.October["-"]

				if ok {
					salesJetty.October["-"] += v.Quantity
				} else {
					salesJetty.October["-"] = v.Quantity
				}
			}

		case 11:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.November += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.November += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.November[v.LoadingPort.Name]

				if ok {
					salesJetty.November[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.November[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.November["-"]

				if ok {
					salesJetty.November["-"] += v.Quantity
				} else {
					salesJetty.November["-"] = v.Quantity
				}
			}

		case 12:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.December += v.Quantity
				saleDetail.Domestic.Total += v.Quantity
			} else {
				saleDetail.Export.December += v.Quantity
				saleDetail.Export.Total += v.Quantity
			}

			if v.LoadingPort != nil {
				if !helperString(saleDetail.JettyList, v.LoadingPort.Name) {
					saleDetail.JettyList = append(saleDetail.JettyList, v.LoadingPort.Name)
				}

				_, ok := salesJetty.December[v.LoadingPort.Name]

				if ok {
					salesJetty.December[v.LoadingPort.Name] += v.Quantity
				} else {
					salesJetty.December[v.LoadingPort.Name] = v.Quantity
				}
			} else {
				_, ok := salesJetty.December["-"]

				if ok {
					salesJetty.December["-"] += v.Quantity
				} else {
					salesJetty.December["-"] = v.Quantity
				}
			}

		}
	}

	sort.Sort(sort.StringSlice(saleDetail.JettyList))

	saleDetail.SalesJetty = salesJetty
	for _, v := range listTransactions {
		if v.ReportDmoId == nil || (v.GroupingVesselDnId != nil && v.SalesSystem != nil && strings.Contains(v.SalesSystem.Name, "Vessel")) {
			continue
		}
		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		if v.IsNotClaim == false {
			if v.DmoId != nil {
				var isAdded = false
				for _, value := range electricAssignmentEndUser {
					if !isAdded {
						if v.DmoDestinationPortId != nil {
							if *v.DmoDestinationPortId == value.PortId {
								isAdded = true
								saleDetail.ElectricAssignment.RealizationQuantity += v.Quantity
							}
						}
					}
				}

				for _, value := range cafAssignmentEndUser {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.CompanyName == value.EndUserString {
							saleDetail.CafAssignment.RealizationQuantity += v.Quantity
						}
					}
				}
			}

			switch periodSplit[0] {
			case "Jan":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.January += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.January += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.January += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["january"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}
					} else {
						if _, ok := saleDetail.NonElectricity.January["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.January["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.January["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.January["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}

						saleDetail.RecapNonElectricity.January += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["january"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Feb":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.February += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.February += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.February += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}
						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["february"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}
					} else {
						if _, ok := saleDetail.NonElectricity.February["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.February["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.February["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.February["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.February += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["february"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Mar":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.March += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.March += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.March += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["march"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.March["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.March["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.March["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.March["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.March += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["march"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Apr":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.April += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.April += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.April += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["april"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.April["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.April["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.April["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.April["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.April += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["april"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "May":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.May += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.May += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.May += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["may"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.May["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.May["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.May["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.May["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.May += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["may"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Jun":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.June += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.June += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.June += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["june"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.June["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.June["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.June["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.June["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.June += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["june"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Jul":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.July += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.July += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.July += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["july"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.July["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.July["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.July["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.July["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.July += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]["-"] += v.Quantity
						}
						saleDetail.DataRecapIndustry["non_electricity"]["july"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Aug":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.August += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.August += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.August += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["august"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.August["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.August["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.August["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.August["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.August += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["august"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Sep":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.September += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.September += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.September += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									fmt.Println(v.Quantity, v.DmoDestinationPort.Name)

									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}

									fmt.Println(saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"])
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["september"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.September["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.September["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.September["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.September["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.September += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["september"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Oct":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.October += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.October += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.October += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["october"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.October["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.October["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.October["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.October["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.October += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["october"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Nov":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.November += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.November += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.November += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["november"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.November["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.November["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.November["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.November["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.November += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]["-"] += v.Quantity
						}
						saleDetail.DataRecapIndustry["non_electricity"]["november"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}

			case "Dec":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Electricity.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.Quantity

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapElectricity.December += v.Quantity
							saleDetail.RecapElectricity.Total += v.Quantity
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.Cement.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapCement.December += v.Quantity
							saleDetail.RecapCement.Total += v.Quantity
						} else {
							if _, ok := saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							} else {
								saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.Quantity

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]["-"] += v.Quantity
								}
							}
							saleDetail.RecapNonElectricity.December += v.Quantity
							saleDetail.RecapNonElectricity.Total += v.Quantity
						}

						if v.DmoBuyer.IndustryType != nil {
							if v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
								if v.DmoDestinationPort != nil {
									if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
										if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.Quantity
									}
								} else {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-") {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], "-")
									}

									if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName]["-"]; okDestination {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName]["-"] += v.Quantity
									} else {
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName] = make(map[string]float64)
										saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName]["-"] = v.Quantity
									}
								}
							} else {
								if _, ok := saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName]; ok {
									if !helperString(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
										saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
								} else {
									saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName] = append(saleDetail.Company[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName][v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
								}

								if _, okDestination := saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory]; okDestination {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.Quantity
								} else {
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName] = make(map[string]float64)
									saleDetail.DataDetailIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] = v.Quantity
								}
							}

							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["december"] += v.Quantity
							saleDetail.DataRecapIndustry[v.DmoBuyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.Quantity
						}

					} else {
						if _, ok := saleDetail.NonElectricity.December["-"]; ok {
							saleDetail.NonElectricity.Total += v.Quantity
							saleDetail.NonElectricity.December["-"]["-"] += v.Quantity
						} else {
							saleDetail.NonElectricity.December["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.December["-"]["-"] += v.Quantity
							saleDetail.NonElectricity.Total += v.Quantity
						}
						saleDetail.RecapNonElectricity.December += v.Quantity
						saleDetail.RecapNonElectricity.Total += v.Quantity

						if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]; ok {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}
							saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]["-"] += v.Quantity
						} else {
							if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
								saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
							}

							saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"] = make(map[string]float64)
							saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]["-"] += v.Quantity
						}

						saleDetail.DataRecapIndustry["non_electricity"]["december"] += v.Quantity
						saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.Quantity
					}
				}
			}
		} else {
			switch periodSplit[0] {
			case "Jan":
				saleDetail.NotClaimable.January += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Feb":
				saleDetail.NotClaimable.February += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Mar":
				saleDetail.NotClaimable.March += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Apr":
				saleDetail.NotClaimable.April += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "May":
				saleDetail.NotClaimable.May += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Jun":
				saleDetail.NotClaimable.June += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Jul":
				saleDetail.NotClaimable.July += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Aug":
				saleDetail.NotClaimable.August += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Sep":
				saleDetail.NotClaimable.September += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Oct":
				saleDetail.NotClaimable.October += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Nov":
				saleDetail.NotClaimable.November += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			case "Dec":
				saleDetail.NotClaimable.December += v.Quantity
				saleDetail.NotClaimable.Total += v.Quantity
			}
		}
	}

	var groupingVessels []groupingvesseldn.GroupingVesselDn

	queryFilterGrouping := fmt.Sprintf("grouping_vessel_dns.iupopk_id = %v AND report_dmos.period LIKE '%%%v' AND grouping_vessel_dns.sales_system = 'Vessel'", iupopkId, year)

	errFindGrouping := r.db.Preload("ReportDmo").Preload("Buyer.IndustryType.CategoryIndustryType").Table("grouping_vessel_dns").Select("grouping_vessel_dns.*").Joins("left join report_dmos on report_dmos.id = grouping_vessel_dns.report_dmo_id").Where(queryFilterGrouping).Find(&groupingVessels).Error

	if errFindGrouping != nil {
		return saleDetail, errFindGrouping
	}

	for _, v := range groupingVessels {
		if v.ReportDmoId == nil {
			continue
		}

		// var tempTransaction transaction.Transaction

		// errFind := r.db.Preload("Dmo").Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("grouping_vessel_dn_id = ? and seller_id = ?", v.ID, iupopkId).First(&tempTransaction).Error

		// if errFind != nil {
		// 	return saleDetail, errFind
		// }

		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		var isAdded = false
		for _, value := range electricAssignmentEndUser {
			if !isAdded {
				if v.DmoDestinationPortId != nil {
					if *v.DmoDestinationPortId == value.PortId {
						isAdded = true
						saleDetail.ElectricAssignment.RealizationQuantity += v.GrandTotalQuantity
					}
				}
			}
		}

		for _, value := range cafAssignmentEndUser {
			if v.Buyer != nil {
				if v.Buyer.CompanyName == value.EndUserString {
					saleDetail.CafAssignment.RealizationQuantity += v.GrandTotalQuantity
				}
			}
		}

		switch periodSplit[0] {
		case "Jan":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.January[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.January[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.January[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.January[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.January[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.January[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.January += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.January[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.January[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.January[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.January[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.January += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.January[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.January[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.January[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.January[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.January[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.January += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["january"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}
			} else {
				if _, ok := saleDetail.NonElectricity.January["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.January["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.January["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.January["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.January += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["january"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["january"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Feb":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.February[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.February[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.February[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.February[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.February[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.February[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.February += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.February[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.February[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.February[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.February[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.February += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.February[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.February[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.February[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.February[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.February[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.February += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["february"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.February["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.February["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.February["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.February["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.February += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["february"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["february"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Mar":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.March[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.March[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.March[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.March[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.March[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.March[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.March += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.March[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.March[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.March[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.March[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.March += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.March[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.March[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.March[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.March[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.March[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.March += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["march"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.March["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.March["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.March["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.March["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.March += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["march"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["march"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Apr":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.April[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.April[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.April[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.April[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.April[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.April[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.April += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.April[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.April[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.April[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.April[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.April += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.April[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.April[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.April[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.April[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.April[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.April += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["april"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.April["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.April["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.April["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.April["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.April += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["april"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["april"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "May":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.May[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.May[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.May[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.May[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.May[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.May[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.May += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.May[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.May[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.May[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.May[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.May += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.May[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.May[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.May[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.May[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.May[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.May += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["may"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.May["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.May["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.May["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.May["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.May += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["may"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["may"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Jun":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.June[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.June[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.June[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.June[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.June[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.June[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.June += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.June[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.June[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.June[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.June[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.June += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.June[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.June[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.June[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.June[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.June[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.June += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["june"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.June["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.June["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.June["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.June["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.June += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["june"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["june"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Jul":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.July[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.July[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.July[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.July[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.July[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.July[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.July += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.July[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.July[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.July[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.July[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.July += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.July[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.July[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.July[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.July[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.July[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.July += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["july"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.July["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.July["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.July["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.July["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.July += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["july"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["july"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Aug":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.August[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.August[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.August[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.August[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.August[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.August[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.August += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.August[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.August[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.August[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.August[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.August += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.August[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.August[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.August[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.August[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.August[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.August += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["august"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.August["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.August["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.August["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.August["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.August += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["august"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["august"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Sep":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.September[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.September[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.September[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.September[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.September[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.September[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.September += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.September[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.September[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.September[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.September[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.September += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.September[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.September[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.September[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.September[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.September[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.September += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["september"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.September["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.September["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.September["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.September["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.September += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["september"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["september"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Oct":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.October[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.October[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.October[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.October[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.October[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.October[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.October += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.October[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.October[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.October[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.October[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.October += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.October[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.October[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.October[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.October[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.October[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.October += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["october"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.October["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.October["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.October["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.October["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.October += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["october"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["october"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Nov":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.November[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.November[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.November[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.November[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.November[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.November[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.November += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.November[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.November[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.November[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.November[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.November += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.November[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.November[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.November[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.November[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.November[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.November += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["november"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.November["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.November["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.November["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.November["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.November += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["november"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["november"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}

		case "Dec":
			if v.Buyer != nil {
				if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					if _, ok := saleDetail.Electricity.December[v.Buyer.CompanyName]; ok {
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.Electricity.December[v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.Electricity.December[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.Electricity.December[v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.Electricity.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Electricity.December[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Electricity.Total += v.GrandTotalQuantity

						if v.DmoDestinationPort != nil {
							if _, ok := companyElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}
							saleDetail.Electricity.December[v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
						} else {
							if !helperString(companyElectricity[v.Buyer.CompanyName], "-") {
								companyElectricity[v.Buyer.CompanyName] = append(companyElectricity[v.Buyer.CompanyName], "-")
							}
							saleDetail.Electricity.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapElectricity.December += v.GrandTotalQuantity
					saleDetail.RecapElectricity.Total += v.GrandTotalQuantity
				} else if v.Buyer.IndustryType != nil && v.Buyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					if _, ok := saleDetail.Cement.December[v.Buyer.CompanyName]; ok {
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.Cement.December[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.Cement.December[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.Cement.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyCement[v.Buyer.CompanyName]; ok {
								if !helperString(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.Cement.December[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyCement[v.Buyer.CompanyName], "-") {
								companyCement[v.Buyer.CompanyName] = append(companyCement[v.Buyer.CompanyName], "-")
							}

							saleDetail.Cement.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapCement.December += v.GrandTotalQuantity
					saleDetail.RecapCement.Total += v.GrandTotalQuantity
				} else {
					if _, ok := saleDetail.NonElectricity.December[v.Buyer.CompanyName]; ok {
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}

							saleDetail.NonElectricity.December[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					} else {
						saleDetail.NonElectricity.December[v.Buyer.CompanyName] = make(map[string]float64)
						saleDetail.NonElectricity.Total += v.GrandTotalQuantity

						if v.Buyer.IndustryType != nil {
							if _, ok := companyNonElectricity[v.Buyer.CompanyName]; ok {
								if !helperString(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
									companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
								}
							} else {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
							saleDetail.NonElectricity.December[v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							if !helperString(companyNonElectricity[v.Buyer.CompanyName], "-") {
								companyNonElectricity[v.Buyer.CompanyName] = append(companyNonElectricity[v.Buyer.CompanyName], "-")
							}

							saleDetail.NonElectricity.December[v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
						}
					}
					saleDetail.RecapNonElectricity.December += v.GrandTotalQuantity
					saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity
				}

				if v.Buyer.IndustryType != nil {
					if v.Buyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
						if v.DmoDestinationPort != nil {
							if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
								if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name) {
									saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
								}
							} else {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.DmoDestinationPort.Name)
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.DmoDestinationPort.Name] = v.GrandTotalQuantity
							}
						} else {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-") {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], "-")
							}

							if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName]["-"]; okDestination {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName]["-"] += v.GrandTotalQuantity
							} else {
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName] = make(map[string]float64)
								saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName]["-"] = v.GrandTotalQuantity
							}
						}
					} else {
						if _, ok := saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName]; ok {
							if !helperString(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory) {
								saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
							}
						} else {
							saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName] = append(saleDetail.Company[v.Buyer.IndustryType.CategoryIndustryType.SystemName][v.Buyer.CompanyName], v.Buyer.IndustryType.SystemCategory)
						}

						if _, okDestination := saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory]; okDestination {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] += v.GrandTotalQuantity
						} else {
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName] = make(map[string]float64)
							saleDetail.DataDetailIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"][v.Buyer.CompanyName][v.Buyer.IndustryType.SystemCategory] = v.GrandTotalQuantity
						}
					}

					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["december"] += v.GrandTotalQuantity
					saleDetail.DataRecapIndustry[v.Buyer.IndustryType.CategoryIndustryType.SystemName]["total"] += v.GrandTotalQuantity
				}

			} else {
				if _, ok := saleDetail.NonElectricity.December["-"]; ok {
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
					saleDetail.NonElectricity.December["-"]["-"] += v.GrandTotalQuantity
				} else {
					saleDetail.NonElectricity.December["-"] = make(map[string]float64)
					if !helperString(companyNonElectricity["-"], "-") {
						companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
					}
					saleDetail.NonElectricity.December["-"]["-"] += v.GrandTotalQuantity
					saleDetail.NonElectricity.Total += v.GrandTotalQuantity
				}
				saleDetail.RecapNonElectricity.December += v.GrandTotalQuantity
				saleDetail.RecapNonElectricity.Total += v.GrandTotalQuantity

				if _, ok := saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]; ok {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}
					saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]["-"] += v.GrandTotalQuantity
				} else {
					if !helperString(saleDetail.Company["non_electricity"]["-"], "-") {
						saleDetail.Company["non_electricity"]["-"] = append(saleDetail.Company["non_electricity"]["-"], "-")
					}

					saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"] = make(map[string]float64)
					saleDetail.DataDetailIndustry["non_electricity"]["december"]["-"]["-"] += v.GrandTotalQuantity
				}

				saleDetail.DataRecapIndustry["non_electricity"]["december"] += v.GrandTotalQuantity
				saleDetail.DataRecapIndustry["non_electricity"]["total"] += v.GrandTotalQuantity
			}
		}
	}

	var jettyBalanceAndLoss []JettyBalanceLoss

	rawQueryJettyBalance := fmt.Sprintf("Select jb.id as id, jb.jetty_id as jetty_id, j.* as jetty, start_balance, total_loss from jetty_balances jb left join jetties j on j.id = jb.jetty_id where jb.iupopk_id = %v and year = '%v'", iupopkId, year)

	errFindJettyBalance := r.db.Preload(clause.Associations).Raw(rawQueryJettyBalance).Find(&jettyBalanceAndLoss).Error

	if errFindJettyBalance != nil {
		return saleDetail, errFindJettyBalance
	}

	var jettyOutputBalance []JettyBalanceLoss

	for _, v := range jetty {
		var jettyBalance JettyBalanceLoss

		rawQueryJettyBalance := fmt.Sprintf("Select jb.id as id, jb.jetty_id as jetty_id, start_balance, total_loss from jetty_balances jb left join jetties j on j.id = jb.jetty_id where jb.iupopk_id = %v and year = '%v' and jb.jetty_id = %v", iupopkId, year, v.ID)

		errJettyBalance := r.db.Preload(clause.Associations).Raw(rawQueryJettyBalance).First(&jettyBalance).Error

		if errJettyBalance != nil {
			jettyBalance.JettyId = v.ID
			jettyBalance.Jetty = v
			jettyBalance.TotalLoss = 0
		} else {
			var pitLoss []pitloss.PitLoss

			errFindPitLoss := r.db.Where("jetty_balance_id = ?", jettyBalance.ID).Find(&pitLoss).Error

			if errFindPitLoss != nil {
				return saleDetail, errFindPitLoss
			}

			for _, p := range pitLoss {
				if _, ok := saleDetail.LossJetty.January[v.Name]; ok {
					saleDetail.LossJetty.January[v.Name] += p.JanuaryLossQuantity
				} else {
					saleDetail.LossJetty.January[v.Name] = p.JanuaryLossQuantity
				}

				if _, ok := saleDetail.LossJetty.February[v.Name]; ok {
					saleDetail.LossJetty.February[v.Name] += p.FebruaryLossQuantity
				} else {
					saleDetail.LossJetty.February[v.Name] = p.FebruaryLossQuantity
				}

				if _, ok := saleDetail.LossJetty.March[v.Name]; ok {
					saleDetail.LossJetty.March[v.Name] += p.MarchLossQuantity
				} else {
					saleDetail.LossJetty.March[v.Name] = p.MarchLossQuantity
				}

				if _, ok := saleDetail.LossJetty.April[v.Name]; ok {
					saleDetail.LossJetty.April[v.Name] += p.AprilLossQuantity
				} else {
					saleDetail.LossJetty.April[v.Name] = p.AprilLossQuantity
				}

				if _, ok := saleDetail.LossJetty.May[v.Name]; ok {
					saleDetail.LossJetty.May[v.Name] += p.MayLossQuantity
				} else {
					saleDetail.LossJetty.May[v.Name] = p.MayLossQuantity
				}

				if _, ok := saleDetail.LossJetty.June[v.Name]; ok {
					saleDetail.LossJetty.June[v.Name] += p.JuneLossQuantity
				} else {
					saleDetail.LossJetty.June[v.Name] = p.JuneLossQuantity
				}

				if _, ok := saleDetail.LossJetty.July[v.Name]; ok {
					saleDetail.LossJetty.July[v.Name] += p.JulyLossQuantity
				} else {
					saleDetail.LossJetty.July[v.Name] = p.JulyLossQuantity
				}

				if _, ok := saleDetail.LossJetty.August[v.Name]; ok {
					saleDetail.LossJetty.August[v.Name] += p.AugustLossQuantity
				} else {
					saleDetail.LossJetty.August[v.Name] = p.AugustLossQuantity
				}

				if _, ok := saleDetail.LossJetty.October[v.Name]; ok {
					saleDetail.LossJetty.October[v.Name] += p.OctoberLossQuantity
				} else {
					saleDetail.LossJetty.October[v.Name] = p.OctoberLossQuantity
				}

				if _, ok := saleDetail.LossJetty.November[v.Name]; ok {
					saleDetail.LossJetty.November[v.Name] += p.NovemberLossQuantity
				} else {
					saleDetail.LossJetty.November[v.Name] = p.NovemberLossQuantity
				}

				if _, ok := saleDetail.LossJetty.December[v.Name]; ok {
					saleDetail.LossJetty.December[v.Name] += p.DecemberLossQuantity
				} else {
					saleDetail.LossJetty.December[v.Name] = p.DecemberLossQuantity
				}
			}
		}

		var production float64
		var sales float64
		var loss float64
		var quantityProduction *float64
		var quantitySales *float64
		var quantityLoss *float64

		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return saleDetail, err
		}

		errProd := r.db.Table("productions").Select("SUM(quantity)").Where("iupopk_id = ? and jetty_id = ? and production_date <= ?", iupopkId, v.ID, fmt.Sprintf("%v-12-31", yearInt-1)).Scan(&quantityProduction).Error
		if errProd != nil {
			return saleDetail, errProd
		}

		errSales := r.db.Table("transactions").Select("SUM(quantity)").Where("seller_id = ? and loading_port_id = ? and shipping_date <= ?", iupopkId, v.ID, fmt.Sprintf("%v-12-31", yearInt-1)).Scan(&quantitySales).Error

		if errProd != nil {
			return saleDetail, errSales
		}

		errLoss := r.db.Table("jetty_balances").Select("SUM(total_loss)").Where("iupopk_id = ? and jetty_id = ? and cast(year AS INTEGER) < ?", iupopkId, v.ID, year).Scan(&quantityLoss).Error

		if errLoss != nil {
			return saleDetail, errLoss
		}

		if quantityProduction != nil {
			production = *quantityProduction
		}

		if quantitySales != nil {
			sales = *quantitySales
		}

		if quantityLoss != nil {
			loss = *quantityLoss
		}

		jettyBalance.StartBalance = production - sales - loss

		jettyOutputBalance = append(jettyOutputBalance, jettyBalance)
	}

	saleDetail.JettyBalanceLoss = jettyOutputBalance
	saleDetail.CompanyElectricity = companyElectricity
	saleDetail.CompanyCement = companyCement
	saleDetail.CompanyNonElectricity = companyNonElectricity
	return saleDetail, nil
}

// Transaction report
func (r *repository) GetTransactionReport(iupopkId int, input TransactionReportInput, typeTransaction string) ([]TransactionReport, error) {
	var transactionReport []TransactionReport

	errFind := r.db.Table("transactions").Preload(clause.Associations).Preload("LoadingPort.Iupopk").Preload("UnloadingPort.PortLocation").Preload("DmoDestinationPort.PortLocation").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("seller_id = ? and transaction_type = ? and shipping_date >= ? and shipping_date <= ?", iupopkId, strings.ToUpper(typeTransaction), input.DateFrom, input.DateTo).Find(&transactionReport).Error

	if errFind != nil {
		return transactionReport, errFind
	}

	return transactionReport, nil
}
