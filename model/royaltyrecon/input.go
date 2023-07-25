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
