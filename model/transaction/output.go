package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
)

type DetailMinerba struct {
	Detail minerba.Minerba `json:"detail"`
	List []Transaction `json:"list"`
}

type DetailDmo struct {
	Detail dmo.Dmo `json:"detail"`
	List []Transaction `json:"list"`
}
