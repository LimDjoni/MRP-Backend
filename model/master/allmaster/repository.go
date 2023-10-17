package allmaster

import (
	"ajebackend/model/master/barge"
	"ajebackend/model/master/categoryindustrytype"
	"ajebackend/model/master/company"
	"ajebackend/model/master/country"
	"ajebackend/model/master/currency"
	"ajebackend/model/master/destination"
	"ajebackend/model/master/documenttype"
	"ajebackend/model/master/industrytype"
	"ajebackend/model/master/insurancecompany"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
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
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListMasterData(iupopkId int) (MasterData, error)
	FindIupopk(iupopkId int) (iupopk.Iupopk, error)
	CreateBarge(input InputBarge) (barge.Barge, error)
	CreateTugboat(input InputTugboat) (tugboat.Tugboat, error)
	CreateVessel(input InputVessel) (vessel.Vessel, error)
	CreatePortLocation(input InputPortLocation) (portlocation.PortLocation, error)
	CreatePort(input InputPort) (ports.Port, error)
	CreateCompany(input InputCompany) (company.Company, error)
	CreateTrader(input InputTrader) (trader.Trader, error)
	CreateIndustryType(input InputIndustryType) (industrytype.IndustryType, error)
	UpdateBarge(id int, input InputBarge) (barge.Barge, error)
	UpdateTugboat(id int, input InputTugboat) (tugboat.Tugboat, error)
	UpdateVessel(id int, input InputVessel) (vessel.Vessel, error)
	UpdatePortLocation(id int, input InputPortLocation) (portlocation.PortLocation, error)
	UpdatePort(id int, input InputPort) (ports.Port, error)
	UpdateCompany(id int, input InputCompany) (company.Company, error)
	UpdateTrader(id int, input InputTrader) (trader.Trader, error)
	UpdateIndustryType(id int, input InputIndustryType) (industrytype.IndustryType, error)
	DeleteBarge(id int) (bool, error)
	DeleteTugboat(id int) (bool, error)
	DeleteVessel(id int) (bool, error)
	DeletePortLocation(id int) (bool, error)
	DeletePort(id int) (bool, error)
	DeleteCompany(id int) (bool, error)
	DeleteTrader(id int) (bool, error)
	DeleteIndustryType(id int) (bool, error)
	ListCompany() ([]company.Company, error)
	ListTrader() ([]trader.Trader, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListMasterData(iupopkId int) (MasterData, error) {
	var masterData MasterData

	var barge []barge.Barge
	var categoryIndustryType []categoryindustrytype.CategoryIndustryType
	var company []company.Company
	var country []country.Country
	var currency []currency.Currency
	var destination []destination.Destination
	var documentType []documenttype.DocumentType
	var industryType []industrytype.IndustryType
	var insuranceCompany []insurancecompany.InsuranceCompany
	var iupopk []iupopk.Iupopk
	var jetty []jetty.Jetty
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

	findBargeErr := r.db.Order("name asc").Order("created_at desc").Find(&barge).Error

	if findBargeErr != nil {
		return masterData, findBargeErr
	}

	findCategoryIndustryTypeErr := r.db.Find(&categoryIndustryType).Error

	if findCategoryIndustryTypeErr != nil {
		return masterData, findCategoryIndustryTypeErr
	}

	findCompanyErr := r.db.Order("company_name asc").Order("created_at desc").Preload(clause.Associations).Preload("IndustryType.CategoryIndustryType").Find(&company).Error
	if findCompanyErr != nil {
		return masterData, findCompanyErr
	}

	findCountryErr := r.db.Order("name asc").Order("created_at desc").Find(&country).Error

	if findCountryErr != nil {
		return masterData, findCountryErr
	}

	findCurrencyErr := r.db.Order("name asc").Order("created_at desc").Find(&currency).Error

	if findCurrencyErr != nil {
		return masterData, findCurrencyErr
	}

	findDestinationErr := r.db.Order("name asc").Order("created_at desc").Find(&destination).Error

	if findDestinationErr != nil {
		return masterData, findDestinationErr
	}

	findDocumentTypeErr := r.db.Order("name asc").Order("created_at desc").Find(&documentType).Error

	if findDocumentTypeErr != nil {
		return masterData, findDocumentTypeErr
	}

	findIndustryTypeErr := r.db.Preload(clause.Associations).Order("name asc").Order("created_at desc").Find(&industryType).Error

	if findIndustryTypeErr != nil {
		return masterData, findIndustryTypeErr
	}

	findInsuranceCompanyErr := r.db.Order("name asc").Order("created_at desc").Find(&insuranceCompany).Error

	if findInsuranceCompanyErr != nil {
		return masterData, findInsuranceCompanyErr
	}

	findIupopkErr := r.db.Order("name asc").Order("created_at desc").Find(&iupopk).Error

	if findIupopkErr != nil {
		return masterData, findIupopkErr
	}

	findJettyErr := r.db.Order("name asc").Preload(clause.Associations).Order("created_at desc").Where("iupopk_id = ?", iupopkId).Find(&jetty).Error

	if findJettyErr != nil {
		return masterData, findJettyErr
	}

	findNavyCompanyErr := r.db.Order("name asc").Order("created_at desc").Find(&navyCompany).Error

	if findNavyCompanyErr != nil {
		return masterData, findNavyCompanyErr
	}

	findNavyShipErr := r.db.Order("name asc").Order("created_at desc").Find(&navyShip).Error

	if findNavyShipErr != nil {
		return masterData, findNavyShipErr
	}

	findPabeanOfficeErr := r.db.Order("name asc").Order("created_at desc").Find(&pabeanOffice).Error

	if findPabeanOfficeErr != nil {
		return masterData, findPabeanOfficeErr
	}

	findPortInswErr := r.db.Order("name asc").Order("created_at desc").Find(&portInsw).Error

	if findPortInswErr != nil {
		return masterData, findPortInswErr
	}

	findPortLocationErr := r.db.Order("name asc").Order("created_at desc").Find(&portLocation).Error

	if findPortLocationErr != nil {
		return masterData, findPortLocationErr
	}

	findPortsErr := r.db.Preload(clause.Associations, func(db *gorm.DB) *gorm.DB {
		return db.Order("name asc").Order("created_at desc")
	}).Find(&ports).Error

	if findPortsErr != nil {
		return masterData, findPortsErr
	}

	findSalesSystemErr := r.db.Order("name asc").Order("created_at desc").Find(&salesSystem).Error

	if findSalesSystemErr != nil {
		return masterData, findSalesSystemErr
	}

	findSurveyorErr := r.db.Order("name asc").Order("created_at desc").Find(&surveyor).Error

	if findSurveyorErr != nil {
		return masterData, findSurveyorErr
	}

	findTraderErr := r.db.Order("trader_name asc").Order("created_at desc").Preload("Company.IndustryType.CategoryIndustryType").Find(&trader).Error

	if findTraderErr != nil {
		return masterData, findTraderErr
	}

	findTugboatErr := r.db.Order("name asc").Order("created_at desc").Find(&tugboat).Error

	if findTugboatErr != nil {
		return masterData, findTugboatErr
	}

	findUnitErr := r.db.Order("name asc").Order("created_at desc").Find(&unit).Error

	if findUnitErr != nil {
		return masterData, findUnitErr
	}

	findVesselErr := r.db.Order("name asc").Order("created_at desc").Find(&vessel).Error

	if findVesselErr != nil {
		return masterData, findVesselErr
	}

	masterData.Barge = barge
	masterData.CategoryIndustryType = categoryIndustryType
	masterData.Company = company
	masterData.Country = country
	masterData.Currency = currency
	masterData.Destination = destination
	masterData.DocumentType = documentType
	masterData.IndustryType = industryType
	masterData.InsuranceCompany = insuranceCompany
	masterData.Iupopk = iupopk
	masterData.Jetty = jetty
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

func (r *repository) CreateBarge(input InputBarge) (barge.Barge, error) {
	var createdBarge barge.Barge

	createdBarge.Name = input.Name
	createdBarge.Height = input.Height
	createdBarge.Deadweight = input.Deadweight
	createdBarge.MinimumQuantity = input.MinimumQuantity
	createdBarge.MaximumQuantity = input.MaximumQuantity

	errCreate := r.db.Create(&createdBarge).Error

	if errCreate != nil {
		return createdBarge, errCreate
	}

	return createdBarge, nil
}

func (r *repository) CreateTugboat(input InputTugboat) (tugboat.Tugboat, error) {
	var createdTugboat tugboat.Tugboat

	createdTugboat.Name = input.Name
	createdTugboat.Height = input.Height
	createdTugboat.Deadweight = input.Deadweight
	createdTugboat.MinimumQuantity = input.MinimumQuantity
	createdTugboat.MaximumQuantity = input.MaximumQuantity

	errCreate := r.db.Create(&createdTugboat).Error

	if errCreate != nil {
		return createdTugboat, errCreate
	}

	return createdTugboat, nil
}

func (r *repository) CreateVessel(input InputVessel) (vessel.Vessel, error) {
	var createdVessel vessel.Vessel

	createdVessel.Name = input.Name
	createdVessel.Deadweight = input.Deadweight
	createdVessel.MinimumQuantity = input.MinimumQuantity
	createdVessel.MaximumQuantity = input.MaximumQuantity

	errCreate := r.db.Create(&createdVessel).Error

	if errCreate != nil {
		return createdVessel, errCreate
	}

	return createdVessel, nil
}

func (r *repository) CreatePortLocation(input InputPortLocation) (portlocation.PortLocation, error) {
	var createdPortLocation portlocation.PortLocation

	createdPortLocation.Name = input.Name

	errCreate := r.db.Create(&createdPortLocation).Error

	if errCreate != nil {
		return createdPortLocation, errCreate
	}

	return createdPortLocation, nil
}

func (r *repository) CreatePort(input InputPort) (ports.Port, error) {
	var createdPort ports.Port

	createdPort.Name = input.Name
	createdPort.PortLocationId = input.PortLocationId
	createdPort.IsLoadingPort = input.IsLoadingPort
	createdPort.IsUnloadingPort = input.IsUnloadingPort
	createdPort.IsDmoDestinationPort = input.IsDmoDestinationPort

	errCreate := r.db.Create(&createdPort).Error

	if errCreate != nil {
		return createdPort, errCreate
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", createdPort.ID).First(&createdPort).Error

	if errFind != nil {
		return createdPort, errFind
	}

	return createdPort, nil
}

func (r *repository) CreateCompany(input InputCompany) (company.Company, error) {
	var createdCompany company.Company

	createdCompany.CompanyName = input.CompanyName
	createdCompany.Address = input.Address
	createdCompany.Province = input.Province
	createdCompany.PhoneNumber = input.PhoneNumber
	createdCompany.FaxNumber = input.FaxNumber
	createdCompany.IndustryTypeId = input.IndustryTypeId
	createdCompany.IsTrader = input.IsTrader
	createdCompany.IsEndUser = input.IsEndUser

	errCreate := r.db.Create(&createdCompany).Error

	if errCreate != nil {
		return createdCompany, errCreate
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", createdCompany.ID).First(&createdCompany).Error

	if errFind != nil {
		return createdCompany, errFind
	}

	return createdCompany, nil
}

func (r *repository) CreateTrader(input InputTrader) (trader.Trader, error) {
	var createdTrader trader.Trader

	var findCompany company.Company
	errFindCompany := r.db.Where("id = ?", input.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return createdTrader, errFindCompany
	}

	createdTrader.TraderName = input.TraderName
	createdTrader.Position = input.Position

	var email string

	if input.Email != nil {
		email = strings.ToLower(*input.Email)
	}
	if email != "" {
		createdTrader.Email = &email
	}

	createdTrader.CompanyId = uint(input.CompanyId)

	errCreate := r.db.Create(&createdTrader).Error

	if errCreate != nil {
		return createdTrader, errCreate
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", createdTrader.ID).First(&createdTrader).Error

	if errFind != nil {
		return createdTrader, errFind
	}

	return createdTrader, nil
}

func (r *repository) CreateIndustryType(input InputIndustryType) (industrytype.IndustryType, error) {
	var createdIndustryType industrytype.IndustryType

	createdIndustryType.Name = input.Name
	createdIndustryType.CategoryIndustryTypeId = input.CategoryIndustryTypeId
	createdIndustryType.SystemCategory = input.SystemCategory

	errCreate := r.db.Create(&createdIndustryType).Error

	if errCreate != nil {
		return createdIndustryType, errCreate
	}

	return createdIndustryType, nil
}

func (r *repository) UpdateBarge(id int, input InputBarge) (barge.Barge, error) {
	var updatedBarge barge.Barge

	errFind := r.db.Where("id = ?", id).First(&updatedBarge).Error

	if errFind != nil {
		return updatedBarge, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name
	upd["height"] = input.Height
	upd["deadweight"] = input.Deadweight
	upd["minimum_quantity"] = input.MinimumQuantity
	upd["maximum_quantity"] = input.MaximumQuantity

	errCreate := r.db.Model(&updatedBarge).Updates(upd).Error

	if errCreate != nil {
		return updatedBarge, errCreate
	}

	return updatedBarge, nil
}

func (r *repository) UpdateTugboat(id int, input InputTugboat) (tugboat.Tugboat, error) {
	var updatedTugboat tugboat.Tugboat

	errFind := r.db.Where("id = ?", id).First(&updatedTugboat).Error

	if errFind != nil {
		return updatedTugboat, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name
	upd["height"] = input.Height
	upd["deadweight"] = input.Deadweight
	upd["minimum_quantity"] = input.MinimumQuantity
	upd["maximum_quantity"] = input.MaximumQuantity

	errCreate := r.db.Model(&updatedTugboat).Updates(upd).Error

	if errCreate != nil {
		return updatedTugboat, errCreate
	}

	return updatedTugboat, nil
}

func (r *repository) UpdateVessel(id int, input InputVessel) (vessel.Vessel, error) {
	var updatedVessel vessel.Vessel

	errFind := r.db.Where("id = ?", id).First(&updatedVessel).Error

	if errFind != nil {
		return updatedVessel, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name
	upd["deadweight"] = input.Deadweight
	upd["minimum_quantity"] = input.MinimumQuantity
	upd["maximum_quantity"] = input.MaximumQuantity

	errCreate := r.db.Model(&updatedVessel).Updates(upd).Error

	if errCreate != nil {
		return updatedVessel, errCreate
	}

	return updatedVessel, nil
}

func (r *repository) UpdatePortLocation(id int, input InputPortLocation) (portlocation.PortLocation, error) {
	var updatedPortLocation portlocation.PortLocation

	errFind := r.db.Where("id = ?", id).First(&updatedPortLocation).Error

	if errFind != nil {
		return updatedPortLocation, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name

	errCreate := r.db.Model(&updatedPortLocation).Updates(upd).Error

	if errCreate != nil {
		return updatedPortLocation, errCreate
	}

	return updatedPortLocation, nil
}

func (r *repository) UpdatePort(id int, input InputPort) (ports.Port, error) {
	var updatedPort ports.Port

	errFind := r.db.Where("id = ?", id).First(&updatedPort).Error

	if errFind != nil {
		return updatedPort, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name
	upd["port_location_id"] = input.PortLocationId
	upd["is_loading_port"] = input.IsLoadingPort
	upd["is_unloading_port"] = input.IsUnloadingPort
	upd["is_dmo_destination_port"] = input.IsDmoDestinationPort

	errCreate := r.db.Model(&updatedPort).Updates(upd).Error

	if errCreate != nil {
		return updatedPort, errCreate
	}

	return updatedPort, nil
}

func (r *repository) UpdateCompany(id int, input InputCompany) (company.Company, error) {
	var updatedCompany company.Company

	errFind := r.db.Where("id = ?", id).First(&updatedCompany).Error

	if errFind != nil {
		return updatedCompany, errFind
	}

	upd := make(map[string]interface{})

	upd["company_name"] = input.CompanyName
	upd["industry_type_id"] = input.IndustryTypeId
	upd["address"] = input.Address
	upd["province"] = input.Province
	upd["phone_number"] = input.PhoneNumber
	upd["fax_number"] = input.FaxNumber
	upd["is_trader"] = input.IsTrader
	upd["is_end_user"] = input.IsEndUser

	updateErr := r.db.Model(&updatedCompany).Updates(upd).Error

	if updateErr != nil {
		return updatedCompany, updateErr
	}

	return updatedCompany, nil
}

func (r *repository) UpdateTrader(id int, input InputTrader) (trader.Trader, error) {
	var updatedTrader trader.Trader

	var email string

	if input.Email != nil {
		email = strings.ToLower(*input.Email)
	}

	errFind := r.db.Preload(clause.Associations).Where("id = ?", id).First(&updatedTrader).Error

	if errFind != nil {
		return updatedTrader, errFind
	}

	var findCompany company.Company

	errFindCompany := r.db.Where("id = ?", input.CompanyId).First(&findCompany).Error

	if errFindCompany != nil {
		return updatedTrader, errFindCompany
	}

	upd := make(map[string]interface{})

	upd["trader_name"] = input.TraderName
	upd["position"] = input.Position
	if email != "" {
		upd["email"] = email
	} else {
		upd["email"] = nil
	}
	upd["company_id"] = input.CompanyId

	updateErr := r.db.Model(&updatedTrader).Updates(upd).Error

	if updateErr != nil {
		return updatedTrader, updateErr
	}

	return updatedTrader, nil
}

func (r *repository) UpdateIndustryType(id int, input InputIndustryType) (industrytype.IndustryType, error) {
	var updatedIndustryType industrytype.IndustryType

	errFind := r.db.Preload(clause.Associations).Where("id = ?", id).First(&updatedIndustryType).Error

	if errFind != nil {
		return updatedIndustryType, errFind
	}

	upd := make(map[string]interface{})

	upd["name"] = input.Name
	upd["category_industry_type_id"] = input.CategoryIndustryTypeId
	upd["system_category"] = input.SystemCategory

	updateErr := r.db.Model(&updatedIndustryType).Updates(upd).Error

	if updateErr != nil {
		return updatedIndustryType, updateErr
	}

	return updatedIndustryType, nil
}

func (r *repository) DeleteBarge(id int) (bool, error) {
	var findBarge barge.Barge

	errFind := r.db.Where("id = ?", id).First(&findBarge).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findBarge).Error

	if errDelete != nil {
		return false, errFind
	}

	return true, nil
}

func (r *repository) DeleteTugboat(id int) (bool, error) {
	var findTugboat tugboat.Tugboat

	errFind := r.db.Where("id = ?", id).First(&findTugboat).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findTugboat).Error

	if errDelete != nil {
		return false, errFind
	}

	return true, nil
}

func (r *repository) DeleteVessel(id int) (bool, error) {
	var findVessel vessel.Vessel

	errFind := r.db.Where("id = ?", id).First(&findVessel).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findVessel).Error

	if errDelete != nil {
		return false, errFind
	}

	return true, nil
}

func (r *repository) DeletePortLocation(id int) (bool, error) {
	var findPortLocation portlocation.PortLocation

	errFind := r.db.Where("id = ?", id).First(&findPortLocation).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findPortLocation).Error

	if errDelete != nil {
		return false, errFind
	}

	return true, nil
}

func (r *repository) DeletePort(id int) (bool, error) {
	var findPort ports.Port

	errFind := r.db.Where("id = ?", id).First(&findPort).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findPort).Error

	if errDelete != nil {
		return false, errFind
	}

	return true, nil
}

func (r *repository) DeleteCompany(id int) (bool, error) {
	var findCompany company.Company

	errFind := r.db.Where("id = ?", id).First(&findCompany).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Delete(&findCompany).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) DeleteTrader(id int) (bool, error) {
	var findTrader trader.Trader

	errFind := r.db.Where("id = ?", id).First(&findTrader).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ?", id).Delete(&findTrader).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) DeleteIndustryType(id int) (bool, error) {
	var findIndustryType industrytype.IndustryType

	errFind := r.db.Where("id = ?", id).First(&findIndustryType).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ?", id).Delete(&findIndustryType).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) ListCompany() ([]company.Company, error) {
	var listCompany []company.Company

	errFind := r.db.Preload(clause.Associations).Preload("IndustryType.CategoryIndustryType").Find(&listCompany).Error

	if errFind != nil {
		return listCompany, errFind
	}

	return listCompany, nil
}

func (r *repository) ListTrader() ([]trader.Trader, error) {
	var listTrader []trader.Trader

	errFind := r.db.Preload(clause.Associations).Preload("Company.IndustryType.CategoryIndustryType").Find(&listTrader).Error

	if errFind != nil {
		return listTrader, errFind
	}

	return listTrader, nil
}
