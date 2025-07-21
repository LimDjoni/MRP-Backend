package history

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	EmployeeId     uint   `json:"employee_id"`
	StatusTerakhir string `json:"status_terakhir"`
	Tanggal        string `json:"tanggal"`
	Keterangan     string `json:"keterangan"`
}
