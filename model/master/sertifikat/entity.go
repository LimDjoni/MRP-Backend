package sertifikat

import (
	"gorm.io/gorm"
)

type Sertifikat struct {
	gorm.Model
	EmployeeId    uint   `json:"employee_id"`
	DateEffective string `json:"date_effective"`
	Sertifikat    string `json:"sertifikat"`
	Remark        string `json:"remark"`
}
