package doh

import (
	"gorm.io/gorm"
)

type DOH struct {
	gorm.Model
	EmployeeId    uint   `json:"employee_id"`
	TanggalDoh    string `json:"tanggal_doh"`
	TanggalEndDoh string `json:"tanggal_end_doh"`
	PT            string `json:"pt"`
	Penempatan    string `json:"penempatan"`
	StatusKontrak string `json:"status_kontrak"`
}
