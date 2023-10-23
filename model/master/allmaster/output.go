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
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/navycompany"
	"ajebackend/model/master/navyship"
	"ajebackend/model/master/pabeanoffice"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/portinsw"
	"ajebackend/model/master/portlocation"
	"ajebackend/model/master/ports"
	"ajebackend/model/master/salessystem"
	"ajebackend/model/master/surveyor"
	"ajebackend/model/master/trader"
	"ajebackend/model/master/tugboat"
	"ajebackend/model/master/unit"
	"ajebackend/model/master/vessel"
)

type MasterData struct {
	Barge                []barge.Barge                               `json:"barge"`
	CategoryIndustryType []categoryindustrytype.CategoryIndustryType `json:"category_industry_type"`
	Company              []company.Company                           `json:"company"`
	Country              []country.Country                           `json:"country"`
	Currency             []currency.Currency                         `json:"currency"`
	Destination          []destination.Destination                   `json:"destination"`
	DocumentType         []documenttype.DocumentType                 `json:"document_type"`
	IndustryType         []industrytype.IndustryType                 `json:"industry_type"`
	InsuranceCompany     []insurancecompany.InsuranceCompany         `json:"insurance_company"`
	Iupopk               []iupopk.Iupopk                             `json:"iupopk"`
	Jetty                []jetty.Jetty                               `json:"jetty"`
	NavyCompany          []navycompany.NavyCompany                   `json:"navy_company"`
	NavyShip             []navyship.NavyShip                         `json:"navy_ship"`
	PabeanOffice         []pabeanoffice.PabeanOffice                 `json:"pabean_office"`
	Pit                  []pit.Pit                                   `json:"pit"`
	Isp                  []isp.Isp                                   `json:"isp"`
	PortInsw             []portinsw.PortInsw                         `json:"port_insw"`
	PortLocation         []portlocation.PortLocation                 `json:"port_location"`
	Ports                []ports.Port                                `json:"ports"`
	SalesSystem          []salessystem.SalesSystem                   `json:"sales_system"`
	Surveyor             []surveyor.Surveyor                         `json:"surveyor"`
	Trader               []trader.Trader                             `json:"trader"`
	Tugboat              []tugboat.Tugboat                           `json:"tugboat"`
	Unit                 []unit.Unit                                 `json:"unit"`
	Vessel               []vessel.Vessel                             `json:"vessel"`
}
