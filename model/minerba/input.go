package minerba

type InputCreateMinerba struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
	ListDataDn []int `json:"list_data_dn" validate:"required,min=1"`
}

type InputUpdateDocumentMinerba struct {
	SP3MEDNDocumentLink string `json:"sp3medn_document_link" validate:"required,url"`
	RecapDmoDocumentLink string `json:"recap_dmo_document_link" validate:"required,url"`
	DetailDmoDocumentLink string `json:"detail_dmo_document_link" validate:"required,url"`
	SP3MELNDocumentLink *string `json:"sp3meln_document_link"`
	INSWExportDocumentLink *string `json:"insw_export_document_link"`
}
