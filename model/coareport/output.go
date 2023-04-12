package coareport

import "ajebackend/model/transaction"

type CoaReportDetail struct {
	Detail          CoaReport                 `json:"detail"`
	ListTransaction []transaction.Transaction `json:"list_transaction"`
}
