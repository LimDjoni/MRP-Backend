package coareport

import (
	"ajebackend/helper"
	"ajebackend/model/transaction"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
	GetDetailTransactionCoaReport(id int, iupopkId int) (CoaReportDetail, error)
	RequestCreateExcelCoaReport(inputRequestCreateExcel InputRequestCreateUploadCoaReport) (map[string]interface{}, error)
	ListCoaReport(page int, sortFilter SortFilterCoaReport, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionCoaReport(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	listTransaction, listTransactionErr := s.repository.GetTransactionCoaReport(dateFrom, dateTo, iupopkId)

	return listTransaction, listTransactionErr
}

func (s *service) GetDetailTransactionCoaReport(id int, iupopkId int) (CoaReportDetail, error) {
	detailCoaReport, detailCoaReportErr := s.repository.GetDetailTransactionCoaReport(id, iupopkId)

	return detailCoaReport, detailCoaReportErr
}

func (s *service) RequestCreateExcelCoaReport(inputRequestCreateExcel InputRequestCreateUploadCoaReport) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/coa_report"
	body, bodyErr := json.Marshal(inputRequestCreateExcel)

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

func (s *service) ListCoaReport(page int, sortFilter SortFilterCoaReport, iupopkId int) (Pagination, error) {
	listCoaReport, listCoaReportErr := s.repository.ListCoaReport(page, sortFilter, iupopkId)

	return listCoaReport, listCoaReportErr
}
