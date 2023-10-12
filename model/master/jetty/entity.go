package jetty

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type Jetty struct {
	gorm.Model
	Name        string        `json:"name" gorm:"UNIQUE"`
	Latitude    string        `json:"latitude"`
	Longitude   string        `json:"longitude"`
	Quantity    float64       `json:"quantity"`
	IupopkId    uint          `json:"iupopk_id"`
	Iupopk      iupopk.Iupopk `json:"iupopk"`
	CreatedById uint          `json:"created_by_id"`
	CreatedBy   user.User     `json:"created_by"`
	UpdatedById uint          `json:"updated_by_id"`
	UpdatedBy   user.User     `json:"updated_by"`
}
