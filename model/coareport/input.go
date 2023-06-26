package coareport

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/transaction"
)

type CoaReportInput struct {
	DateFrom string `json:"date_from" validate:"required,DateValidation"`
	DateTo   string `json:"date_to" validate:"required,DateValidation"`
}

type InputUpdateDocumentCoaReport struct {
	Data []map[string]interface{} `json:"data"`
}

type InputRequestCreateUploadCoaReport struct {
	Authorization   string                    `json:"authorization"`
	CoaReport       CoaReport                 `json:"coa_report"`
	ListTransaction []transaction.Transaction `json:"list_transaction"`
	Iupopk          iupopk.Iupopk             `json:"iupopk"`
}

type SortFilterCoaReport struct {
	Field     string
	Sort      string
	DateStart string
	DateEnd   string
}
