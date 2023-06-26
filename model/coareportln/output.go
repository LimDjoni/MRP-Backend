package coareportln

import "ajebackend/model/transaction"

type CoaReportLnDetail struct {
	Detail          CoaReportLn               `json:"detail"`
	ListTransaction []transaction.Transaction `json:"list_transaction"`
}
