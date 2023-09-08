package masterreport

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/electricassignment"
	"ajebackend/model/electricassignmentenduser"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/production"
	"ajebackend/model/rkab"
	"ajebackend/model/transaction"
	"fmt"
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
					reportDmoOuput.RecapElectricity.January += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.January += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.January += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Feb":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.February += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Mar":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.March += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Apr":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.April += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "May":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.May += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Jun":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.June += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Jul":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.July += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Aug":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.August += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Sep":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.September += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Oct":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.October += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Nov":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.November += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case "Dec":
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					reportDmoOuput.RecapElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
					reportDmoOuput.RecapCement.December += v.QuantityUnloading
					reportDmoOuput.RecapCement.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Smelter" {
					reportDmoOuput.RecapNonElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			}
		} else {
			switch periodSplit[0] {
			case "Jan":
				reportDmoOuput.NotClaimable.January += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Feb":
				reportDmoOuput.NotClaimable.February += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Mar":
				reportDmoOuput.NotClaimable.March += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Apr":
				reportDmoOuput.NotClaimable.April += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "May":
				reportDmoOuput.NotClaimable.May += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Jun":
				reportDmoOuput.NotClaimable.June += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Jul":
				reportDmoOuput.NotClaimable.July += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Aug":
				reportDmoOuput.NotClaimable.August += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Sep":
				reportDmoOuput.NotClaimable.September += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Oct":
				reportDmoOuput.NotClaimable.October += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Nov":
				reportDmoOuput.NotClaimable.November += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case "Dec":
				reportDmoOuput.NotClaimable.December += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			}
		}
	}

	var groupingVessels []groupingvesseldn.GroupingVesselDn

	queryFilterGrouping := fmt.Sprintf("grouping_vessel_dns.iupopk_id = %v AND report_dmos.period LIKE '%%%v' AND grouping_vessel_dns.sales_system = 'Vessel'", iupopkId, year)

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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.Quantity = v.QuantityUnloading
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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
			transactionTemp.QualityCaloriesAr = *v.QualityCaloriesAr
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

	// Production Query
	var listProduction []production.Production

	queryFilterProduction := fmt.Sprintf("production_date >= '%s' AND production_date <= '%s' AND iupopk_id = %v", startFilter, endFilter, iupopkId)

	errFindProduction := r.db.Where(queryFilterProduction).Order("id ASC").Find(&listProduction).Error

	if errFindProduction != nil {
		return saleDetail, errFindProduction
	}

	for _, v := range listProduction {
		date, _ := time.Parse("2006-01-02T00:00:00Z", v.ProductionDate)
		_, month, _ := date.Date()
		switch int(month) {
		case 1:
			saleDetail.Production.January += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 2:
			saleDetail.Production.February += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 3:
			saleDetail.Production.March += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 4:
			saleDetail.Production.April += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 5:
			saleDetail.Production.May += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 6:
			saleDetail.Production.June += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 7:
			saleDetail.Production.July += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 8:
			saleDetail.Production.August += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 9:
			saleDetail.Production.September += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 10:
			saleDetail.Production.October += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 11:
			saleDetail.Production.November += v.Quantity
			saleDetail.Production.Total += v.Quantity
		case 12:
			saleDetail.Production.December += v.Quantity
			saleDetail.Production.Total += v.Quantity
		}
	}

	// Query Transaction
	var listCurrentTransactions []transaction.Transaction

	queryCurrentFilter := fmt.Sprintf("seller_id = %v AND shipping_date >= '%s' AND shipping_date <= '%s'", iupopkId, startFilter, endFilter)

	errCurrentFind := r.db.Where(queryCurrentFilter).Order("shipping_date ASC").Find(&listCurrentTransactions).Error

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

	queryRkab := fmt.Sprintf("year = '%s' AND iupopk_id = %v", year, iupopkId)

	errFindRkab := r.db.Where(queryRkab).Order("id ASC").Find(&rkabs).Error

	if errFindRkab != nil {
		return saleDetail, errFindRkab
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
				saleDetail.Domestic.January += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.January += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 2:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.February += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.February += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 3:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.March += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.March += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 4:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.April += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.April += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 5:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.May += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.May += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 6:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.June += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.June += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 7:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.July += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.July += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 8:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.August += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.August += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 9:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.September += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.September += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 10:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.October += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.October += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 11:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.November += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.November += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}

		case 12:
			if v.TransactionType == "DN" {
				saleDetail.Domestic.December += v.QuantityUnloading
				saleDetail.Domestic.Total += v.QuantityUnloading
			} else {
				saleDetail.Export.December += v.QuantityUnloading
				saleDetail.Export.Total += v.QuantityUnloading
			}
		}
	}

	for _, v := range listTransactions {
		if v.DmoId != nil && v.SalesSystem != nil && !strings.Contains(v.SalesSystem.Name, "Vessel") {
			var isAdded = false
			for _, value := range electricAssignmentEndUser {
				if !isAdded {
					if v.CustomerId != nil && value.SupplierId != nil && v.DmoDestinationPortId != nil {
						if v.Customer.CompanyName == value.Supplier.CompanyName && *v.DmoDestinationPortId == value.PortId {
							isAdded = true
							saleDetail.ElectricAssignment.RealizationQuantity += v.QuantityUnloading

						}
					} else {
						if v.DmoDestinationPortId != nil {
							if v.CustomerId == nil && value.SupplierId == nil && *v.DmoDestinationPortId == value.PortId {
								isAdded = true
								saleDetail.ElectricAssignment.RealizationQuantity += v.QuantityUnloading

							}
						}
					}
				}
			}

			for _, value := range cafAssignmentEndUser {
				if v.DmoBuyer != nil {
					if v.DmoBuyer.CompanyName == value.EndUserString {
						saleDetail.CafAssignment.RealizationQuantity += v.QuantityUnloading
					}
				}
			}
		}

		if v.ReportDmoId == nil || (v.GroupingVesselDnId != nil && v.SalesSystem != nil && strings.Contains(v.SalesSystem.Name, "Vessel")) {
			continue
		}

		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		if v.IsNotClaim == false {
			if v.DmoId != nil {
				var isAdded = false
				for _, value := range electricAssignmentEndUser {
					if !isAdded {
						if v.CustomerId != nil && value.SupplierId != nil && v.DmoDestinationPortId != nil {
							if v.Customer.CompanyName == value.Supplier.CompanyName && *v.DmoDestinationPortId == value.PortId {
								isAdded = true
								saleDetail.ElectricAssignment.RealizationQuantity += v.QuantityUnloading

							}
						} else {
							if v.DmoDestinationPortId != nil {
								if v.CustomerId == nil && value.SupplierId == nil && *v.DmoDestinationPortId == value.PortId {
									isAdded = true
									saleDetail.ElectricAssignment.RealizationQuantity += v.QuantityUnloading

								}
							}
						}
					}
				}

				for _, value := range cafAssignmentEndUser {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.CompanyName == value.EndUserString {
							saleDetail.CafAssignment.RealizationQuantity += v.QuantityUnloading
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
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.January[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.January += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.January += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.January += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.January["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.January["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.January["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.January["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.January += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Feb":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.February[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.February += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.February += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.February += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.February["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.February["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.February["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.February["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.February += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Mar":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.March[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.March += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.March += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.March += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.March["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.March["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.March["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.March["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.March += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Apr":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.April[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.April += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.April += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.April += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.April["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.April["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.April["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.April["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.April += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "May":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.May[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.May += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.May += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.May += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.May["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.May["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.May["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.May["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.May += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Jun":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.June[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.June += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.June += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.June += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.June["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.June["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.June["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.June["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.June += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Jul":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.July[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.July += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.July += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.July += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.July["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.July["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.July["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.July["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.July += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Aug":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.August[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.August += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.August += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.August += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.August["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.August["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.August["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.August["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.August += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Sep":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.September[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.September += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.September += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.September += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.September["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.September["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.September["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.September["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.September += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Oct":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.October[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.October += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.October += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.October += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.October["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.October["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.October["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.October["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.October += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Nov":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.November[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.November += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.November += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.November += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.November["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.November["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.November["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.November["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.November += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}

			case "Dec":
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
							if _, ok := saleDetail.Electricity.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}

									if _, okDestination := saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name]; okDestination {
										saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
									} else {
										saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] = v.QuantityUnloading
									}
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Electricity.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Electricity.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Electricity.Total += v.QuantityUnloading

								if v.DmoDestinationPort != nil {
									if _, ok := companyElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name) {
											companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
										}
									} else {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], v.DmoDestinationPort.Name)
									}
									saleDetail.Electricity.December[v.DmoBuyer.CompanyName][v.DmoDestinationPort.Name] += v.QuantityUnloading
								} else {
									if !helperString(companyElectricity[v.DmoBuyer.CompanyName], "-") {
										companyElectricity[v.DmoBuyer.CompanyName] = append(companyElectricity[v.DmoBuyer.CompanyName], "-")
									}
									saleDetail.Electricity.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapElectricity.December += v.QuantityUnloading
							saleDetail.RecapElectricity.Total += v.QuantityUnloading
						} else if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.CategoryIndustryType.Name == "Semen" {
							if _, ok := saleDetail.Cement.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.Cement.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.Cement.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyCement[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.Cement.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyCement[v.DmoBuyer.CompanyName], "-") {
										companyCement[v.DmoBuyer.CompanyName] = append(companyCement[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.Cement.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapCement.December += v.QuantityUnloading
							saleDetail.RecapCement.Total += v.QuantityUnloading
						} else {
							if _, ok := saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							} else {
								saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName] = make(map[string]float64)
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.SystemCategory)
									}
									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.SystemCategory] += v.QuantityUnloading
								} else {
									if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], "-") {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], "-")
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]["-"] += v.QuantityUnloading
								}
							}
							saleDetail.RecapNonElectricity.December += v.QuantityUnloading
							saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
						}
					} else {
						if _, ok := saleDetail.NonElectricity.December["-"]; ok {
							saleDetail.NonElectricity.Total += v.QuantityUnloading
							saleDetail.NonElectricity.December["-"]["-"] += v.QuantityUnloading
						} else {
							saleDetail.NonElectricity.December["-"] = make(map[string]float64)
							if !helperString(companyNonElectricity["-"], "-") {
								companyNonElectricity["-"] = append(companyNonElectricity["-"], "-")
							}
							saleDetail.NonElectricity.December["-"]["-"] += v.QuantityUnloading
							saleDetail.NonElectricity.Total += v.QuantityUnloading
						}
						saleDetail.RecapNonElectricity.December += v.QuantityUnloading
						saleDetail.RecapNonElectricity.Total += v.QuantityUnloading
					}
				}
			}
		} else {
			switch periodSplit[0] {
			case "Jan":
				saleDetail.NotClaimable.January += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Feb":
				saleDetail.NotClaimable.February += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Mar":
				saleDetail.NotClaimable.March += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Apr":
				saleDetail.NotClaimable.April += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "May":
				saleDetail.NotClaimable.May += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Jun":
				saleDetail.NotClaimable.June += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Jul":
				saleDetail.NotClaimable.July += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Aug":
				saleDetail.NotClaimable.August += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Sep":
				saleDetail.NotClaimable.September += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Oct":
				saleDetail.NotClaimable.October += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Nov":
				saleDetail.NotClaimable.November += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case "Dec":
				saleDetail.NotClaimable.December += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
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

		var tempTransaction transaction.Transaction

		errFind := r.db.Preload("Dmo").Preload("Customer.IndustryType.CategoryIndustryType").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("grouping_vessel_dn_id = ? and seller_id = ?", v.ID, iupopkId).First(&tempTransaction).Error

		if errFind != nil {
			return saleDetail, errFind
		}

		periodSplit := strings.Split(v.ReportDmo.Period, " ")

		var isAdded = false
		for _, value := range electricAssignmentEndUser {
			if !isAdded {
				if tempTransaction.CustomerId != nil && value.SupplierId != nil && v.DmoDestinationPortId != nil {
					if tempTransaction.Customer.CompanyName == value.Supplier.CompanyName && *v.DmoDestinationPortId == value.PortId {
						isAdded = true
						saleDetail.ElectricAssignment.RealizationQuantity += v.Quantity

					}
				} else {
					if v.DmoDestinationPortId != nil {
						if tempTransaction.CustomerId == nil && value.SupplierId == nil && *v.DmoDestinationPortId == value.PortId {
							isAdded = true
							saleDetail.ElectricAssignment.RealizationQuantity += v.Quantity
						}
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
			}
		}

	}

	saleDetail.CompanyElectricity = companyElectricity
	saleDetail.CompanyCement = companyCement
	saleDetail.CompanyNonElectricity = companyNonElectricity
	return saleDetail, nil
}

// Transaction report
func (r *repository) GetTransactionReport(iupopkId int, input TransactionReportInput, typeTransaction string) ([]TransactionReport, error) {
	var transactionReport []TransactionReport

	errFind := r.db.Table("transactions").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoDestinationPort.PortLocation").Preload("DmoBuyer.IndustryType.CategoryIndustryType").Where("seller_id = ? and transaction_type = ? and shipping_date >= ? and shipping_date <= ?", iupopkId, strings.ToUpper(typeTransaction), input.DateFrom, input.DateTo).Find(&transactionReport).Error

	if errFind != nil {
		return transactionReport, errFind
	}

	return transactionReport, nil
}
