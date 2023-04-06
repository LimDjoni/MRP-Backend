package allmaster

import (
	"ajebackend/model/master/barge"
	"ajebackend/model/master/company"
	"ajebackend/model/master/country"
	"ajebackend/model/master/currency"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/documenttype"
	"ajebackend/model/master/industrytype"
	"ajebackend/model/master/insurancecompany"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/navycompany"
	"ajebackend/model/master/navyship"
	"ajebackend/model/master/pabeanoffice"
	"ajebackend/model/master/portinsw"
	"ajebackend/model/master/portlocation"
	"ajebackend/model/master/ports"
	"ajebackend/model/master/salessystem"
	"ajebackend/model/master/surveyor"
	"ajebackend/model/master/trader"
	"ajebackend/model/master/tugboat"
	"ajebackend/model/master/unit"
	"ajebackend/model/master/vessel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListMasterData() (MasterData, error)
	FindIupopk(iupopkId int) (iupopk.Iupopk, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListMasterData() (MasterData, error) {
	var masterData MasterData

	var barge []barge.Barge
	var company []company.Company
	var country []country.Country
	var currency []currency.Currency
	var destination []destination.Destination
	var documentType []documenttype.DocumentType
	var industryType []industrytype.IndustryType
	var insuranceCompany []insurancecompany.InsuranceCompany
	var iupopk []iupopk.Iupopk
	var navyCompany []navycompany.NavyCompany
	var navyShip []navyship.NavyShip
	var pabeanOffice []pabeanoffice.PabeanOffice
	var portInsw []portinsw.PortInsw
	var portLocation []portlocation.PortLocation
	var ports []ports.Port
	var salesSystem []salessystem.SalesSystem
	var surveyor []surveyor.Surveyor
	var trader []trader.Trader
	var tugboat []tugboat.Tugboat
	var unit []unit.Unit
	var vessel []vessel.Vessel

	findBargeErr := r.db.Order("id desc").Find(&barge).Error

	if findBargeErr != nil {
		return masterData, findBargeErr
	}

	findCompanyErr := r.db.Order("id desc").Preload(clause.Associations).Find(&company).Error

	if findCompanyErr != nil {
		return masterData, findCompanyErr
	}

	findCountryErr := r.db.Order("id desc").Find(&country).Error

	if findCountryErr != nil {
		return masterData, findCountryErr
	}

	findCurrencyErr := r.db.Order("id desc").Find(&currency).Error

	if findCurrencyErr != nil {
		return masterData, findCurrencyErr
	}

	findDestinationErr := r.db.Order("id desc").Find(&destination).Error

	if findDestinationErr != nil {
		return masterData, findDestinationErr
	}

	findDocumentTypeErr := r.db.Order("id desc").Find(&documentType).Error

	if findDocumentTypeErr != nil {
		return masterData, findDocumentTypeErr
	}

	findIndustryTypeErr := r.db.Order("id desc").Find(&industryType).Error

	if findIndustryTypeErr != nil {
		return masterData, findIndustryTypeErr
	}

	findInsuranceCompanyErr := r.db.Order("id desc").Find(&insuranceCompany).Error

	if findInsuranceCompanyErr != nil {
		return masterData, findInsuranceCompanyErr
	}

	findIupopkErr := r.db.Order("id desc").Find(&iupopk).Error

	if findIupopkErr != nil {
		return masterData, findIupopkErr
	}

	findNavyCompanyErr := r.db.Order("id desc").Find(&navyCompany).Error

	if findNavyCompanyErr != nil {
		return masterData, findNavyCompanyErr
	}

	findNavyShipErr := r.db.Order("id desc").Find(&navyShip).Error

	if findNavyShipErr != nil {
		return masterData, findNavyShipErr
	}

	findPabeanOfficeErr := r.db.Order("id desc").Find(&pabeanOffice).Error

	if findPabeanOfficeErr != nil {
		return masterData, findPabeanOfficeErr
	}

	findPortInswErr := r.db.Order("id desc").Find(&portInsw).Error

	if findPortInswErr != nil {
		return masterData, findPortInswErr
	}

	findPortLocationErr := r.db.Order("id desc").Find(&portLocation).Error

	if findPortLocationErr != nil {
		return masterData, findPortLocationErr
	}

	findPortsErr := r.db.Order("id desc").Preload(clause.Associations).Find(&ports).Error

	if findPortsErr != nil {
		return masterData, findPortsErr
	}

	findSalesSystemErr := r.db.Order("id desc").Find(&salesSystem).Error

	if findSalesSystemErr != nil {
		return masterData, findSalesSystemErr
	}

	findSurveyorErr := r.db.Order("id desc").Find(&surveyor).Error

	if findSurveyorErr != nil {
		return masterData, findSurveyorErr
	}

	findTraderErr := r.db.Order("id desc").Preload("Company.IndustryType").Find(&trader).Error

	if findTraderErr != nil {
		return masterData, findTraderErr
	}

	findTugboatErr := r.db.Order("id desc").Find(&tugboat).Error

	if findTugboatErr != nil {
		return masterData, findTugboatErr
	}

	findUnitErr := r.db.Order("id desc").Find(&unit).Error

	if findUnitErr != nil {
		return masterData, findUnitErr
	}

	findVesselErr := r.db.Order("id desc").Find(&vessel).Error

	if findVesselErr != nil {
		return masterData, findVesselErr
	}

	masterData.Barge = barge
	masterData.Company = company
	masterData.Country = country
	masterData.Currency = currency
	masterData.Destination = destination
	masterData.DocumentType = documentType
	masterData.IndustryType = industryType
	masterData.InsuranceCompany = insuranceCompany
	masterData.Iupopk = iupopk
	masterData.NavyCompany = navyCompany
	masterData.NavyShip = navyShip
	masterData.PabeanOffice = pabeanOffice
	masterData.PortInsw = portInsw
	masterData.PortLocation = portLocation
	masterData.Ports = ports
	masterData.SalesSystem = salesSystem
	masterData.Surveyor = surveyor
	masterData.Trader = trader
	masterData.Tugboat = tugboat
	masterData.Unit = unit
	masterData.Vessel = vessel

	return masterData, nil
}

func (r *repository) FindIupopk(iupopkId int) (iupopk.Iupopk, error) {
	var iupopk iupopk.Iupopk

	errFind := r.db.Where("id = ?", iupopkId).First(&iupopk).Error

	if errFind != nil {
		return iupopk, errFind
	}

	return iupopk, nil
}
