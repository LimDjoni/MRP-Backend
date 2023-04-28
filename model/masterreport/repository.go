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

	queryFilter := fmt.Sprintf("seller_id = %v AND transaction_type = 'DN' AND shipping_date >= '%s' AND shipping_date <= '%s'", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("id ASC").Find(&listTransactions).Error

	if errFind != nil {
		return reportDmoOuput, errFind
	}

	// Query Electric Assignment
	var electricAssignment electricassignment.ElectricAssignment

	errFindElectricAssigment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&electricAssignment).Error

	if errFindElectricAssigment != nil {
		return reportDmoOuput, errFindElectricAssigment
	}

	reportDmoOuput.ElectricAssignment = electricAssignment
	var electricAssignmentEndUser []electricassignmentenduser.ElectricAssignmentEndUser

	errFindElectricAssigmentEndUser := r.db.Where("electric_assignment_id = ?", electricAssignment.ID).Order("id desc").Find(&electricAssignmentEndUser).Error

	if errFindElectricAssigmentEndUser != nil {
		return reportDmoOuput, errFindElectricAssigmentEndUser
	}

	// Query Caf Assignment
	var cafAssignment cafassignment.CafAssignment

	errFindCafAssignment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&cafAssignment).Error

	if errFindCafAssignment != nil {
		return reportDmoOuput, errFindCafAssignment
	}

	reportDmoOuput.CafAssignment = cafAssignment

	var cafAssignmentEndUser []cafassignmentenduser.CafAssignmentEndUser

	errFindCafAssigmentEndUser := r.db.Where("caf_assignment_id = ?", cafAssignment.ID).Order("id desc").Find(&cafAssignmentEndUser).Error

	if errFindCafAssigmentEndUser != nil {
		return reportDmoOuput, errFindCafAssigmentEndUser
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

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.January += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 2:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.February += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 3:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.March += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 4:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.April += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 5:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.May += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 6:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.June += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 7:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.July += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 8:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.August += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 9:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.September += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 10:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.October += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 11:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.November += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 12:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					reportDmoOuput.RecapElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							reportDmoOuput.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					reportDmoOuput.RecapNonElectricity.December += v.QuantityUnloading
					reportDmoOuput.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							reportDmoOuput.RealizationCafAssignment += v.QuantityUnloading
						}
					}
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

	queryFilter := fmt.Sprintf("seller_id = %v AND transaction_type = 'DN' AND shipping_date >= '%s' AND shipping_date <= '%s' AND is_not_claim = false", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType").Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("id ASC").Find(&listTransactions).Error

	if errFind != nil {
		return realizationOutput, errFind
	}

	// Query Electric Assignment
	var electricAssignment electricassignment.ElectricAssignment

	errFindElectricAssigment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&electricAssignment).Error

	if errFindElectricAssigment != nil {
		return realizationOutput, errFindElectricAssigment
	}

	var electricAssignmentEndUser []electricassignmentenduser.ElectricAssignmentEndUser

	errFindElectricAssigmentEndUser := r.db.Where("electric_assignment_id = ?", electricAssignment.ID).Order("id desc").Find(&electricAssignmentEndUser).Error

	if errFindElectricAssigmentEndUser != nil {
		return realizationOutput, errFindElectricAssigmentEndUser
	}

	// Query Caf Assignment
	var cafAssignment cafassignment.CafAssignment

	errFindCafAssignment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&cafAssignment).Error

	if errFindCafAssignment != nil {
		return realizationOutput, errFindCafAssignment
	}

	var cafAssignmentEndUser []cafassignmentenduser.CafAssignmentEndUser

	errFindCafAssigmentEndUser := r.db.Where("caf_assignment_id = ?", cafAssignment.ID).Order("id desc").Find(&cafAssignmentEndUser).Error

	if errFindCafAssigmentEndUser != nil {
		return realizationOutput, errFindCafAssigmentEndUser
	}

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()
		var transactionTemp RealizationTransaction

		switch int(month) {
		case 1:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.January = append(realizationOutput.Caf.January, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.January = append(realizationOutput.Electric.January, transactionTemp)
					}
				}
			}
		case 2:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.February = append(realizationOutput.Caf.February, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.February = append(realizationOutput.Electric.February, transactionTemp)
					}
				}
			}
		case 3:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.March = append(realizationOutput.Caf.March, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.March = append(realizationOutput.Electric.March, transactionTemp)
					}
				}
			}
		case 4:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.April = append(realizationOutput.Caf.April, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.April = append(realizationOutput.Electric.April, transactionTemp)
					}
				}
			}
		case 5:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.May = append(realizationOutput.Caf.May, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.May = append(realizationOutput.Electric.May, transactionTemp)
					}
				}
			}
		case 6:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.June = append(realizationOutput.Caf.June, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.June = append(realizationOutput.Electric.June, transactionTemp)
					}
				}
			}
		case 7:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.July = append(realizationOutput.Caf.July, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.July = append(realizationOutput.Electric.July, transactionTemp)
					}
				}
			}
		case 8:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.August = append(realizationOutput.Caf.August, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.August = append(realizationOutput.Electric.August, transactionTemp)
					}
				}
			}
		case 9:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.September = append(realizationOutput.Caf.September, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.September = append(realizationOutput.Electric.September, transactionTemp)
					}
				}
			}
		case 10:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.October = append(realizationOutput.Caf.October, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.October = append(realizationOutput.Electric.October, transactionTemp)
					}
				}
			}
		case 11:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.November = append(realizationOutput.Caf.November, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.November = append(realizationOutput.Electric.November, transactionTemp)
					}
				}
			}
		case 12:
			if v.DmoBuyer != nil {
				for _, caf := range cafAssignmentEndUser {
					if caf.EndUserString == v.DmoBuyer.CompanyName {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Caf.December = append(realizationOutput.Caf.December, transactionTemp)
					}
				}
			}

			if v.DmoDestinationPortId != nil {
				for _, elec := range electricAssignmentEndUser {
					if elec.PortId == *v.DmoDestinationPortId {
						transactionTemp.ShippingDate = *v.ShippingDate
						if v.Customer != nil {
							transactionTemp.Trader = *v.Customer
						}

						if v.DmoBuyer != nil {
							transactionTemp.EndUser = *v.DmoBuyer
						}
						transactionTemp.Quantity = v.QuantityUnloading
						transactionTemp.QualityCaloriesAr = v.QualityCaloriesAr
						if v.Dmo != nil && v.Dmo.IsBastDocumentSigned {
							transactionTemp.IsBastOk = true
						} else {
							transactionTemp.IsBastOk = false
						}

						realizationOutput.Electric.December = append(realizationOutput.Electric.December, transactionTemp)
					}
				}
			}
		}
	}

	return realizationOutput, nil
}

func (r *repository) SaleDetailReport(year string, iupopkId int) (SaleDetail, error) {
	var saleDetail SaleDetail

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	var companyElectricity []string
	var companyNonElectricity []string

	var realizationCompanyElectricity []string
	var realizationCompanyCaf []string

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

	queryFilter := fmt.Sprintf("seller_id = %v AND transaction_type = 'DN' AND shipping_date >= '%s' AND shipping_date <= '%s'", iupopkId, startFilter, endFilter)

	errFind := r.db.Preload(clause.Associations).Preload("Customer.IndustryType").Preload("DmoBuyer.IndustryType").Where(queryFilter).Order("id ASC").Find(&listTransactions).Error

	if errFind != nil {
		return saleDetail, errFind
	}

	// Query Electric Assignment
	var electricAssignment electricassignment.ElectricAssignment

	errFindElectricAssigment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&electricAssignment).Error

	if errFindElectricAssigment != nil {
		return saleDetail, errFindElectricAssigment
	}

	saleDetail.ElectricAssignment = electricAssignment

	var electricAssignmentEndUser []electricassignmentenduser.ElectricAssignmentEndUser

	errFindElectricAssigmentEndUser := r.db.Preload(clause.Associations).Where("electric_assignment_id = ?", electricAssignment.ID).Order("id desc").Find(&electricAssignmentEndUser).Error

	if errFindElectricAssigmentEndUser != nil {
		return saleDetail, errFindElectricAssigmentEndUser
	}

	// Query Caf Assignment
	var cafAssignment cafassignment.CafAssignment

	errFindCafAssignment := r.db.Where("year = ? AND iupopk_id = ?", year, iupopkId).First(&cafAssignment).Error

	if errFindCafAssignment != nil {
		return saleDetail, errFindCafAssignment
	}

	saleDetail.CafAssignment = cafAssignment

	var cafAssignmentEndUser []cafassignmentenduser.CafAssignmentEndUser

	errFindCafAssigmentEndUser := r.db.Preload(clause.Associations).Where("caf_assignment_id = ?", cafAssignment.ID).Order("id desc").Find(&cafAssignmentEndUser).Error

	if errFindCafAssigmentEndUser != nil {
		return saleDetail, errFindCafAssigmentEndUser
	}

	// Rkabs Query
	var rkabs []rkab.Rkab

	queryRkab := fmt.Sprintf("year = '%s' AND iupopk_id = %v", year, iupopkId)

	errFindRkab := r.db.Where(queryRkab).Order("id ASC").Find(&rkabs).Error

	if errFindRkab != nil {
		return saleDetail, errFindRkab
	}

	saleDetail.Rkabs = rkabs

	saleDetail.Electricity.January = make(map[string]float64)
	saleDetail.Electricity.February = make(map[string]float64)
	saleDetail.Electricity.March = make(map[string]float64)
	saleDetail.Electricity.April = make(map[string]float64)
	saleDetail.Electricity.May = make(map[string]float64)
	saleDetail.Electricity.June = make(map[string]float64)
	saleDetail.Electricity.July = make(map[string]float64)
	saleDetail.Electricity.August = make(map[string]float64)
	saleDetail.Electricity.September = make(map[string]float64)
	saleDetail.Electricity.October = make(map[string]float64)
	saleDetail.Electricity.November = make(map[string]float64)
	saleDetail.Electricity.December = make(map[string]float64)

	saleDetail.NonElectricity.January = make(map[string]float64)
	saleDetail.NonElectricity.February = make(map[string]float64)
	saleDetail.NonElectricity.March = make(map[string]float64)
	saleDetail.NonElectricity.April = make(map[string]float64)
	saleDetail.NonElectricity.May = make(map[string]float64)
	saleDetail.NonElectricity.June = make(map[string]float64)
	saleDetail.NonElectricity.July = make(map[string]float64)
	saleDetail.NonElectricity.August = make(map[string]float64)
	saleDetail.NonElectricity.September = make(map[string]float64)
	saleDetail.NonElectricity.October = make(map[string]float64)
	saleDetail.NonElectricity.November = make(map[string]float64)
	saleDetail.NonElectricity.December = make(map[string]float64)

	saleDetail.DetailRealizationCafAssignment.January = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.February = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.March = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.April = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.May = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.June = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.July = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.August = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.September = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.October = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.November = make(map[string]float64)
	saleDetail.DetailRealizationCafAssignment.December = make(map[string]float64)

	saleDetail.DetailRealizationElectricAssignment.January = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.February = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.March = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.April = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.May = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.June = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.July = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.August = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.September = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.October = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.November = make(map[string]float64)
	saleDetail.DetailRealizationElectricAssignment.December = make(map[string]float64)

	for _, v := range electricAssignmentEndUser {
		realizationCompanyElectricity = append(realizationCompanyElectricity, v.Port.Name)
	}

	for _, v := range cafAssignmentEndUser {
		realizationCompanyCaf = append(realizationCompanyCaf, v.EndUserString)
	}

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

		if v.IsNotClaim == false {
			switch int(month) {
			case 1:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.January += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.January += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.January[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.January[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.January[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.January += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.January[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.January[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.January[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.January[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.January[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.January += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							if _, ok := saleDetail.DetailRealizationCafAssignment.January[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.January[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.January[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 2:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.February += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.February += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.February[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.February[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.February[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.February += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.February[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.February[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.February[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.February[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.February[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.February += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							if _, ok := saleDetail.DetailRealizationCafAssignment.February[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.February[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.February[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 3:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.March += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.March += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.March[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.March[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.March[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.March += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.March[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.March[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.March[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.March[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.March[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.March += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {
							if _, ok := saleDetail.DetailRealizationCafAssignment.March[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.March[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.March[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 4:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.April += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.April += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.April[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.April[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.April[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.April += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.April[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.April[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.April[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.April[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.April[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.April += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.April[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.April[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.April[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 5:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.May += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.May += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.May[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.May[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.May[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.May += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.May[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.May[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.May[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.May[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.May[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.May += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.May[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.May[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.May[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 6:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.June += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.June += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.June[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.June[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.June[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.June += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.June[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.June[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.June[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.June[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.June[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.June += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.June[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.June[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.June[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 7:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.July += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.July += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.July[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.July[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.July[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.July += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.July[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.July[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.July[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.July[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.July[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.July += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.July[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.July[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.July[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 8:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.August += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.August += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.August[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.August[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.August[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.August += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.August[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.August[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.August[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.August[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.August[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.August += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.August[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.August[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.August[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 9:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.September += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.September += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.September[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.September[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.September[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.September += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.September[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.September[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.September[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.September[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.September[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.September += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.September[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.September[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.September[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 10:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.October += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.October += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.October[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.October[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.October[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.October += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.October[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.October[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.October[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.October[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.October[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.October += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.October[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.October[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.October[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 11:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.November += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.November += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.November[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.November[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.November[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.November += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.November[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.November[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.November[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.November[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.November[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.November += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.November[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.November[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.November[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
					}
				}
			case 12:
				if v.Destination != nil && v.Destination.Name == "Domestik" {
					saleDetail.Domestic.December += v.QuantityUnloading
				} else if v.Destination != nil && v.Destination.Name == "Ekspor" {
					saleDetail.Export.December += v.QuantityUnloading
				}
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := saleDetail.Electricity.December[v.DmoBuyer.CompanyName]; ok {
						saleDetail.Electricity.December[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.Electricity.December[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.Electricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapElectricity.December += v.QuantityUnloading
					saleDetail.RecapElectricity.Total += v.QuantityUnloading

					for _, electric := range electricAssignmentEndUser {

						if electric.PortId == *v.DmoDestinationPortId {
							if _, ok := saleDetail.DetailRealizationElectricAssignment.December[electric.Port.Name]; ok {
								saleDetail.DetailRealizationElectricAssignment.December[electric.Port.Name] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationElectricAssignment.December[electric.Port.Name] = v.QuantityUnloading
							}
							saleDetail.RealizationElectricAssignment += v.QuantityUnloading
						}
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := saleDetail.Electricity.December[v.DmoBuyer.CompanyName]; ok {
						saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						saleDetail.NonElectricity.December[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						saleDetail.NonElectricity.Total += v.QuantityUnloading
					}
					saleDetail.RecapNonElectricity.December += v.QuantityUnloading
					saleDetail.RecapNonElectricity.Total += v.QuantityUnloading

					for _, caf := range cafAssignmentEndUser {
						if caf.EndUserString == v.DmoBuyer.CompanyName {

							if _, ok := saleDetail.DetailRealizationCafAssignment.December[caf.EndUserString]; ok {
								saleDetail.DetailRealizationCafAssignment.December[caf.EndUserString] += v.QuantityUnloading
							} else {
								saleDetail.DetailRealizationCafAssignment.December[caf.EndUserString] = v.QuantityUnloading
							}
							saleDetail.RealizationCafAssignment += v.QuantityUnloading
						}
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
	saleDetail.RealizationCompanyElectricity = realizationCompanyElectricity
	saleDetail.RealizationCompanyCaf = realizationCompanyCaf
	return saleDetail, nil
}
