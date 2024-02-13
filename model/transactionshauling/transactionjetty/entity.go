package transactionjetty

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/truck"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type TransactionJetty struct {
	gorm.Model
	IdNumber              string        `json:"id_number"`
	Truck                 truck.Truck   `json:"truck" `
	TruckId               uint          `json:"truck_id"`
	NettQuantity          float64       `json:"nett_quantity"`
	IupopkId              uint          `json:"iupopk_id"`
	Iupopk                iupopk.Iupopk `json:"iupopk"`
	IspCode               *string       `json:"isp_code"`
	Isp                   *isp.Isp      `json:"isp" gorm:"foreignKey:IspCode;references:Code"`
	PitCode               *string       `json:"pit_code"`
	Pit                   *pit.Pit      `json:"pit" gorm:"foreignKey:PitCode;references:Code"`
	JettyCode             string        `json:"jetty_code"`
	Jetty                 jetty.Jetty   `json:"jetty" gorm:"foreignKey:JettyCode;references:Code"`
	Seam                  string        `json:"seam"`
	Gar                   float64       `json:"gar"`
	ClockInDate           string        `json:"clock_in_date" gorm:"DATETIME"`
	TopTruckPhotoLink     string        `json:"top_truck_photo_link"`
	TopTruckPhotoPath     string        `json:"top_truck_photo_path"`
	LambungTruckPhotoLink string        `json:"lambung_truck_photo_link"`
	LambungTruckPhotoPath string        `json:"lambung_truck_photo_path"`
	WeightBridgePhotoLink string        `json:"weight_bridge_photo_link"`
	WeightBridgePhotoPath string        `json:"weight_bridge_photo_path"`
	CreatedById           uint          `json:"created_by_id"`
	CreatedBy             user.User     `json:"created_by"`
	UpdatedById           uint          `json:"updated_by_id"`
	UpdatedBy             user.User     `json:"updated_by"`
	IsFailedUpload        bool          `json:"is_failed_upload"`
}
