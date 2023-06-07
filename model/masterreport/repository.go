package masterreport

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/electricassignment"
	"ajebackend/model/electricassignmentenduser"
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

	queryFilter := fmt.Sprintf("seller_id = %v AND transaction_type = 'DN' AND shipping_date >= '%s' AND shipping_date <= '%s' AND dmo_id IS NOT NULL", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

	if errFind != nil {
		return reportDmoOuput, errFind
	}

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

		if v.IsNotClaim == false {
			switch int(month) {
			case 1:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.January += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.January += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 2:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 3:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 4:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 5:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 6:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 7:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 8:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 9:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 10:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 11:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			case 12:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading
				}
			}
		} else {
			switch int(month) {
			case 1:
				reportDmoOuput.NotClaimable.January += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 2:
				reportDmoOuput.NotClaimable.February += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 3:
				reportDmoOuput.NotClaimable.March += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 4:
				reportDmoOuput.NotClaimable.April += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 5:
				reportDmoOuput.NotClaimable.May += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 6:
				reportDmoOuput.NotClaimable.June += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 7:
				reportDmoOuput.NotClaimable.July += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 8:
				reportDmoOuput.NotClaimable.August += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 9:
				reportDmoOuput.NotClaimable.September += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 10:
				reportDmoOuput.NotClaimable.October += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 11:
				reportDmoOuput.NotClaimable.November += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			case 12:
				reportDmoOuput.NotClaimable.December += v.QuantityUnloading
				reportDmoOuput.NotClaimable.Total += v.QuantityUnloading
			}
		}
	}

	return reportDmoOuput, nil
}

func (r *repository) RealizationReport(year string, iupopkId int) (RealizationOutput, error) {
	var realizationOutput RealizationOutput

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	var listTransactions []transaction.Transaction

	queryFilter := fmt.Sprintf("seller_id = %v AND transaction_type = 'DN' AND shipping_date >= '%s' AND shipping_date <= '%s' AND is_not_claim = false AND dmo_id IS NOT NULL", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType").Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

	if errFind != nil {
		return realizationOutput, errFind
	}

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()
		var transactionTemp RealizationTransaction

		switch int(month) {
		case 1:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.January = append(realizationOutput.Electric.January, transactionTemp)
					} else {
						realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.January = append(realizationOutput.NonElectric.January, transactionTemp)
			}
		case 2:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.February = append(realizationOutput.Electric.February, transactionTemp)
					} else {
						realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.February = append(realizationOutput.NonElectric.February, transactionTemp)
			}
		case 3:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.March = append(realizationOutput.Electric.March, transactionTemp)
					} else {
						realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.March = append(realizationOutput.NonElectric.March, transactionTemp)
			}
		case 4:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.April = append(realizationOutput.Electric.April, transactionTemp)
					} else {
						realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.April = append(realizationOutput.NonElectric.April, transactionTemp)
			}
		case 5:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.May = append(realizationOutput.Electric.May, transactionTemp)
					} else {
						realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.May = append(realizationOutput.NonElectric.May, transactionTemp)
			}
		case 6:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.June = append(realizationOutput.Electric.June, transactionTemp)
					} else {
						realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.June = append(realizationOutput.NonElectric.June, transactionTemp)
			}
		case 7:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.July = append(realizationOutput.Electric.July, transactionTemp)
					} else {
						realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.July = append(realizationOutput.NonElectric.July, transactionTemp)
			}
		case 8:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.August = append(realizationOutput.Electric.August, transactionTemp)
					} else {
						realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.August = append(realizationOutput.NonElectric.August, transactionTemp)
			}
		case 9:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.September = append(realizationOutput.Electric.September, transactionTemp)
					} else {
						realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.September = append(realizationOutput.NonElectric.September, transactionTemp)
			}
		case 10:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.October = append(realizationOutput.Electric.October, transactionTemp)
					} else {
						realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.October = append(realizationOutput.NonElectric.October, transactionTemp)
			}
		case 11:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.November = append(realizationOutput.Electric.November, transactionTemp)
					} else {
						realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
					}
				} else {
					realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
				}
			} else {
				realizationOutput.NonElectric.November = append(realizationOutput.NonElectric.November, transactionTemp)
			}
		case 12:
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
					if v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
						realizationOutput.Electric.December = append(realizationOutput.Electric.December, transactionTemp)
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
	var listTransactions []transaction.Transaction

	queryFilter := fmt.Sprintf("seller_id = %v AND shipping_date >= '%s' AND shipping_date <= '%s'", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType").Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("shipping_date ASC").Find(&listTransactions).Error

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

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

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

			switch int(month) {
			case 1:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.January += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.January += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 2:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.February += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.February += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 3:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.March += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.March += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 4:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.April += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.April += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 5:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.May += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.May += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}
				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 6:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.June += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.June += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 7:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.July += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.July += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 8:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.August += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.August += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 9:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.September += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.September += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 10:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.October += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.October += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 11:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.November += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.November += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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

			case 12:
				if v.TransactionType == "DN" {
					saleDetail.Domestic.December += v.QuantityUnloading
					saleDetail.Domestic.Total += v.QuantityUnloading
				} else {
					saleDetail.Export.December += v.QuantityUnloading
					saleDetail.Export.Total += v.QuantityUnloading
				}

				if v.DmoId != nil {
					if v.DmoBuyer != nil {
						if v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
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
						} else {
							if _, ok := saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName]; ok {
								saleDetail.NonElectricity.Total += v.QuantityUnloading

								if v.DmoBuyer.IndustryType != nil {
									if _, ok := companyNonElectricity[v.DmoBuyer.CompanyName]; ok {
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}

									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
										if !helperString(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name) {
											companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
										}
									} else {
										companyNonElectricity[v.DmoBuyer.CompanyName] = append(companyNonElectricity[v.DmoBuyer.CompanyName], v.DmoBuyer.IndustryType.Name)
									}
									saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName][v.DmoBuyer.IndustryType.Name] += v.QuantityUnloading
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
			switch int(month) {
			case 1:
				saleDetail.NotClaimable.January += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 2:
				saleDetail.NotClaimable.February += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 3:
				saleDetail.NotClaimable.March += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 4:
				saleDetail.NotClaimable.April += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 5:
				saleDetail.NotClaimable.May += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 6:
				saleDetail.NotClaimable.June += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 7:
				saleDetail.NotClaimable.July += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 8:
				saleDetail.NotClaimable.August += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 9:
				saleDetail.NotClaimable.September += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 10:
				saleDetail.NotClaimable.October += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 11:
				saleDetail.NotClaimable.November += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			case 12:
				saleDetail.NotClaimable.December += v.QuantityUnloading
				saleDetail.NotClaimable.Total += v.QuantityUnloading
			}
		}
	}

	saleDetail.CompanyElectricity = companyElectricity
	saleDetail.CompanyNonElectricity = companyNonElectricity
	return saleDetail, nil
}
