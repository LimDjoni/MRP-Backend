package transactionrequestreport

import "ajebackend/model/master/iupopk"

type InputRequestJobReportTransaction struct {
	Authorization string        `json:"authorization"`
	Id            uint          `json:"id"`
	Iupopk        iupopk.Iupopk `json:"iupopk"`
}
