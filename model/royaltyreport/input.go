package royaltyreport

import "ajebackend/model/master/iupopk"

type InputUpdateDocumentRoyaltyReport struct {
	Data []map[string]interface{} `json:"data"`
}

type InputRequestCreateUploadRoyaltyReport struct {
	Authorization   string              `json:"authorization"`
	RoyaltyReport   RoyaltyReport       `json:"royalty_report"`
	ListTransaction []RoyaltyReportData `json:"list_transaction"`
	Iupopk          iupopk.Iupopk       `json:"iupopk"`
}
