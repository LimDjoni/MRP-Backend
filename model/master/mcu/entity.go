package mcu

import (
	"gorm.io/gorm"
)

type MCU struct {
	gorm.Model
	EmployeeId uint   `json:"employee_id"`
	DateMCU    string `json:"date_mcu"`
	DateEndMCU string `json:"date_end_mcu"`
	HasilMCU   string `json:"hasil_mcu"`
	MCU        string `json:"mcu"`
}
