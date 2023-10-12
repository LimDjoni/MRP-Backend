package transactiontoisp

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/pit"
	"ajebackend/model/master/site"
	"ajebackend/model/master/truck"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type TransactionToIsp struct {
	gorm.Model
	IdNumber              string        `json:"id_number"`
	TruckId               *uint         `json:"truck_id"`
	Truck                 *truck.Truck  `json:"truck"`
	NettQuantity          float64       `json:"nett_quantity"`
	RitaseQuantity        float64       `json:"ritase_quantity"`
	Gar                   float64       `json:"gar"`
	IupopkId              uint          `json:"iupopk_id"`
	Iupopk                iupopk.Iupopk `json:"iupopk"`
	PitId                 uint          `json:"pit_id"`
	Pit                   pit.Pit       `json:"pit"`
	IspId                 uint          `json:"isp_id"`
	Isp                   isp.Isp       `json:"isp"`
	SiteId                uint          `json:"site_id"`
	Site                  site.Site     `json:"site"`
	Seam                  string        `json:"seam"`
	Category              string        `json:"category"`
	TopTruckPhotoLink     string        `json:"top_truck_photo_link"`
	TopTruckPhotoPath     string        `json:"top_truck_photo_path"`
	LambungTruckPhotoLink string        `json:"lambung_truck_photo_link"`
	LambungTruckPhotoPath string        `json:"lambung_truck_photo_path"`
	SurveyDocumentLink    string        `json:"survey_document_link"`
	CreatedById           uint          `json:"created_by_id"`
	CreatedBy             user.User     `json:"created_by"`
	UpdatedById           uint          `json:"updated_by_id"`
	UpdatedBy             user.User     `json:"updated_by"`
}
