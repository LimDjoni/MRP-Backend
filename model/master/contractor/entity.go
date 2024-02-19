package contractor

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type Contractor struct {
	gorm.Model
	Name        string        `json:"name"`
	Address     string        `json:"address"`
	PhoneNumber string        `json:"phone_number"`
	Email       string        `json:"email" gorm:"UNIQUE"`
	IupopkId    uint          `json:"iupopk_id"`
	Iupopk      iupopk.Iupopk `json:"iupopk"`
	CreatedById uint          `json:"created_by_id"`
	CreatedBy   user.User     `json:"created_by"`
	UpdatedById uint          `json:"updated_by_id"`
	UpdatedBy   user.User     `json:"updated_by"`
}
