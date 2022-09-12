package transaction

import "ajebackend/model/minerba"

type DetailMinerba struct {
	Detail minerba.Minerba `json:"detail"`
	List []Transaction `json:"list"`
}
