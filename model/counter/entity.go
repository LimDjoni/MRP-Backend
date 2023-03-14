package counter

import (
	"ajebackend/model/master/iupopk"

	"gorm.io/gorm"
)

type Counter struct {
	gorm.Model
	IupopkId      uint           `json:"iupopk_id"`
	Iupopk        *iupopk.Iupopk `json:"iupopk" gorm:"constraint:OnDelete:CASCADE;"`
	TransactionDn int            `json:"transaction_dn`
	TransactionLn int            `json:"transaction_ln`
	GroupingMvDn  int            `json:"grouping_mv_dn`
	GroupingMvLn  int            `json:"grouping_mv_ln`
	Sp3medn       int            `json:"sp3medn`
	Sp3meln       int            `json:"sp3meln`
	BaEndUser     int            `json:"ba_end_user`
	Dmo           int            `json:"dmo`
	Production    int            `json:"production`
	Insw          int            `json:"insw"`
}
