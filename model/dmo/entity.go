package dmo

import (
	"gorm.io/gorm"
)

type Dmo struct {
	gorm.Model
	IdNumber                                      *string `json:"id_number" gorm:"UNIQUE"`
	Type                                          string  `json:"type"`
	Period                                        string  `json:"period"`
	DocumentDate                                  string  `json:"document_date" gorm:"type:DATE"`
	BargeTotalQuantity                            float64 `json:"barge_total_quantity"`
	BargeAdjustment                               float64 `json:"barge_adjustment"`
	BargeGrandTotalQuantity                       float64 `json:"barge_grand_total_quantity"`
	VesselTotalQuantity                           float64 `json:"vessel_total_quantity"`
	VesselAdjustment                              float64 `json:"vessel_adjustment"`
	VesselGrandTotalQuantity                      float64 `json:"vessel_grand_total_quantity"`
	IsDocumentCustom                              bool    `json:"is_document_custom"`
	ReconciliationLetterDocumentLink              *string `json:"reconciliation_letter_document_link"`
	IsReconciliationLetterDownloaded              bool    `json:"is_reconciliation_letter_downloaded"`
	IsReconciliationLetterSigned                  bool    `json:"is_reconciliation_letter_signed"`
	SignedReconciliationLetterDocumentLink        *string `json:"signed_reconciliation_letter_document_link"`
	BASTDocumentLink                              *string `json:"bast_document_link"`
	IsBastDocumentDownloaded                      bool    `json:"is_bast_document_downloaded"`
	IsBastDocumentSigned                          bool    `json:"is_bast_document_signed"`
	SignedBASTDocumentLink                        *string `json:"signed_bast_document_link"`
	StatementLetterDocumentLink                   *string `json:"statement_letter_document_link"`
	IsStatementLetterDownloaded                   bool    `json:"is_statement_letter_downloaded"`
	IsStatementLetterSigned                       bool    `json:"is_statement_letter_signed"`
	SignedStatementLetterDocumentLink             *string `json:"signed_statement_letter_document_link"`
	ReconciliationLetterEndUserDocumentLink       *string `json:"reconciliation_letter_end_user_document_link"`
	IsReconciliationLetterEndUserDownloaded       bool    `json:"is_reconciliation_letter_end_user_downloaded"`
	IsReconciliationLetterEndUserSigned           bool    `json:"is_reconciliation_letter_end_user_signed"`
	SignedReconciliationLetterEndUserDocumentLink *string `json:"signed_reconciliation_letter_end_user_document_link"`
	RecapDmoDocumentLink                          *string `json:"recap_dmo_document_link"`
	DetailDmoDocumentLink                         *string `json:"detail_dmo_document_link"`
}
