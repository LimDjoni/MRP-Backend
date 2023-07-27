package royaltyrecon

import "ajebackend/model/master/iupopk"

type InputUpdateDocumentRoyaltyRecon struct {
	Data []map[string]interface{} `json:"data"`
}

type InputRequestCreateUploadRoyaltyRecon struct {
	Authorization   string             `json:"authorization"`
	RoyaltyRecon    RoyaltyRecon       `json:"royalty_recon"`
	ListTransaction []RoyaltyReconData `json:"list_transaction"`
	Iupopk          iupopk.Iupopk      `json:"iupopk"`
}

type InputRoyaltyRecon struct {
	DateFrom string `json:"date_from" validate:"required,DateValidation"`
	DateTo   string `json:"date_to" validate:"required,DateValidation"`
}
