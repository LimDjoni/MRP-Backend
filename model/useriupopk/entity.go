package useriupopk

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type UserIupopk struct {
	gorm.Model
	UserId   uint          `json:"user_id"`
	User     user.User     `json:"user"`
	IupopkId uint          `json:"iupopk_id"`
	Iupopk   iupopk.Iupopk `json:"iupopk" gorm:"constraint:OnDelete:CASCADE;"`
}
