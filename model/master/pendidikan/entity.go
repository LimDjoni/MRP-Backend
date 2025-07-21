package pendidikan

import (
	"gorm.io/gorm"
)

type Pendidikan struct {
	gorm.Model
	PendidikanLabel    string `json:"pendidikan_label"`
	PendidikanTerakhir string `json:"pendidikan_terakhir"`
	Jurusan            string `json:"jurusan"`
}
