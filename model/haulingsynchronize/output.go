package haulingsynchronize

import (
	"ajebackend/model/master/contractor"
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/role"
	"ajebackend/model/user"
	"ajebackend/model/useriupopk"
	"ajebackend/model/userrole"
)

type MasterDataIsp struct {
	Contractor []contractor.Contractor `json:"contractor"`
	Isp        []isp.Isp               `json:"isp"`
	Iupopk     []iupopk.Iupopk         `json:"iupopk"`
	Jetty      []jetty.Jetty           `json:"jetty"`
	Pit        []pit.Pit               `json:"pit"`
	Truck      []TruckOutput           `json:"truck"`
	Role       []role.Role             `json:"role"`
	User       []user.User             `json:"user"`
	UserIupopk []useriupopk.UserIupopk `json:"user_iupopk"`
	UserRole   []userrole.UserRole     `json:"user_role"`
}

type MasterDataJetty struct {
	Contractor []contractor.Contractor `json:"contractor"`
	Isp        []isp.Isp               `json:"isp"`
	Jetty      []jetty.Jetty           `json:"jetty"`
	Pit        []pit.Pit               `json:"pit"`
	Iupopk     []iupopk.Iupopk         `json:"iupopk"`
	Role       []role.Role             `json:"role"`
	User       []user.User             `json:"user"`
	UserIupopk []useriupopk.UserIupopk `json:"user_iupopk"`
	UserRole   []userrole.UserRole     `json:"user_role"`
}

type TruckOutput struct {
	UpdatedAt     string                 `json:"updated_at" gorm:"DATETIME"`
	CreatedAt     string                 `json:"created_at" gorm:"DATETIME"`
	Rfid          *string                `json:"rfid" gorm:"UNIQUE"`
	Code          string                 `json:"code" gorm:"UNIQUE;primaryKey"`
	NumberLambung string                 `json:"number_lambung"`
	TruckModel    string                 `json:"truck_model"`
	Tara          float64                `json:"tara"`
	Capacity      float64                `json:"capacity"`
	ContractorId  *uint                  `json:"contractor_id"`
	Contractor    *contractor.Contractor `json:"contractor"`
	CreatedById   uint                   `json:"created_by_id"`
	CreatedBy     user.User              `json:"created_by"`
	UpdatedById   uint                   `json:"updated_by_id"`
	UpdatedBy     user.User              `json:"updated_by"`
	IsActive      bool                   `json:"is_active"`
}
