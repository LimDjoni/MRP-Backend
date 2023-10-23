package isp

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/master/site"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type Isp struct {
	gorm.Model
	Name        string        `json:"name" gorm:"UNIQUE"`
	Latitude    string        `json:"latitude"`
	Longitude   string        `json:"longitude"`
	Quantity    float64       `json:"quantity"`
	SiteId      uint          `json:"site_id"`
	Site        site.Site     `json:"site"`
	IupopkId    uint          `json:"iupopk_id"`
	Iupopk      iupopk.Iupopk `json:"iupopk"`
	CreatedById uint          `json:"created_by_id"`
	CreatedBy   user.User     `json:"created_by"`
	UpdatedById uint          `json:"updated_by_id"`
	UpdatedBy   user.User     `json:"updated_by"`
}
