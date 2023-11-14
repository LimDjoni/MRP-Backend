package transactionjetty

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/truck"
	"ajebackend/model/user"
)

type InputTransactionJetty struct {
	IdNumber              string        `json:"id_number"`
	Truck                 truck.Truck   `json:"truck" `
	TruckId               uint          `json:"truck_id"`
	NettQuantity          float64       `json:"nett_quantity"`
	IupopkId              uint          `json:"iupopk_id"`
	Iupopk                iupopk.Iupopk `json:"iupopk"`
	IspId                 *uint         `json:"isp_id"`
	Isp                   *isp.Isp      `json:"isp"`
	PitId                 *uint         `json:"pit_id"`
	Pit                   *pit.Pit      `json:"pit"`
	JettyId               uint          `json:"jetty_id"`
	Jetty                 jetty.Jetty   `json:"jetty"`
	Seam                  string        `json:"seam"`
	ClockInDate           string        `json:"clock_in_date" gorm:"DATETIME"`
	TopTruckPhotoLink     string        `json:"top_truck_photo_link"`
	TopTruckPhotoPath     string        `json:"top_truck_photo_path"`
	LambungTruckPhotoLink string        `json:"lambung_truck_photo_link"`
	LambungTruckPhotoPath string        `json:"lambung_truck_photo_path"`
	CreatedById           uint          `json:"created_by_id"`
	CreatedBy             user.User     `json:"created_by"`
	UpdatedById           uint          `json:"updated_by_id"`
	UpdatedBy             user.User     `json:"updated_by"`
}
