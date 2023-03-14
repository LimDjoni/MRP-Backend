package notification

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/user"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Status   string        `json:"status"`
	Type     string        `json:"type"`
	Period   string        `json:"period"`
	Document string        `json:"document"`
	EndUser  string        `json:"end_user"`
	UserId   uint          `json:"user_id"`
	User     user.User     `json:"user"`
	IupopkId uint          `json:"iupopk_id"`
	Iupopk   iupopk.Iupopk `json:"iupopk"`
}
