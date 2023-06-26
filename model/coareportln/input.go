package coareportln

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/transaction"
)

type CoaReportLnInput struct {
	DateFrom string `json:"date_from" validate:"required,DateValidation"`
	DateTo   string `json:"date_to" validate:"required,DateValidation"`
}

type InputUpdateDocumentCoaReportLn struct {
	Data []map[string]interface{} `json:"data"`
}

type InputRequestCreateUploadCoaReportLn struct {
	Authorization   string                    `json:"authorization"`
	CoaReport       CoaReportLn               `json:"coa_report"`
	ListTransaction []transaction.Transaction `json:"list_transaction"`
	Iupopk          iupopk.Iupopk             `json:"iupopk"`
}

type SortFilterCoaReportLn struct {
	Field     string
	Sort      string
	DateStart string
	DateEnd   string
}
