package allmaster

import (
	"ajebackend/model/master/barge"
	"ajebackend/model/master/company"
	"ajebackend/model/master/industrytype"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/portlocation"
	"ajebackend/model/master/ports"
	"ajebackend/model/master/trader"
	"ajebackend/model/master/tugboat"
	"ajebackend/model/master/vessel"
)

type Service interface {
	ListMasterData() (MasterData, error)
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

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListMasterData() (MasterData, error) {
	listMasterData, listMasterDataErr := s.repository.ListMasterData()

	return listMasterData, listMasterDataErr
}

func (s *service) FindIupopk(iupopkId int) (iupopk.Iupopk, error) {
	findIupopk, findIupopkErr := s.repository.FindIupopk(iupopkId)

	return findIupopk, findIupopkErr
}

func (s *service) CreateBarge(input InputBarge) (barge.Barge, error) {
	data, dataErr := s.repository.CreateBarge(input)
	return data, dataErr
}

func (s *service) CreateTugboat(input InputTugboat) (tugboat.Tugboat, error) {
	data, dataErr := s.repository.CreateTugboat(input)
	return data, dataErr
}

func (s *service) CreateVessel(input InputVessel) (vessel.Vessel, error) {
	data, dataErr := s.repository.CreateVessel(input)
	return data, dataErr
}

func (s *service) CreateCompany(input InputCompany) (company.Company, error) {
	data, dataErr := s.repository.CreateCompany(input)
	return data, dataErr
}

func (s *service) CreateTrader(input InputTrader) (trader.Trader, error) {
	data, dataErr := s.repository.CreateTrader(input)
	return data, dataErr
}

func (s *service) CreateIndustryType(input InputIndustryType) (industrytype.IndustryType, error) {
	data, dataErr := s.repository.CreateIndustryType(input)
	return data, dataErr
}

func (s *service) CreatePortLocation(input InputPortLocation) (portlocation.PortLocation, error) {
	data, dataErr := s.repository.CreatePortLocation(input)
	return data, dataErr
}

func (s *service) CreatePort(input InputPort) (ports.Port, error) {
	data, dataErr := s.repository.CreatePort(input)
	return data, dataErr
}

func (s *service) UpdateBarge(id int, input InputBarge) (barge.Barge, error) {
	data, dataErr := s.repository.UpdateBarge(id, input)
	return data, dataErr
}

func (s *service) UpdateTugboat(id int, input InputTugboat) (tugboat.Tugboat, error) {
	data, dataErr := s.repository.UpdateTugboat(id, input)
	return data, dataErr
}

func (s *service) UpdateVessel(id int, input InputVessel) (vessel.Vessel, error) {
	data, dataErr := s.repository.UpdateVessel(id, input)
	return data, dataErr
}

func (s *service) UpdatePortLocation(id int, input InputPortLocation) (portlocation.PortLocation, error) {
	data, dataErr := s.repository.UpdatePortLocation(id, input)
	return data, dataErr
}

func (s *service) UpdatePort(id int, input InputPort) (ports.Port, error) {
	data, dataErr := s.repository.UpdatePort(id, input)
	return data, dataErr
}

func (s *service) UpdateCompany(id int, input InputCompany) (company.Company, error) {
	data, dataErr := s.repository.UpdateCompany(id, input)
	return data, dataErr
}

func (s *service) UpdateTrader(id int, input InputTrader) (trader.Trader, error) {
	data, dataErr := s.repository.UpdateTrader(id, input)
	return data, dataErr
}

func (s *service) UpdateIndustryType(id int, input InputIndustryType) (industrytype.IndustryType, error) {
	data, dataErr := s.repository.UpdateIndustryType(id, input)
	return data, dataErr
}

func (s *service) DeleteBarge(id int) (bool, error) {
	data, dataErr := s.repository.DeleteBarge(id)
	return data, dataErr
}

func (s *service) DeleteTugboat(id int) (bool, error) {
	data, dataErr := s.repository.DeleteTugboat(id)
	return data, dataErr
}

func (s *service) DeleteVessel(id int) (bool, error) {
	data, dataErr := s.repository.DeleteVessel(id)
	return data, dataErr
}

func (s *service) DeletePortLocation(id int) (bool, error) {
	data, dataErr := s.repository.DeletePortLocation(id)
	return data, dataErr
}

func (s *service) DeletePort(id int) (bool, error) {
	data, dataErr := s.repository.DeletePort(id)
	return data, dataErr
}

func (s *service) DeleteCompany(id int) (bool, error) {
	data, dataErr := s.repository.DeleteCompany(id)
	return data, dataErr
}

func (s *service) DeleteTrader(id int) (bool, error) {
	data, dataErr := s.repository.DeleteTrader(id)
	return data, dataErr
}

func (s *service) DeleteIndustryType(id int) (bool, error) {
	data, dataErr := s.repository.DeleteIndustryType(id)
	return data, dataErr
}

func (s *service) ListCompany() ([]company.Company, error) {
	list, listErr := s.repository.ListCompany()

	return list, listErr
}

func (s *service) ListTrader() ([]trader.Trader, error) {
	list, listErr := s.repository.ListTrader()

	return list, listErr
}
