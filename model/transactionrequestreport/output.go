package transactionrequestreport

import (
	"ajebackend/model/masterreport"
)

type TransactionRequestReportDetail struct {
	Detail             TransactionRequestReport         `json:"detail"`
	ListDnTransactions []masterreport.TransactionReport `json:"list_dn_transactions"`
	ListLnTransactions []masterreport.TransactionReport `json:"list_ln_transactions"`
}

type TransactionRequestReportPreview struct {
	ListDnTransactions []masterreport.TransactionReport `json:"list_dn_transactions"`
	ListLnTransactions []masterreport.TransactionReport `json:"list_ln_transactions"`
}
