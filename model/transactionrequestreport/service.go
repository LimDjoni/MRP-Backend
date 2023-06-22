package transactionrequestreport

import (
	"ajebackend/helper"
	"ajebackend/model/masterreport"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	PreviewTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReportPreview, error)
	DetailTransactionReport(id int, iupopkId int) (TransactionRequestReportDetail, error)
	CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReportError(id int, iupopkId int) (TransactionRequestReport, error)
	DeleteTransactionReport() (bool, error)
	DeleteTransactionReportById(id int, iupopkId int) (bool, error)
	ListTransactionReport(page int, iupopkId int) (Pagination, error)
	ListDeletedTransactionReport() ([]TransactionRequestReport, error)
	ReqJobCreateTransactionReport(input InputRequestJobReportTransaction) (map[string]interface{}, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) PreviewTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReportPreview, error) {
	preview, previewErr := s.repository.PreviewTransactionReport(input, iupopkId)

	return preview, previewErr
}

func (s *service) DetailTransactionReport(id int, iupopkId int) (TransactionRequestReportDetail, error) {
	detail, detailErr := s.repository.DetailTransactionReport(id, iupopkId)

	return detail, detailErr
}

func (s *service) CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error) {
	transaction, transactionErr := s.repository.CreateTransactionReport(input, iupopkId)

	return transaction, transactionErr
}

func (s *service) UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error) {
	transaction, transactionErr := s.repository.UpdateTransactionReport(DnDocumentLink, LnDocumentLink, id, iupopkId)

	return transaction, transactionErr
}

func (s *service) UpdateTransactionReportError(id int, iupopkId int) (TransactionRequestReport, error) {
	transaction, transactionErr := s.repository.UpdateTransactionReportError(id, iupopkId)

	return transaction, transactionErr
}

func (s *service) DeleteTransactionReport() (bool, error) {
	isDeleted, isDeletedErr := s.repository.DeleteTransactionReport()

	return isDeleted, isDeletedErr
}

func (s *service) DeleteTransactionReportById(id int, iupopkId int) (bool, error) {
	isDeleted, isDeletedErr := s.repository.DeleteTransactionReportById(id, iupopkId)

	return isDeleted, isDeletedErr
}

func (s *service) ListTransactionReport(page int, iupopkId int) (Pagination, error) {
	list, listErr := s.repository.ListTransactionReport(page, iupopkId)

	return list, listErr
}

func (s *service) ListDeletedTransactionReport() ([]TransactionRequestReport, error) {
	list, listErr := s.repository.ListDeletedTransactionReport()

	return list, listErr
}

func (s *service) ReqJobCreateTransactionReport(input InputRequestJobReportTransaction) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/report/transactionrequest"
	body, bodyErr := json.Marshal(input)

	if bodyErr != nil {
		return res, bodyErr
	}
	var payload = bytes.NewBufferString(string(body))

	req, doReqErr := http.NewRequest("POST", urlPost, payload)

	if req != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}
	client := &http.Client{}
	resp, doReqErr := client.Do(req)

	if doReqErr != nil {
		return res, doReqErr
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, doReqErr
}
