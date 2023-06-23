package masterreport

type TransactionReportInput struct {
	DateFrom string `json:"date_from" validate:"required,DateValidation"`
	DateTo   string `json:"date_to" validate:"required,DateValidation"`
}
